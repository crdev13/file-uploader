package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "errors"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"

    "github.com/gofiber/fiber/v3"
    "github.com/google/uuid"
)

var encryptionKey = []byte("a very very very very secret key") // 32 bytes

var storeName = "uploads"

var allowedTypes = map[string]bool{
    "image/jpeg": true,
    "image/png":  true,
    "text/plain": false,
}

func main() {
    os.MkdirAll(storeName, os.ModePerm)

    app := fiber.New(
        fiber.Config{
            BodyLimit: 10 * 1024 * 1024, // 10MB in bytes
        },
    )

    app.Get("/", func(ctx fiber.Ctx) error {
        return ctx.SendString("Hello, World 👋!")
    })

    app.Post("/upload", func(ctx fiber.Ctx) error {
        file, err := ctx.FormFile("file")
        if err != nil {
            return err
        }

        allowedType := allowedTypes[file.Header.Get("Content-Type")]
        if !allowedType {
            return errors.New("invalid file type")
        }

        data, err := file.Open()
        if err != nil {
            return err
        }

        fileName := uuid.NewString()
        if err := saveFile(fileName, data); err != nil {
            return err
        }

        return ctx.SendString(fmt.Sprintf("uploaded: %s", fileName))
    })

    // Start the server on port 3000
    log.Fatal(app.Listen(":3000"))
}

func saveFile(fileName string, file io.Reader) error {
    data, err := io.ReadAll(file)
    if err != nil {
        return errors.New("failed to read file")
    }

    encryptedData, err := encrypt(data)
    if err != nil {
        return errors.New("failed to encrypt file")
    }

    dst, err := os.Create(filepath.Join(storeName, fileName))
    if err != nil {
        return errors.New("unable to create file on server")
    }
    defer dst.Close()

    _, err = dst.Write(encryptedData)
    if err != nil {
        return errors.New("unable to save file")
    }
    return nil
}

func encrypt(data []byte) ([]byte, error) {
    block, err := aes.NewCipher(encryptionKey)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    return gcm.Seal(nonce, nonce, data, nil), nil
}
