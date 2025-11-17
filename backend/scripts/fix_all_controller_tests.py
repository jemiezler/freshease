#!/usr/bin/env python3
import re
import sys

files = [
    "modules/deliveries/controller_test.go",
    "modules/payments/controller_test.go",
    "modules/recipes/controller_test.go",
    "modules/reviews/controller_test.go",
]

for filepath in files:
    with open(filepath, 'r') as f:
        content = f.read()
    
    # Fix struct field definitions
    content = re.sub(
        r'(\s+)expectedBody\s+map\[string\]interface{}\{\}',
        r'\1expectedMessage string',
        content
    )
    
    # Fix references in assertions (should already be fixed but double-check)
    content = re.sub(
        r'tt\.expectedBody\["message"\]',
        r'tt.expectedMessage',
        content
    )
    
    with open(filepath, 'w') as f:
        f.write(content)
    print(f"Fixed {filepath}")

print("All controller test files fixed!")

