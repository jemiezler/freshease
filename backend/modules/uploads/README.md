# Image Upload API Documentation

This module provides image upload functionality using MinIO object storage. It supports uploading, retrieving, and deleting images with automatic validation and secure URL generation.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Usage Examples](#usage-examples)
- [Response Formats](#response-formats)
- [Error Handling](#error-handling)
- [File Validation](#file-validation)
- [Integration Examples](#integration-examples)

## Prerequisites

1. **MinIO Server**: Ensure MinIO is running via Docker Compose
   ```bash
   docker-compose up -d minio
   ```

2. **MinIO Access**: The MinIO console is available at `http://localhost:9001`
   - Default credentials: `minioadmin` / `minioadmin1234`

## Configuration

Configure MinIO settings via environment variables or `.env` file:

```env
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY_ID=minioadmin
MINIO_SECRET_ACCESS_KEY=minioadmin1234
MINIO_BUCKET=freshease
MINIO_USE_SSL=false
```

**Default Values:**
- `MINIO_ENDPOINT`: `localhost:9000`
- `MINIO_ACCESS_KEY_ID`: `minioadmin`
- `MINIO_SECRET_ACCESS_KEY`: `minioadmin1234`
- `MINIO_BUCKET`: `freshease`
- `MINIO_USE_SSL`: `false`

The bucket will be automatically created if it doesn't exist.

## API Endpoints

### 1. Upload Image

Upload an image file with optional folder specification.

**Endpoint:** `POST /api/uploads/images`

**Content-Type:** `multipart/form-data`

**Parameters:**
- `file` (required): Image file to upload
- `folder` (optional): Folder path to store the image (default: `images`)

**Example:**
```bash
curl -X POST http://localhost:8080/api/uploads/images \
  -F "file=@product-image.jpg" \
  -F "folder=products"
```

### 2. Upload Image to Specific Folder

Upload an image to a specific folder path via URL parameter.

**Endpoint:** `POST /api/uploads/images/:folder`

**Content-Type:** `multipart/form-data`

**Parameters:**
- `folder` (path parameter): Folder path (e.g., `products`, `users/avatars`)
- `file` (form data): Image file to upload

**Example:**
```bash
curl -X POST http://localhost:8080/api/uploads/images/users/avatars \
  -F "file=@avatar.png"
```

### 3. Delete Image

Delete an image from storage.

**Endpoint:** `DELETE /api/uploads/images/:path`

**Parameters:**
- `path` (path parameter): Object path (e.g., `images/uuid.jpg` or `products/uuid.png`)

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/uploads/images/products/550e8400-e29b-41d4-a716-446655440000.jpg
```

## Usage Examples

### cURL

**Basic Upload:**
```bash
curl -X POST http://localhost:8080/api/uploads/images \
  -F "file=@image.jpg"
```

**Upload to Products Folder:**
```bash
curl -X POST http://localhost:8080/api/uploads/images \
  -F "file=@product.jpg" \
  -F "folder=products"
```

**Upload to Nested Folder:**
```bash
curl -X POST http://localhost:8080/api/uploads/images/users/avatars \
  -F "file=@avatar.png"
```

**Delete Image:**
```bash
curl -X DELETE http://localhost:8080/api/uploads/images/products/550e8400-e29b-41d4-a716-446655440000.jpg
```

### JavaScript (Fetch API)

**Upload Image:**
```javascript
const formData = new FormData();
formData.append('file', fileInput.files[0]);
formData.append('folder', 'products');

const response = await fetch('http://localhost:8080/api/uploads/images', {
  method: 'POST',
  body: formData
});

const data = await response.json();
console.log('Image URL:', data.url);
console.log('Object Name:', data.object_name);
```

**Upload to Specific Folder:**
```javascript
const formData = new FormData();
formData.append('file', fileInput.files[0]);

const folder = 'users/avatars';
const response = await fetch(`http://localhost:8080/api/uploads/images/${folder}`, {
  method: 'POST',
  body: formData
});

const data = await response.json();
```

**Delete Image:**
```javascript
const objectPath = 'products/550e8400-e29b-41d4-a716-446655440000.jpg';
const response = await fetch(`http://localhost:8080/api/uploads/images/${objectPath}`, {
  method: 'DELETE'
});

const data = await response.json();
```

### React Example

```jsx
import { useState } from 'react';

function ImageUpload() {
  const [imageUrl, setImageUrl] = useState('');
  const [loading, setLoading] = useState(false);

  const handleUpload = async (event) => {
    const file = event.target.files[0];
    if (!file) return;

    setLoading(true);
    const formData = new FormData();
    formData.append('file', file);
    formData.append('folder', 'products');

    try {
      const response = await fetch('http://localhost:8080/api/uploads/images', {
        method: 'POST',
        body: formData
      });

      const data = await response.json();
      if (response.ok) {
        setImageUrl(data.url);
        // Store data.object_name in your database
      } else {
        console.error('Upload failed:', data.message);
      }
    } catch (error) {
      console.error('Error:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <input type="file" accept="image/*" onChange={handleUpload} />
      {loading && <p>Uploading...</p>}
      {imageUrl && <img src={imageUrl} alt="Uploaded" />}
    </div>
  );
}
```

### Go Example

```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "os"
)

func uploadImage(filePath, folder string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    var requestBody bytes.Buffer
    writer := multipart.NewWriter(&requestBody)

    // Add file
    part, err := writer.CreateFormFile("file", filePath)
    if err != nil {
        return err
    }
    io.Copy(part, file)

    // Add folder
    writer.WriteField("folder", folder)
    writer.Close()

    req, err := http.NewRequest("POST", "http://localhost:8080/api/uploads/images", &requestBody)
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
    return nil
}
```

## Response Formats

### Success Response (Upload)

**Status Code:** `200 OK`

```json
{
  "message": "Image uploaded successfully",
  "object_name": "products/550e8400-e29b-41d4-a716-446655440000.jpg",
  "url": "http://localhost:9000/freshease/products/550e8400-e29b-41d4-a716-446655440000.jpg?X-Amz-Algorithm=..."
}
```

**Fields:**
- `message`: Success message
- `object_name`: The path/name of the uploaded file in MinIO (store this in your database)
- `url`: Presigned URL valid for 7 days (use this for immediate display)

### Success Response (Delete)

**Status Code:** `200 OK`

```json
{
  "message": "Image deleted successfully"
}
```

### Error Response

**Status Code:** `400 Bad Request` or `500 Internal Server Error`

```json
{
  "message": "failed to upload image",
  "error": "invalid file type. Allowed types: [.jpg .jpeg .png .gif .webp]"
}
```

## Error Handling

Common error scenarios:

1. **Missing File:**
   ```json
   {
     "message": "file is required",
     "error": "..."
   }
   ```

2. **Invalid File Type:**
   ```json
   {
     "message": "failed to upload image",
     "error": "invalid file type. Allowed types: [.jpg .jpeg .png .gif .webp]"
   }
   ```

3. **File Too Large:**
   ```json
   {
     "message": "failed to upload image",
     "error": "file size exceeds 10MB limit"
   }
   ```

4. **MinIO Connection Error:**
   ```json
   {
     "message": "failed to upload image",
     "error": "failed to create MinIO client: ..."
   }
   ```

## File Validation

### Allowed File Types
- `.jpg` / `.jpeg`
- `.png`
- `.gif`
- `.webp`

### File Size Limit
- **Maximum:** 10 MB (10,485,760 bytes)

### File Naming
- Files are automatically renamed using UUIDs to prevent conflicts
- Format: `{folder}/{uuid}.{extension}`
- Example: `products/550e8400-e29b-41d4-a716-446655440000.jpg`

## Integration Examples

### Integrating with Product Creation

When creating a product, upload the image first, then use the returned `object_name` or `url`:

```javascript
// 1. Upload image
const uploadResponse = await fetch('http://localhost:8080/api/uploads/images', {
  method: 'POST',
  body: formData
});
const uploadData = await uploadResponse.json();

// 2. Create product with image URL
const productData = {
  name: "Fresh Apples",
  price: 5.99,
  description: "Organic red apples",
  image_url: uploadData.url,  // or uploadData.object_name
  // ... other fields
};

await fetch('http://localhost:8080/api/products', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify(productData)
});
```

### Integrating with User Avatar Upload

```javascript
// Upload avatar
const avatarFormData = new FormData();
avatarFormData.append('file', avatarFile);
avatarFormData.append('folder', 'users/avatars');

const uploadResponse = await fetch('http://localhost:8080/api/uploads/images', {
  method: 'POST',
  body: avatarFormData
});
const { url, object_name } = await uploadResponse.json();

// Update user profile
await fetch(`http://localhost:8080/api/users/${userId}`, {
  method: 'PUT',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ avatar: url })
});
```

### Deleting Image When Entity is Deleted

```javascript
// When deleting a product, also delete its image
async function deleteProduct(productId) {
  // 1. Get product to find image path
  const product = await fetch(`http://localhost:8080/api/products/${productId}`).then(r => r.json());
  
  // 2. Extract object_name from image_url or use stored object_name
  const objectName = extractObjectName(product.image_url);
  
  // 3. Delete image
  await fetch(`http://localhost:8080/api/uploads/images/${objectName}`, {
    method: 'DELETE'
  });
  
  // 4. Delete product
  await fetch(`http://localhost:8080/api/products/${productId}`, {
    method: 'DELETE'
  });
}
```

## Best Practices

1. **Store Object Names**: Store the `object_name` in your database instead of the full URL, as presigned URLs expire after 7 days.

2. **Regenerate URLs**: When serving images, generate new presigned URLs using the stored `object_name` if needed.

3. **Folder Organization**: Use meaningful folder names:
   - `products/` for product images
   - `users/avatars/` for user avatars
   - `users/covers/` for user cover images
   - `categories/` for category images

4. **Error Handling**: Always handle upload errors gracefully and provide user feedback.

5. **File Validation**: Validate files on the client side before uploading to improve user experience.

6. **Cleanup**: Always delete associated images when deleting entities to prevent orphaned files.

7. **Security**: Consider adding authentication middleware to upload endpoints in production.

## Troubleshooting

### MinIO Connection Issues

If you get connection errors:
1. Verify MinIO is running: `docker ps | grep minio`
2. Check MinIO logs: `docker logs minio`
3. Verify endpoint configuration matches Docker Compose settings

### Bucket Not Found

The bucket is automatically created on first use. If creation fails:
1. Check MinIO credentials
2. Verify network connectivity
3. Check MinIO console at `http://localhost:9001`

### Presigned URL Expiration

Presigned URLs expire after 7 days. To get a new URL:
- Store the `object_name` in your database
- Generate a new presigned URL when needed using the MinIO client

## Additional Resources

- [MinIO Documentation](https://min.io/docs/)
- [Fiber Documentation](https://docs.gofiber.io/)
- [Swagger API Docs](http://localhost:8080/swagger/index.html) (when server is running)

