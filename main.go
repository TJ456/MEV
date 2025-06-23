package main

import (
    "MEV/handlers"
    "MEV/utils"
    "github.com/gin-gonic/gin"
    "log"
)

func main() {
    // Load environment variables using your utils function
    utils.LoadEnv()

    // Initialize Gin router
    r := gin.Default()

    // Setup POST route for Flashbots MEV-protected transaction
    r.POST("/send-mev-protected-tx", handlers.SendMEVProtectedTx)

    // Start server
    log.Println("🚀 Backend running on http://localhost:8080")
    r.Run(":8080")
}
