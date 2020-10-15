package tests


import(
	"fmt"
	"context"
    "log"
    "math/big"
    "github.com/ethereum/go-ethereum/ethclient"
)


func AddNonceAt(){
	client, err := ethclient.Dial("http://192.168.1.63:8545")
	pid,_:=new(big.Int).SetString("10358737893495713720",10)
    GeneratorState,err:= client.AddNonceAt(context.Background(),pid,big.NewInt(21100100),big.NewInt(1))
    if err!=nil{
        log.Fatal(err)
	}
	
	fmt.Println(GeneratorState)
}