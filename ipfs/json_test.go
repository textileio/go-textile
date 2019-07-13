package ipfs_test

import (
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
    },
    {
        "_id": "5d28de69181cc13f4098b43c",
        "index": 2,
        "guid": "72d3b854-ad3b-48a0-8586-933bcf4db91b",
        "isActive": true,
        "balance": "$1,880.80",
        "picture": "http://placehold.it/32x32",
        "age": 39,
        "eyeColor": "brown",
        "name": {
            "first": "Booker",
            "last": "Clements"
        },
        "company": "ANARCO",
        "email": "booker.clements@anarco.biz",
        "phone": "+1 (948) 427-3054",
        "address": "137 Ingraham Street, Fingerville, Pennsylvania, 5567",
        "about": "Culpa eiusmod mollit ad occaecat sunt in dolore velit est proident excepteur. Ut ullamco nostrud incididunt minim non nulla laborum ullamco reprehenderit. Esse labore ea non ad nulla aliquip ut nulla officia. Amet velit ea sit consectetur in.",
        "registered": "Monday, June 4, 2018 2:13 AM",
        "latitude": "-36.779417",
        "longitude": "-109.215102",
        "tags": [
            "magna",
            "amet",
            "id",
            "sit",
            "occaecat"
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
                "name": "Terri Nicholson"
            },
            {
                "id": 1,
                "name": "Eaton Moreno"
            },
            {
                "id": 2,
                "name": "Ola Byrd"
            }
        ],
        "greeting": "Hello, Booker! You have 10 unread messages.",
        "favoriteFruit": "banana"
    },
    {
        "_id": "5d28de696a60f30f0c35440a",
        "index": 3,
        "guid": "053914d8-0f07-457c-a0b0-237b143f7103",
        "isActive": false,
        "balance": "$2,379.77",
        "picture": "http://placehold.it/32x32",
        "age": 23,
        "eyeColor": "green",
        "name": {
            "first": "Maribel",
            "last": "Lynch"
        },
        "company": "QUALITERN",
        "email": "maribel.lynch@qualitern.net",
        "phone": "+1 (996) 598-2639",
        "address": "194 Jamaica Avenue, Foxworth, New Hampshire, 1097",
        "about": "Do tempor excepteur velit in ex. Deserunt amet officia mollit incididunt laborum in tempor ullamco dolor proident excepteur. Nulla ut cillum in aliqua mollit amet cupidatat elit sint aliquip labore Lorem. Culpa enim ad Lorem ipsum aliquip. Cupidatat aliquip officia eiusmod mollit laboris nulla labore duis velit enim irure occaecat. Duis incididunt adipisicing elit nulla. Aliquip incididunt qui ullamco cillum Lorem ipsum esse.",
        "registered": "Thursday, August 6, 2015 3:44 AM",
        "latitude": "-25.139304",
        "longitude": "156.595556",
        "tags": [
            "dolor",
            "veniam",
            "laboris",
            "sit",
            "ex"
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
                "name": "Noreen Lyons"
            },
            {
                "id": 1,
                "name": "Colon Walsh"
            },
            {
                "id": 2,
                "name": "Lessie Donovan"
            }
        ],
        "greeting": "Hello, Maribel! You have 6 unread messages.",
        "favoriteFruit": "banana"
    },
    {
        "_id": "5d28de69ea4ed4ee3a2256b3",
        "index": 4,
        "guid": "8506fff0-cf24-49d7-9bbf-9ea07b25f96e",
        "isActive": false,
        "balance": "$3,758.72",
        "picture": "http://placehold.it/32x32",
        "age": 35,
        "eyeColor": "brown",
        "name": {
            "first": "Estela",
            "last": "Hopkins"
        },
        "company": "GEEKMOSIS",
        "email": "estela.hopkins@geekmosis.tv",
        "phone": "+1 (952) 537-2497",
        "address": "947 Belvidere Street, Watchtower, Arkansas, 3323",
        "about": "Enim commodo laboris eu deserunt ut enim eu velit veniam id ullamco aliquip labore. Ad ex consectetur voluptate nulla non esse commodo velit nostrud magna eiusmod labore excepteur Lorem. Aliqua laboris nisi labore irure laboris laboris excepteur incididunt fugiat laboris ad aute enim. Veniam culpa esse culpa laboris quis nisi consequat exercitation tempor nisi ullamco et aliqua eu. Occaecat culpa magna sunt excepteur ipsum labore culpa enim exercitation ut laboris consequat aliqua.",
        "registered": "Sunday, May 10, 2015 5:00 AM",
        "latitude": "-39.669387",
        "longitude": "-25.634747",
        "tags": [
            "irure",
            "non",
            "ex",
            "do",
            "amet"
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
                "name": "Donaldson Grant"
            },
            {
                "id": 1,
                "name": "Eva Daniel"
            },
            {
                "id": 2,
                "name": "Carney Howe"
            }
        ],
        "greeting": "Hello, Estela! You have 6 unread messages.",
        "favoriteFruit": "apple"
    },
    {
        "_id": "5d28de69c0b7b3317204c008",
        "index": 5,
        "guid": "9562fac7-ee3b-459d-babc-714901666565",
        "isActive": true,
        "balance": "$2,101.76",
        "picture": "http://placehold.it/32x32",
        "age": 36,
        "eyeColor": "green",
        "name": {
            "first": "Gayle",
            "last": "Jordan"
        },
        "company": "OPTICALL",
        "email": "gayle.jordan@opticall.ca",
        "phone": "+1 (834) 408-2269",
        "address": "907 Dank Court, Welda, Nebraska, 2593",
        "about": "Ut fugiat ea elit excepteur nostrud qui pariatur aliqua quis cillum proident. Cupidatat ullamco non consectetur minim. Ad ad amet non sint duis occaecat culpa culpa. Dolor aliqua cupidatat ad reprehenderit esse velit adipisicing. Dolore dolore magna ullamco ex dolor enim consectetur dolor deserunt sunt reprehenderit laboris dolore ullamco. Eu et veniam sunt amet do. Sunt velit velit cupidatat est voluptate laborum Lorem.",
        "registered": "Wednesday, July 27, 2016 1:40 AM",
        "latitude": "87.365274",
        "longitude": "71.566986",
        "tags": [
            "commodo",
            "anim",
            "do",
            "eu",
            "adipisicing"
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
                "name": "Alexandra Peters"
            },
            {
                "id": 1,
                "name": "Krista Hartman"
            },
            {
                "id": 2,
                "name": "Alisha Garrett"
            }
        ],
        "greeting": "Hello, Gayle! You have 6 unread messages.",
        "favoriteFruit": "banana"
    }
]
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
	_, err := AddJSON(vars.node.Ipfs(), vars.input1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIPFS_AddJSON_Array(t *testing.T) {
	_, err := AddJSON(vars.node.Ipfs(), vars.input2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIPFS_Teardown(t *testing.T) {
	_ = vars.node.Stop()
	vars.node = nil
}
