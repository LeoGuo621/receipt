geth --datadir node2/data --networkid 621 init poachain.json 2>> node2.log
cp -rf node2keystore/keystore node2/data
echo "node2 has been initialized"
sleep 1
geth --datadir ./node2/data --networkid 621 --port 30313 --ws --ws.port 8556 --http --http.port 8555 --http.corsdomain '*' --http.api 'web3,eth,debug,personal,net' --allow-insecure-unlock --rpc.allow-unprotected-txs 2>> node2.log
echo "node2 has started. Connected to network (621)"
