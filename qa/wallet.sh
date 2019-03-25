#!/bin/bash

set -e

while getopts n: option
do
case "${option}"
in
n) num=${OPTARG};;
esac
done

[[ -z "$num" ]] && { echo "Please specify the number of account peers, e.g., -n 3" ; exit 1; }

function getPort(){
    echo $(( ((RANDOM<<15)|RANDOM) % 48128 + 1024 ))
}

function prepRepo() {
    repo="$1"

    rm -rf ${repo}
    mkdir -p ${repo}/logs/
    touch ${repo}/logs/textile.log
}

function createPeer() {
    repo="$1"
    seed="$2"
    api="$3"

    rm -rf ${repo}
    mkdir -p ${repo}/logs/
    touch ${repo}/logs/textile.log

    textile init -s $(echo "$wallet" | tail -n1) -a 127.0.0.1:${api} -g 127.0.0.1:$(getPort) -r ${repo} -d
    textile daemon -r ${repo} -d
}

echo "--> Initializing a new wallet..."
echo ""
wallet=$(textile wallet init)
echo "$wallet"
seed=$(echo "$wallet" | tail -n1)

echo ""
echo "--> Creating $num peers from seed $seed"
echo ""

declare -a ports

logs=""
for i in `seq 0 $((num - 1))`;
do
    prepRepo "/tmp/peer$i"

    log="/tmp/peer$i/logs/textile.log "
    logs=${logs}${log}
done

for i in `seq 0 $((num - 1))`;
do
    ports[i]=$(getPort)
    createPeer "/tmp/peer$i" ${seed} ${ports[i]} &
    sleep 5 # give a little breathing room
done

tail -f ${logs} &



wait
