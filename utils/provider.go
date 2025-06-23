package utils

import (
    "github.com/ethereum/go-ethereum/ethclient"
    "os"
)

func GetClient() (*ethclient.Client, error) {
    rpcUrl := os.Getenv("ETHEREUM_RPC_URL")
    return ethclient.Dial(rpcUrl)
}
