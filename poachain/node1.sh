geth --datadir node1/data --networkid 621 init poachain.json 2>> node1.log
cp -rf node1keystore/keystore node1/data
echo "node1 has been initialized"
sleep 1
geth --datadir ./node1/data --networkid 621 --port 30303 --ws --ws.port 8546 --http --http.port 8545 --http.corsdomain '*' --http.api 'web3,eth,debug,personal,net' --allow-insecure-unlock --rpc.allow-unprotected-txs 2>> node1.log
echo "node1 has started. Connected to network (621)"
