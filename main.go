package main

import (
    "log"
    "os"

    "github.com/gofiber/fiber/v3"
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

    // Start the server on port 3000
    log.Fatal(app.Listen(":3000"))
}
