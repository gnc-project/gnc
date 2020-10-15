package tests

import (
    "context"
    "crypto/ecdsa"
    "fmt"
    "log"
    "math/big"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    // client, err := ethclient.Dial("http://47.112.224.137:8545")
    client, err := ethclient.Dial("http://192.168.1.63:8545")
    // client, err := ethclient.Dial("http://localhost:8545")
    if err != nil {
        log.Fatal(err)
    }
    var data []byte

    data=[]byte("")
    var txType types.TxType

    txType=1
//0xf5403E4F120901407eF221E2419583D1F3556953
    // privateKey, err := crypto.HexToECDSA("c6bacf5a75ed8c46ce49723974d77912c75f23a062e5f9c6a0937778af27323a")
    // privateKey, err := crypto.HexToECDSA("0978f416fe608a0f530b448e79477c202bcf40605fc19e921f191b17607c9fd0")
// "0xa086ff707591f0b3a051aa31a3198a74d22d9c4f"
    // privateKey, err := crypto.HexToECDSA("460547f84edf49d2cac7aab46852deb3e1ab6385a1fca0329d4803d51d442030")
// wl_793d13778effe395edb9d4255ddff69b73c72971
    privateKey, err := crypto.HexToECDSA("FA6E7EC9827F5E332725589D3612ECD803164469C95147FE343A85D54E501AD3")
    // privateKey, err := crypto.HexToECDSA("cfcfa295cab51ccae9110ed6932c2c68dc0b94dba300baf5b8890906b248b50b")
    if err != nil {
        log.Fatal(err)
    }
    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
    }

    fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
    
    if err != nil {
        log.Fatal(err)
    }
    // fmt.Println(fromAddress)
    value,_:= new(big.Int).SetString("1",10)
//     // gasPrice, err := client.SuggestGasPrice(context.Background())
//     // if err != nil {
//     //     log.Fatal(err)
//     // }
    toAddress := common.HexToAddress("GNC91e1cb367d192b049c4c36732c09004fc3371d0a")
    fmt.Println(toAddress.Hex())
    
    var signedTx *types.Transaction
    for i:=0;i<1;i++{
        nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
        fmt.Println(nonce)
        tx := types.NewTransaction(txType,nonce, toAddress, value,21000, big.NewInt(81000), data)
        chainID, err := client.NetworkID(context.Background())
        if err != nil {
            log.Fatal(err)
        }
        signedTx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
        if err != nil {
            log.Fatal(err)    
        }
        
        // fmt.Println(signedTx)
        err = client.SendTransaction(context.Background(), signedTx)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("tx sent: %s\n", signedTx.Hash().Hex()) 
    }
	for {
    tx, isPending, err := client.TransactionByHash(context.Background(), signedTx.Hash())
    if err != nil {
        log.Fatal(err)
    }
    if isPending==false{
         fmt.Println("transaction is successful!!")
		 receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
        }
        if receipt.Status==0{
            log.Fatal( "Error: Transaction has been reverted by the EVM")
        }
		fmt.Printf("receipt.Status:%v\n",receipt.Status)
		return 
    }
   }
   

}
