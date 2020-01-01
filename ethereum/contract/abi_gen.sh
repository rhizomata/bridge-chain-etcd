contractpath=$GOPATH/github.com/rhizomata/bridge-chain-etcd/ethereum/contract
echo $contractpath
docker run -it --rm -v $contractpath:/sources \
  ethereum/solc:0.4.20 bash