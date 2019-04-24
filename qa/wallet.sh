#!/bin/bash

set -e

for i in "$@"; do
	case ${i} in
	-n=* | --num=*)
		num="${i#*=}"
		shift
		;;
	--wait)
		wait=yes
		shift
		;;
	*) ;;

	esac
done

[[ -z "$num" ]] && {
	echo "Please specify the number of account peers, e.g., -n=3"
	exit 1
}

declare -a ports

function getPort() {
	echo $((((RANDOM << 15) | RANDOM) % 48128 + 1024))
}

function getApi() {
	echo "http://127.0.0.1:${ports[$1]}"
}

function createPeer() {
	repo="/tmp/peer$1"
	seed="$2"
	api="$3"

	exec > >(sed "s/^/peer$1: /")
	exec 2> >(sed "s/^/peer$1: (stderr) /" >&2)

	rm -rf ${repo}
	mkdir -p ${repo}

	textile init -s $(echo "$wallet" | tail -n1) -a 127.0.0.1:${api} -g 127.0.0.1:$(getPort) -r ${repo} -d -n
	textile daemon -r ${repo} -d
}

trap "exit" INT TERM
trap "kill 0" EXIT

echo "--> Initializing a new wallet..."
echo ""
wallet=$(textile wallet init)
echo "$wallet"
seed=$(echo "$wallet" | tail -n1)

echo ""
echo "--> Creating $num peers from seed $seed"
echo ""

thread=""

for i in $(seq 0 $((num - 1))); do
	if [[ "$i" -eq "0" ]]; then
		ports[i]=40600
	else
		ports[i]=$(getPort)
	fi
	createPeer ${i} ${seed} ${ports[i]} &
	sleep 5

	# bootstrap some content on peer 0
	if [[ "$i" -eq "0" ]]; then
		# create a private thread
		thread=$(textile threads add "test" --api=$(getApi i) | jq ".id")

		# add a message
		textile messages add "ping" -t ${thread} --api=$(getApi i)
	else
		sleep 5
		# add a message to synced thread
		textile messages add "pong" -t ${thread} --api=$(getApi i)
	fi
done

sleep 5

# peer 0's summary
echo ""
echo "--> peer0's summary:"
summary0=$(textile summary --api=$(getApi 0))
apc0=$(echo ${summary0} | jq ".account_peer_count")
tc0=$(echo ${summary0} | jq ".thread_count")
fc0=$(echo ${summary0} | jq ".files_count")
cc0=$(echo ${summary0} | jq ".contact_count")
echo ${summary0} | jq .

# peer 0's account thread
echo ""
echo "--> peer0's account thread:"
account0=$(textile threads get $(textile config Account.Thread --api=$(getApi 0)) --api=$(getApi 0))
ath0=$(echo ${account0} | jq ".head")
atbc0=$(echo ${account0} | jq ".block_count")
atpc0=$(echo ${account0} | jq ".peer_count")
echo ${account0} | jq .

# peer 0's test thread
echo ""
echo "--> peer0's test thread:"
test0=$(textile threads get ${thread} --api=$(getApi 0))
th0=$(echo ${test0} | jq ".head")
tbc0=$(echo ${test0} | jq ".block_count")
tpc0=$(echo ${test0} | jq ".peer_count")
echo ${test0} | jq .

