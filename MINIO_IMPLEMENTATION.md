# MinIO File Storage Implementation

## Overview

This module implements production-ready file storage using MinIO (S3-compatible object storage) with clean architecture principles, comprehensive validation, and RBAC integration.

## Architecture

The MinIO module follows the same architectural pattern as other modules in the project:

```
internal/minio/
├── client.go      # MinIO client initialization
├── repository.go  # Data access layer for MinIO operations
├── service.go     # Business logic layer
└── validator.go   # File validation logic
```

## Features

✅ **Thread-safe MinIO client** with automatic bucket creation  
✅ **Repository pattern** for clean separation of concerns  
✅ **Comprehensive file validation** (size, MIME type, extension)  
✅ **Scalable folder structure** with timestamps and collision prevention  
✅ **Automatic cleanup** on failed operations  
✅ **RBAC integration** (users can only manage their own images)  
✅ **URL storage optimization** (storing objectName vs full URL)  
✅ **camelCase JSON responses**  

## Configuration

### Environment Variables

Add to `.env`:

```bash
# MinIO Configuration
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=fiber-app
MINIO_USE_SSL=false
```

### Docker Setup

MinIO is included in `docker-compose.yml`:

```yaml
minio:
  image: minio/minio:latest
  ports:
    - "9000:9000"  # API
    - "9001:9001"  # Console
  environment:
    - MINIO_ROOT_USER=minioadmin
    - MINIO_ROOT_PASSWORD=minioadmin
```

**MinIO Console**: <http://localhost:9001>  
**Credentials**: minioadmin / minioadmin

## File Storage Structure

### Best Practice Folder Organization

```
users/{userId}/profile/{timestamp}-{randomHash}.{ext}
```

**Example:**

```
users/550e8400-e29b-41d4/profile/2025-11-28T12:11:45Z-a81fbc2.png
```

### Why This Structure?

1. **Organized by User**: Easy to find all files for a specific user
2. **Timestamp (UTC)**: Chronological ordering and audit trail
3. **Random Hash (7 chars)**: Prevents filename collisions
4. **Lowercase Extension**: Consistent naming convention
5. **Scalable**: Works with millions of users and files

### URL Storage Strategy

**Stored in Database**: Full public URL  
**Example**: `http://localhost:9000/fiber-app/users/92/profile/2025-11-28T12:11:45Z-a81fbc2.png`

**Why Full URL?**

- ✅ Frontend can use directly without transformation
- ✅ CDN-friendly (can be replaced with CDN URL in future)
- ✅ No need to reconstruct URLs in application code
- ✅ Easier debugging and logging

**Alternative Approach (objectName only):**

```go
// Store: "users/92/profile/2025-11-28T12:11:45Z-a81fbc2.png"
// Reconstruct: baseURL + "/" + bucket + "/" + objectName
```

This approach is more flexible if you change storage providers, but requires URL construction logic throughout the application.

## File Validation

### Image Validation Rules

| Rule | Value | Description |
|------|-------|-------------|
| **Max Size** | 5 MB | Prevents server overload |
| **Allowed MIME Types** | `image/jpeg`, `image/png`, `image/webp` | Secure formats only |
| **Allowed Extensions** | `.jpg`, `.jpeg`, `.png`, `.webp` | Prevents disguised files |
| **Filename** | Required | Must not be empty |

### Validation Example

```go
// Validator checks:
1. File size ≤ 5MB
2. Extension in [.jpg, .jpeg, .png, .webp]
3. Content-Type matches allowed MIME types
4. Filename is not empty

// Returns descriptive errors:
{
  "field": "file",
  "message": "file size exceeds maximum allowed size of 5242880 bytes (5.00 MB)"
}
```

## API Endpoints

### 1. Upload Profile Image

```http
POST /api/v1/users/:id/profile-image
Authorization: Bearer <access_token>
Content-Type: multipart/form-data

Form Data:
- image: [file] (required)
```

**RBAC:**

- USER: Can only upload their own image (id must match authenticated user)
- ADMIN: Can upload for any user
- SUPER_USER: Can upload for any user

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Profile image uploaded successfully",
  "data": {
    "imageUrl": "http://localhost:9000/fiber-app/users/92/profile/2025-11-28T12:11:45Z-a81fbc2.png"
  }
}
```

**Error Response (422 Unprocessable Entity):**

```json
{
  "success": false,
  "message": "file size exceeds maximum allowed size of 5242880 bytes (5.00 MB)",
  "error": null
}
```

### 2. Update Profile Image

```http
PUT /api/v1/users/:id/profile-image
Authorization: Bearer <access_token>
Content-Type: multipart/form-data

