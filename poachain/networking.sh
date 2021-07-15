echo "get the enode address of all nodes"
node1_enode_addr=$(geth attach "rpc:http://localhost:8545" << EOF | grep "Data: " | sed "s/Data: //")
var addr = admin.nodeInfo.enode;
console.log("Data: " + addr);
exit;
EOF
echo "node1 enode address: ${node1_enode_addr}"

node2_enode_addr=$(geth attach "rpc:http://localhost:8555" << EOF | grep "Data: " | sed "s/Data: //")
var addr = admin.nodeInfo.enode;
console.log("Data: " + addr);
exit;
EOF
echo "node2 enode address: ${node2_enode_addr}"

node3_enode_addr=$(geth attach "rpc:http://localhost:8565" << EOF | grep "Data: " | sed "s/Data: //")
var addr = admin.nodeInfo.enode;
console.log("Data: " + addr);
exit;
EOF
echo "node3 enode address: ${node3_enode_addr}"

# Networking
echo "Networking ..."

conn_status_13=$(geth attach "rpc:http://localhost:8545" << EOF | grep "Data: " | sed "s/Data: //")
var connstatus = admin.addPeer("$node3_enode_addr");
console.log("Data: " + connstatus);
exit;
EOF
echo "node1 has connected to node3 with status=${conn_status_13}"

conn_status_21=$(geth attach "rpc:http://localhost:8555" << EOF | grep "Data: " | sed "s/Data: //")
var connstatus = admin.addPeer("$node1_enode_addr");
console.log("Data: " + connstatus);
exit;
EOF
echo "node2 has connected to node1 with status=${conn_status_21}"

conn_status_32=$(geth attach "rpc:http://localhost:8565" << EOF | grep "Data: " | sed "s/Data: //")
var connstatus = admin.addPeer("$node2_enode_addr");
console.log("Data: " + connstatus);
exit;
EOF
echo "node3 has connected to node2 with status=${conn_status_32}"
echo "all nodes have connected to the same network(621)."

peerInfo_1=$(geth attach "rpc:http://localhost:8545" << EOF | grep "Data: " | sed "s/Data: //")
var peerinfo = net.peerCount;
console.log("Data: " + peerinfo);
exit;
EOF
echo "Checking connection on node1: ${peerInfo_1} < should see 2"

peerInfo_2=$(geth attach "rpc:http://localhost:8555" << EOF | grep "Data: " | sed "s/Data: //")
var peerinfo = net.peerCount;
console.log("Data: " + peerinfo);
exit;
EOF
echo "Checking connection on node2: ${peerInfo_2} < should see 2"

peerInfo_3=$(geth attach "rpc:http://localhost:8565" << EOF | grep "Data: " | sed "s/Data: //")
var peerinfo = net.peerCount;
console.log("Data: " + peerinfo);
exit;
EOF
echo "Checking connection on node3: ${peerInfo_3} < should see 2"

sleep 1
# Start mining
sealer1_stat=$(geth attach "rpc:http://localhost:8545" << EOF | grep "Data: " | sed "s/Data: //")
personal.unlockAccount(eth.accounts[0], "123", 9999999)
miner.setEtherbase(eth.coinbase);
var stat = miner.start(1);
console.log("Data: " + stat);
exit;
EOF
echo "sealer1 mining stat: ${sealer1_stat} < should see null"

sealer2_stat=$(geth attach "rpc:http://localhost:8555" << EOF | grep "Data: " | sed "s/Data: //")
personal.unlockAccount(eth.accounts[0], "123", 9999999)
miner.setEtherbase(eth.coinbase);
var stat = miner.start(1);
console.log("Data: " + stat);
exit;
EOF
echo "sealer2 mining stat: ${sealer2_stat} < should see null"

sealer3_stat=$(geth attach "rpc:http://localhost:8565" << EOF | grep "Data: " | sed "s/Data: //")
personal.unlockAccount(eth.accounts[0], "123", 9999999)
miner.setEtherbase(eth.coinbase);
var stat = miner.start(1);
console.log("Data: " + stat);
exit;
EOF
echo "sealer3 mining stat: ${sealer3_stat} < should see null"
echo "3 sealers have started mining."
