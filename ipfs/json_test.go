package ipfs_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/gateway"
	. "github.com/textileio/go-textile/ipfs"
)

var vars = struct {
	nodePath string
	node     *core.Textile

	input1 string
	input2 string
	input3 string
}{
	nodePath: "./testdata/.textile",

	input1: `
{
    "_id": "5d28de69de4a0a091097e507",
    "index": 0,
    "guid": "5058db26-6318-448b-98ba-287fa52104be",
    "isActive": false,
    "balance": "$1,086.88",
    "picture": "http://placehold.it/32x32",
    "age": 38,
    "eyeColor": "brown",
    "name": {
        "first": "Liza",
        "last": "Collins"
    },
    "company": "EURON",
    "email": "liza.collins@euron.me",
    "phone": "+1 (950) 451-3373",
    "address": "129 Cranberry Street, Otranto, Kentucky, 7502",
    "about": "Laborum irure culpa cillum do amet proident. Do anim aute aute labore. Aute laboris laboris nostrud velit minim nulla aliqua quis. Non aute ex et sit fugiat Lorem ut quis reprehenderit nisi qui in ut. Nisi incididunt magna nulla minim aliquip mollit reprehenderit adipisicing occaecat.",
    "registered": "Tuesday, July 28, 2015 8:57 AM",
    "latitude": "-30.373386",
    "longitude": "-137.985157",
    "friends": [
        {
            "id": 0,
            "name": "Cornelia Delacruz"
        },
        {
            "id": 1,
            "name": "Navarro Hawkins"
        },
        {
            "id": 2,
            "name": "Nadia Finley"
        }
    ],
    "greeting": "Hello, Liza! You have 7 unread messages.",
    "favoriteFruit": "banana"
}
	`,
	input2: `
[
    {
        "_id": "5d28de69de4a0a091097e507",
        "index": 0,
        "guid": "5058db26-6318-448b-98ba-287fa52104be",
        "isActive": false,
        "balance": "$1,086.88",
        "picture": "http://placehold.it/32x32",
        "age": 38,
        "eyeColor": "brown",
        "name": {
            "first": "Liza",
            "last": "Collins"
        },
        "company": "EURON",
        "email": "liza.collins@euron.me",
        "phone": "+1 (950) 451-3373",
        "address": "129 Cranberry Street, Otranto, Kentucky, 7502",
        "about": "Laborum irure culpa cillum do amet proident. Do anim aute aute labore. Aute laboris laboris nostrud velit minim nulla aliqua quis. Non aute ex et sit fugiat Lorem ut quis reprehenderit nisi qui in ut. Nisi incididunt magna nulla minim aliquip mollit reprehenderit adipisicing occaecat.",
        "registered": "Tuesday, July 28, 2015 8:57 AM",
        "latitude": "-30.373386",
        "longitude": "-137.985157",
        "tags": [
            "dolore",
            "labore",
            "sunt",
            "qui",
            "cillum"
        ],
        "range": [
            0,
            1,
            2,
            3,
            4,
            5,
            6,
            7,
            8,
            9
        ],
        "friends": [
            {
                "id": 0,
                "name": "Cornelia Delacruz"
            },
            {
                "id": 1,
                "name": "Navarro Hawkins"
            },
            {
                "id": 2,
                "name": "Nadia Finley"
            }
        ],
        "greeting": "Hello, Liza! You have 7 unread messages.",
        "favoriteFruit": "banana"
    },
    {
        "_id": "5d28de69ab99c3483a8309bb",
        "index": 1,
        "guid": "9254d657-94af-4f64-a3a7-16b9a46b8f95",
        "isActive": true,
        "balance": "$1,128.81",
        "picture": "http://placehold.it/32x32",
        "age": 37,
        "eyeColor": "green",
        "name": {
            "first": "Tracey",
            "last": "Osborn"
        },
        "company": "RECRITUBE",
        "email": "tracey.osborn@recritube.com",
        "phone": "+1 (848) 405-3368",
        "address": "560 Hopkins Street, Loomis, Rhode Island, 7543",
        "about": "Voluptate occaecat adipisicing officia sit. Aute commodo ipsum id minim ipsum minim mollit aliquip ut do nisi labore. Laboris excepteur amet aliquip enim non minim officia est ea. Est labore aliquip fugiat dolore aute cillum Lorem esse id culpa. Nisi aliquip occaecat excepteur ut ad anim ex ullamco. Incididunt deserunt voluptate deserunt aliqua non anim veniam ipsum aliquip fugiat in officia eu deserunt.",
        "registered": "Monday, August 20, 2018 6:39 AM",
        "latitude": "-89.638045",
        "longitude": "-111.183917",
        "tags": [
            "velit",
            "irure",
            "non",
            "qui",
            "ea"
        ],
        "range": [
            0,
            1,
            2,
            3,
            4,
            5,
            6,
            7,
            8,
            9
        ],
        "friends": [
            {
                "id": 0,
                "name": "Angela Taylor"
            },
            {
                "id": 1,
                "name": "Kristen Cortez"
            },
            {
                "id": 2,
                "name": "Cleo Stone"
            }
        ],
        "greeting": "Hello, Tracey! You have 8 unread messages.",
        "favoriteFruit": "banana"
    }
]
	`,
	input3: `
{
    "default": 7,
    "accounts": {
        "P5PiuxRn7qiYM2Wgdjzzfihvo2Stgow3XQ8HyCqM1Xukr4Rb": "WRITE"
    }
}
	`,
}

func TestIPFS_Setup(t *testing.T) {
	var err error
	vars.node, err = core.CreateAndStartPeer(core.InitConfig{
		RepoPath: vars.nodePath,
		Debug:    true,
	}, true)
	if err != nil {
		t.Fatal(err)
	}

	gateway.Host = &gateway.Gateway{
		Node: vars.node,
	}
	gateway.Host.Start(vars.node.Config().Addresses.Gateway)
}

func TestIPFS_AddJSON_Object(t *testing.T) {
	cid, err := AddJSON(vars.node.Ipfs(), vars.input1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cid.String())

	j, err := JSONAtPath(vars.node.Ipfs(), cid.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(j))

	if bytes.Equal(jsonBytes(string(j)), jsonBytes(vars.input1)) {
		t.Fatal("output does not equal input")
	}
}

func TestIPFS_AddJSON_Array(t *testing.T) {
	cid, err := AddJSON(vars.node.Ipfs(), vars.input2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cid.String())

	j, err := JSONAtPath(vars.node.Ipfs(), cid.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(j))

	if bytes.Equal(jsonBytes(string(j)), jsonBytes(vars.input2)) {
		t.Fatal("output does not equal input")
	}
}

func TestIPFS_AddJSON_Roles(t *testing.T) {
	cid, err := AddJSON(vars.node.Ipfs(), vars.input3)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cid.String())

	j, err := JSONAtPath(vars.node.Ipfs(), cid.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(j))

	if bytes.Equal(jsonBytes(string(j)), jsonBytes(vars.input3)) {
		t.Fatal("output does not equal input")
	}
}

func TestIPFS_Teardown(t *testing.T) {
	_ = vars.node.Stop()
	vars.node = nil
}

func jsonBytes(input string) []byte {
	var i interface{}
	_ = json.Unmarshal([]byte(input), &i)
	b, _ := json.Marshal(input)
	return b
}