Form Data:
- image: [file] (required)
```

**Behavior:**

1. Validates new file
2. Deletes old image from MinIO (if exists)
3. Uploads new image
4. Updates database with new URL
5. Rollback on failure (deletes new upload if DB update fails)

**RBAC:** Same as upload

### 3. Delete Profile Image

```http
DELETE /api/v1/users/:id/profile-image
Authorization: Bearer <access_token>
```

**Behavior:**

1. Checks if user has a profile image
2. Deletes from MinIO
3. Sets `profileImageUrl` to empty string in database

**RBAC:** Same as upload

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Profile image deleted successfully",
  "data": null
}
```

**Error Response (404 Not Found):**

```json
{
  "success": false,
  "message": "user has no profile image to delete",
  "error": null
}
```

## Usage Examples

### cURL Examples

#### Upload Image

```bash
curl -X POST http://localhost:8080/api/v1/users/123/profile-image \
  -H "Authorization: Bearer <token>" \
  -F "image=@/path/to/profile.jpg"
```

#### Update Image

```bash
curl -X PUT http://localhost:8080/api/v1/users/123/profile-image \
  -H "Authorization: Bearer <token>" \
  -F "image=@/path/to/new-profile.png"
```

#### Delete Image

```bash
curl -X DELETE http://localhost:8080/api/v1/users/123/profile-image \
  -H "Authorization: Bearer <token>"
```

### JavaScript/TypeScript Example

```typescript
// Upload profile image
async function uploadProfileImage(userId: string, file: File, token: string) {
  const formData = new FormData();
  formData.append('image', file);

  const response = await fetch(`http://localhost:8080/api/v1/users/${userId}/profile-image`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`
    },
    body: formData
  });

  const data = await response.json();
  
  if (data.success) {
    console.log('Image uploaded:', data.data.imageUrl);
    return data.data.imageUrl;
  } else {
    throw new Error(data.message);
  }
}

// Update profile image
async function updateProfileImage(userId: string, file: File, token: string) {
  const formData = new FormData();
  formData.append('image', file);

  const response = await fetch(`http://localhost:8080/api/v1/users/${userId}/profile-image`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`
    },
    body: formData
  });

  return await response.json();
}

// Delete profile image
async function deleteProfileImage(userId: string, token: string) {
  const response = await fetch(`http://localhost:8080/api/v1/users/${userId}/profile-image`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });

  return await response.json();
}
```

### React Component Example

```tsx
import { useState } from 'react';

function ProfileImageUpload({ userId, token }: { userId: string; token: string }) {
  const [uploading, setUploading] = useState(false);
  const [imageUrl, setImageUrl] = useState<string | null>(null);

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    // Validate file size (5MB)
    if (file.size > 5 * 1024 * 1024) {
      alert('File size must be less than 5MB');
      return;
    }

    // Validate file type
    const allowedTypes = ['image/jpeg', 'image/png', 'image/webp'];
    if (!allowedTypes.includes(file.type)) {
      alert('Only JPEG, PNG, and WebP images are allowed');
      return;
    }

    setUploading(true);

    try {
      const formData = new FormData();
      formData.append('image', file);

      const response = await fetch(`/api/v1/users/${userId}/profile-image`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
        },
        body: formData
      });

      const data = await response.json();

      if (data.success) {
        setImageUrl(data.data.imageUrl);
        alert('Image uploaded successfully!');
      } else {
        alert(`Upload failed: ${data.message}`);
      }
    } catch (error) {
      alert('Upload failed');
    } finally {
      setUploading(false);
    }
  };

  return (
    <div>
      <input 
        type="file" 
        accept="image/jpeg,image/png,image/webp"
        onChange={handleFileChange}
        disabled={uploading}
      />
      {uploading && <p>Uploading...</p>}
      {imageUrl && <img src={imageUrl} alt="Profile" width={150} />}
    </div>
  );
}
```

## Extending to Other File Types

The validator is designed to be reusable:

```go
// Custom file type validation
var AllowedDocumentTypes = map[string]bool{
    "application/pdf": true,
    "application/msword": true,
}

