#!/bin/sh
if [ "$1" = "docker" ]; then
    docker="docker exec hubble-geth"
fi

i=1
while [ "$i" -le 120 ]; do
    echo "Trying to fund second account (attempt #$i)... "
    response=$($docker geth attach --exec "eth.sendTransaction({from: eth.accounts[0], to: eth.accounts[1], value: 10e36}) http://localhost:8545")
    res=$(echo "$response" | cut -c1-3)
    if [ "$res" = "\"0x" ]; then
        break
    fi
    sleep 0.5

    i=$(expr $i + 1)
done

echo "Second account funded!"
