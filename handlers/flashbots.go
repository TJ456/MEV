package handlers

import (
    "context"
    "MEV/utils"
    "log"
    "math/big"
    "net/http"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/gin-gonic/gin"
)

type TxRequest struct {
    To       string `json:"to"`
    Value    string `json:"value"`    // In wei (as string)
    GasLimit uint64 `json:"gasLimit"` // Example: 21000
    Data     string `json:"data"`     // Optional data for contract call
}

func SendMEVProtectedTx(c *gin.Context) {
    var req TxRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request"})
        return
    }

    client, err := utils.GetClient()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "RPC connection error"})
        return
    }

    key, err := crypto.HexToECDSA(utils.GetEnv("FLASHBOTS_RELAYER_KEY"))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid relayer private key"})
        return
    }

    fromAddress := crypto.PubkeyToAddress(key.PublicKey)
    nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get nonce"})
        return
    }

    // Convert value (wei string) to big.Int
value := new(big.Int)
_, ok := value.SetString(req.Value, 10)
if !ok {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value format"})
    return
}


    gasPrice, err := client.SuggestGasPrice(context.Background())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get gas price"})
        return
    }

    toAddress := common.HexToAddress(req.To)

    // Validate data
    txData := common.FromHex(req.Data)

    tx := types.NewTransaction(nonce, toAddress, value, req.GasLimit, gasPrice, txData)

    chainID, err := client.NetworkID(context.Background())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get network ID"})
        return
    }

    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), key)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction signing failed"})
        return
    }

    err = client.SendTransaction(context.Background(), signedTx)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    log.Println("✅ Transaction sent:", signedTx.Hash().Hex())

    c.JSON(http.StatusOK, gin.H{
        "status":    "Transaction sent via backend",
        "txHash":    signedTx.Hash().Hex(),
        "to":        req.To,
        "from":      fromAddress.Hex(),
        "value_wei": req.Value,
        "gas_price": gasPrice.String(),
        "gas_limit": req.GasLimit,
    })
}
