# Challenge
Route that takes a file upload, creates a new key, encrypts the file and stores in S3 (or file system)

- **Run:**
```bash
go run main.go
```

- **Upload:**
```bash
curl -X POST -F "file=@image.png" http://localhost:3000/upload
```

- **Download:**
```bash
curl -O "http://localhost:3000/download/:imageId"
```
