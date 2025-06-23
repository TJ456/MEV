package handlers

import (
    "context"
    "crypto/ecdsa"
    "encoding/json"
    "flashbots-backend/utils"
    "log"
    "math/big"
    "net/http"
    "os"

    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/gin-gonic/gin"
)

type TxRequest struct {
    To       string `json:"to"`
    Value    string `json:"value"`
    GasLimit uint64 `json:"gasLimit"`
    Data     string `json:"data"`
}

func SendMEVProtectedTx(c *gin.Context) {
    var req TxRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    client, err := utils.GetClient()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "RPC connection error"})
        return
    }

    key, err := crypto.HexToECDSA(os.Getenv("FLASHBOTS_RELAYER_KEY"))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid relayer key"})
        return
    }

    fromAddress := crypto.PubkeyToAddress(key.PublicKey)
    nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get nonce"})
        return
    }

    value := new(big.Int)
    value.SetString(req.Value, 10)

    gasPrice, _ := client.SuggestGasPrice(context.Background())
    toAddress := common.HexToAddress(req.To)
    tx := types.NewTransaction(nonce, toAddress, value, req.GasLimit, gasPrice, common.FromHex(req.Data))

    chainID, _ := client.NetworkID(context.Background())
    signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), key)
    err = client.SendTransaction(context.Background(), signedTx)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "Transaction sent via backend",
        "txHash": signedTx.Hash().Hex(),
    })
}