# compare against other peer content
for i in $(seq 1 $((num - 1))); do
	export API=$(getApi i)

	# summary
	summary=$(textile summary)
	apc=$(echo ${summary} | jq ".account_peer_count")
	tc=$(echo ${summary} | jq ".thread_count")
	fc=$(echo ${summary} | jq ".files_count")
	cc=$(echo ${summary} | jq ".contact_count")

	echo ""
	echo "--> Checking peer${i}'s account peer count..."
	if [[ "$apc" -ne "$apc0" ]]; then
		echo ERROR: peer ${i} has incorrect "account_peer_count": ${apc}. Expected ${apc0}.
		exit 1
	else
		echo "--> Got ${apc} ✓"
	fi

	echo ""
	echo "--> Checking peer${i}'s thread count..."
	if [[ "$tc" -ne "$tc0" ]]; then
		echo ERROR: peer ${i} has incorrect "thread_count": ${tc}. Expected ${tc0}.
		exit 1
	else
		echo "--> Got ${tc} ✓"
	fi

	echo ""
	echo "--> Checking peer${i}'s files count..."
	if [[ "$fc" -ne "$fc0" ]]; then
		echo ERROR: peer ${i} has incorrect "files_count": ${fc}. Expected ${fc0}.
		exit 1
	else
		echo "--> Got ${fc} ✓"
	fi

	echo ""
	echo "--> Checking peer${i}'s contact count..."
	if [[ "$cc" -ne "$cc0" ]]; then
		echo ERROR: peer ${i} has incorrect "contact_count": ${cc}. Expected ${cc0}.
		exit 1
	else
		echo "--> Got ${cc} ✓"
	fi

	# account thread
	account=$(textile threads get $(textile config Account.Thread))
	ath=$(echo ${account} | jq ".head")
	atbc=$(echo ${account} | jq ".block_count")
	atpc=$(echo ${account} | jq ".peer_count")

	echo ""
	echo "--> Checking peer${i}'s account thread HEAD..."
	if [[ "$ath" != "$ath0" ]]; then
		echo ERROR: peer ${i} has incorrect "head": ${ath}. Expected ${ath0}.
		exit 1
	else
		echo "--> Got ${ath} ✓"
	fi

	echo ""
	echo "--> Checking peer${i}'s account thread block count..."
	if [[ "$atbc" -ne "$atbc0" ]]; then
		echo ERROR: peer ${i} has incorrect "block_count": ${atbc}. Expected ${atbc0}.
		exit 1
	else
		echo "--> Got ${atbc} ✓"
	fi

	echo ""
	echo "--> Checking peer${i}'s account thread peer count..."
	if [[ "$atpc" -ne "$atpc0" ]]; then
		echo ERROR: peer ${i} has incorrect "peer_count": ${atpc}. Expected ${atpc0}.
		exit 1
	else
		echo "--> Got ${atpc} ✓"
	fi

	# test thread
	test=$(textile threads get ${thread})
	th=$(echo ${test} | jq ".head")
	tbc=$(echo ${test} | jq ".block_count")
	tpc=$(echo ${test} | jq ".peer_count")

	echo ""
	echo "--> Checking peer${i}'s test thread HEAD..."
	if [[ "$th" != "$th0" ]]; then
		echo ERROR: peer ${i} has incorrect "head": ${th}. Expected ${th0}.
		exit 1
	else
		echo "--> Got ${th} ✓"
	fi

	echo ""
	echo "--> Checking peer${i}'s test thread block count..."
	if [[ "$tbc" -ne "$tbc0" ]]; then
		echo ERROR: peer ${i} has incorrect "block_count": ${tbc}. Expected ${tbc0}.
		exit 1
	else
		echo "--> Got ${tbc} ✓"
	fi

	echo ""
	echo "--> Checking peer${i}'s test thread peer count..."
	if [[ "$tpc" -ne "$tpc0" ]]; then
		echo ERROR: peer ${i} has incorrect "peer_count": ${tpc}. Expected ${tpc0}.
		exit 1
	else
		echo "--> Got ${tpc} ✓"
	fi
done

# draw thread graphs
for i in $(seq 1 $((num - 1))); do
	export API=$(getApi i)
	textile blocks ls -t $(textile config Account.Thread) -l 100 -d | dot -Tpng -o /tmp/peer${i}_account_blocks.png
	textile blocks ls -t ${thread} -l 100 -d | dot -Tpng -o /tmp/peer${i}_thread_blocks.png
done

echo "Success."

if [[ "$wait" == "yes" ]]; then
	wait
else
	exit 0
fi
