package main

import (
    "flashbots-backend/handlers"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "log"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    r := gin.Default()
    r.POST("/send-mev-protected-tx", handlers.SendMEVProtectedTx)

    log.Println("Server running on :8080")
    r.Run(":8080")
}