var AllowedDocumentExtensions = map[string]bool{
    ".pdf": true,
    ".doc": true,
    ".docx": true,
}

// Use generic validator
err := minio.ValidateFile(
    header, 
    10*1024*1024, // 10MB
    AllowedDocumentTypes,
    AllowedDocumentExtensions,
)
```

## Security Best Practices

### ✅ Implemented

1. **File Type Validation**: Checks both extension and MIME type
2. **Size Limits**: Prevents large file uploads (DoS protection)
3. **RBAC**: Users can only manage their own files
4. **Collision Prevention**: Random hash in filename
5. **Public Read Only**: Bucket policy allows read, not write
6. **Atomic Operations**: Rollback on failure

### 🔒 Additional Recommendations

1. **Virus Scanning**: Integrate ClamAV for malware detection
2. **Image Processing**: Use library to strip EXIF data (privacy)
3. **Rate Limiting**: Limit upload frequency per user
4. **CDN Integration**: Use CloudFront/CloudFlare for better performance
5. **Signed URLs**: For private files, use presigned URLs
6. **Backup Strategy**: Regular MinIO backups

## Error Handling

### Common Errors

| Error | Status | Cause | Solution |
|-------|--------|-------|----------|
| "Image file is required" | 400 | Missing `image` field | Include file in form data |
| "file size exceeds maximum" | 422 | File > 5MB | Compress or resize image |
| "invalid file extension" | 422 | Wrong file type | Use .jpg, .png, or .webp |
| "invalid content type" | 422 | MIME type mismatch | Ensure proper file encoding |
| "You can only upload your own profile image" | 403 | USER accessing other's profile | Use correct user ID |
| "user has no profile image to delete" | 404 | Deleting non-existent image | Check if image exists first |

## Monitoring & Logging

### Recommended Logging

```go
// Log all file operations
log.Printf("File uploaded: user=%s, size=%d, url=%s", userID, fileSize, fileURL)
log.Printf("File deleted: user=%s, url=%s", userID, oldURL)
log.Printf("File update: user=%s, old=%s, new=%s", userID, oldURL, newURL)
```

### Metrics to Track

- Upload success/failure rate
- Average file size
- Storage usage per user
- Upload latency
- Validation failure reasons

## Testing

### Unit Tests

```go
func TestValidateImageFile(t *testing.T) {
    tests := []struct {
        name    string
        size    int64
        mime    string
        ext     string
        wantErr bool
    }{
        {
            name:    "Valid JPEG",
            size:    1024 * 1024, // 1MB
            mime:    "image/jpeg",
            ext:     ".jpg",
            wantErr: false,
        },
        {
            name:    "File too large",
            size:    6 * 1024 * 1024, // 6MB
            mime:    "image/jpeg",
            ext:     ".jpg",
            wantErr: true,
        },
        {
            name:    "Invalid extension",
            size:    1024 * 1024,
            mime:    "image/jpeg",
            ext:     ".bmp",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Production Checklist

- [ ] MinIO configured with strong credentials
- [ ] SSL/TLS enabled for production (`MINIO_USE_SSL=true`)
- [ ] Bucket policy correctly set
- [ ] File size limits appropriate for use case
- [ ] RBAC rules tested for all user types
- [ ] Error handling covers all edge cases
- [ ] Logging configured for audit trail
- [ ] Backup strategy in place
- [ ] CDN configured (optional but recommended)
- [ ] Rate limiting enabled

## Troubleshooting

### Issue: "Failed to create MinIO client"

**Solution**: Check MinIO endpoint and credentials in `.env`

### Issue: "Failed to upload file to MinIO"

**Solution**: Verify bucket exists and credentials have write permission

### Issue: "Object not found" when deleting

**Solution**: File may have been manually deleted from MinIO console

### Issue: Database update fails after upload

**Solution**: File is automatically deleted from MinIO (rollback), check DB connection

## Resources

- [MinIO Go SDK Documentation](https://min.io/docs/minio/linux/developers/go/minio-go.html)
- [S3 Best Practices](https://docs.aws.amazon.com/AmazonS3/latest/userguide/optimizing-performance.html)
- [Image Optimization Guide](https://web.dev/fast/#optimize-your-images)
