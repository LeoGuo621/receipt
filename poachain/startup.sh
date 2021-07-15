gnome-terminal --tab -t "node1" -- bash -c "sh ./node1.sh;exec bash" 
sleep 1

gnome-terminal --tab -t "node2" -- bash -c "sh ./node2.sh;exec bash"
sleep 1

gnome-terminal --tab -t "node3" -- bash -c "sh ./node3.sh;exec bash"
echo "All nodes has been initialized (number of nodes == 3)."
sleep 10
gnome-terminal --tab -t "networking" -- bash -c "sh ./networking.sh;exec bash"
