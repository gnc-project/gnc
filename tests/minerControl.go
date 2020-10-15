package tests


import (
    "fmt"
	"os"
	"context"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/common"
)
func test(){

	if len(os.Args) != 3 {
		fmt.Println("Args insufficient")
		os.Exit(1)
	}
	  client, err := ethclient.Dial(os.Args[1])

    if err != nil {
        fmt.Println("Failed to connect to the Ethereum client:", err)
        os.Exit(0)
	}
	if (os.Args[2]=="start"){
		result,err:=client.MinerStart(context.Background())
		if err!=nil{
			panic(err)
		}
		fmt.Println("minerStart",result)
	}
	if (os.Args[2]=="stop"){
		result,err:=client.MinerStop(context.Background())
		if err!=nil{
			panic(err)
		}
		fmt.Println("minerStop",result)
	}
	if common.IsHexAddress(os.Args[2]){
		boolen,err:=client.SuperNodeAt(context.Background(),common.HexToAddress(os.Args[2]),"latest")
		if err!=nil{
			panic(err)
		}
		if boolen||common.HexToAddress(os.Args[2])==common.HexToAddress("793d13778effe395edb9d4255ddff69b73c72971"){
			result,err:=client.SetEtherBase(context.Background(),common.HexToAddress(os.Args[2]))
			if err!=nil{
				panic(err)
			}
			fmt.Printf("SetEtherBase:%v,coinbase=%v\n",result,os.Args[2])
		}else{
			fmt.Println("coinbase is not a superNode")
			return
		}
		
	}

	return
	
}