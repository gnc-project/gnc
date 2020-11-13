## GNC

```txt
You can download the source code and compile it yourself, or download the binary code file directly from https://github.com/gnc-project/gnc-node
```

### how to compile

mkdir /home/gnccode

cd /home/gnccode

download : git clone https://github.com/gnc-project/gnc

1: cd /home/gnccode/gnc && make all

2: cd /home/gnccode/gnc/bulid/bin && mv gnc /home/gncnode

4: cd /home/gncnode && mkdir data

4: nohup ./gnc  --syncmode 'full' --cache 256 --datadir data --rpc  --rpcaddr 0.0.0.0 --rpcapi 'net,db,eth,web3,txpool' --rpcport 8545 --port 30303 --rpccorsdomain "*"  --rpcvhosts "*" --allow-insecure-unlock --networkid 37021 --ws --wsaddr 0.0.0.0 --wsport "8560" --wsapi "net,db,eth,web3,txpool" --wsorigins "*" >data/gnc.log 2>&1 &
