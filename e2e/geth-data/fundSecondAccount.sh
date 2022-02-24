#!/bin/sh

for i in {1..120} ; do
    echo "Trying to fund second account (attempt #$i)... "
    response=$(geth attach http://localhost:8545 --exec "eth.sendTransaction({from: eth.accounts[0], to: eth.accounts[1], value: 10e36})")
    res=$(echo "$response" | cut -c1-5)
    if [ "$res" != "Fatal" ]; then
        break
    fi
    sleep 0.5
done

echo "Second account funded!"
