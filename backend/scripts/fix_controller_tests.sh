#!/bin/bash

# Fix controller test files to use expectedMessage instead of expectedBody

fix_file() {
    local file=$1
    echo "Fixing $file..."
    
    # Replace expectedBody map[string]interface{} with expectedMessage string in test structs
    sed -i 's/expectedBody[[:space:]]*map\[string\]interface{}{}/expectedMessage string/g' "$file"
    
    # Replace expectedBody assignments with expectedMessage
    sed -i 's/expectedBody:[[:space:]]*map\[string\]string{"message": "\([^"]*\)"}/expectedMessage: "\1"/g' "$file"
    
    # Replace references to expectedBody["message"] with expectedMessage
    sed -i 's/tt\.expectedBody\["message"\]/tt.expectedMessage/g' "$file"
    
    echo "Fixed $file"
}

# Fix all controller test files
for file in modules/orders/controller_test.go \
            modules/deliveries/controller_test.go \
            modules/payments/controller_test.go \
            modules/recipes/controller_test.go \
            modules/reviews/controller_test.go; do
    if [ -f "$file" ]; then
        fix_file "$file"
    fi
done

echo "All controller test files fixed!"

