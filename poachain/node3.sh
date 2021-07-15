geth --datadir node3/data --networkid 621 init poachain.json 2>> node3.log
cp -rf node3keystore/keystore node3/data
echo "node2 has been initialized"
sleep 1
geth --datadir ./node3/data --networkid 621 --port 30323 --ws --ws.port 8566 --http --http.port 8565 --http.corsdomain '*' --http.api 'web3,eth,debug,personal,net' --allow-insecure-unlock --rpc.allow-unprotected-txs 2>> node3.log
echo "node3 has started. Connected to network (621)"
