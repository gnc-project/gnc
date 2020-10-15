package tests

import (
	"fmt"
	"context"
    "log"
    "github.com/gnc-project/gnc/ethclient"
)

func MinerInfoAt() {
    client, err := ethclient.Dial("http://192.168.1.63:8545")
    if err != nil {
        log.Fatal(err)
	}
	    
    minerIofoResult,err:=client.MinerInfoAt(context.Background())
    if err!=nil{
        log.Fatal(err)
    }
    fmt.Println(minerIofoResult)
}