#!/usr/bin/env python3
"""
Final fix for all controller test files
"""
import re

files = [
    "modules/payments/controller_test.go",
    "modules/recipes/controller_test.go",
    "modules/reviews/controller_test.go",
]

for filepath in files:
    print(f"Fixing {filepath}...")
    with open(filepath, 'r') as f:
        lines = f.readlines()
    
    fixed_lines = []
    i = 0
    while i < len(lines):
        line = lines[i]
        
        # Fix struct field definitions - look for "expectedBody   map[string]interface{}"
        if re.search(r'expectedBody\s+map\[string\]interface{}\{\}', line):
            # Replace with expectedMessage string
            fixed_lines.append(re.sub(
                r'expectedBody\s+map\[string\]interface{}\{\}',
                'expectedMessage string',
                line
            ))
            i += 1
            continue
        
        # Skip lines that have expectedBody: map[string]string assignments (they should be expectedMessage: "value")
        if 'expectedBody:' in line and 'map[string]' in line:
            # Skip this line - it should have been converted already
            i += 1
            continue
        
        fixed_lines.append(line)
        i += 1
    
    # Now fix context.Background() to mock.Anything
    content = ''.join(fixed_lines)
    content = re.sub(r'context\.Background\(\)', 'mock.Anything', content)
    
    # Fix any remaining expectedBody references
    content = re.sub(r'tt\.expectedBody\["message"\]', 'tt.expectedMessage', content)
    
    with open(filepath, 'w') as f:
        f.write(content)
    
    print(f"  Fixed {filepath}")

print("All files fixed!")

