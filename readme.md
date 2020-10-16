## GNC

```txt
You can download the source code and compile it yourself, or download the binary code file directly from https://github.com/gnc-project/gnc-node
```

### how to compile

mkdir /home/gnccode

cd /home/gnccode

download : git clone https://github.com/gnc-project/gnc

1: cd /home/gnccode/gnc && make all

2: mkdir /home/gncnode 

3: cd /home/gnccode/gnc/bulid/bin && mv gnc /home/gncnode

4: cd /home/gnccode/gnc && mv gnc.json /home/gncnode

4: cd /home/gncnode && ./gnc --datadir data init ./gnc.json

4: nohup ./gnc  --syncmode 'full' --cache 256 --datadir data --rpc  --rpcaddr 0.0.0.0 --rpcapi 'net,db,eth,web3,txpool' --rpcport 8545 --port 30303 --rpccorsdomain "*"  --rpcvhosts "*" --allow-insecure-unlock --networkid 37021 --ws --wsaddr 0.0.0.0 --wsport "8560" --wsapi "net,db,eth,web3,txpool" --bootnodes "enode://d01987c09dc1149c7da115c8e9335b531716510fce9a0dec74d226a8a6d582cbbf8e455400e5c446359bee7e92341a0a62f69004474b20e9fcac302bed478c32@47.75.203.235:30303" --wsorigins "*" >data/gnc.log 2>&1 &

5: If the node connection or Sync fails, you can add this node

```txt
"enode://7ec759c185382e169e3fbd4718ca9907f26897109a5401aadb6128c26bce7bc9bc26613d202e67db6019137bec412d87a1f1411877cd5e165a544515c57203c9@47.57.116.216:30303"

"enode://b09b539b1bf9b0bb4545fe89b969e440361a6c790809f8769d1bc91e4d631462e21be17db6db9f1df67f0380cc23ff39e32843496ca501e7d5022ac378ec7484@47.57.115.222:30303",

"enode://d01987c09dc1149c7da115c8e9335b531716510fce9a0dec74d226a8a6d582cbbf8e455400e5c446359bee7e92341a0a62f69004474b20e9fcac302bed478c32@47.75.203.235:30303"

"enode://10b45f2f3d3f27a77d19c74caabb24027814beea634ca540828ced46cf726b0808c59a3aafe7d76955fe55b5d4a6fb749b3611d27529fb0802a7ff0bccbac505@47.115.113.25:30303"

```
