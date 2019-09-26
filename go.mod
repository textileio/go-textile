module github.com/textileio/go-textile

go 1.12

replace github.com/textileio/go-textile-core v0.0.1 => ../go-bots/go-textile-core

require (
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412
	github.com/ahmetb/govvv v0.2.0 // indirect
	github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc
	github.com/chzyer/readline v0.0.0-20160726135117-62c6fe619375
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/disintegration/imaging v1.6.0
	github.com/evanphx/json-patch v4.1.0+incompatible
	github.com/fatih/color v1.7.0
	github.com/gin-contrib/location v0.0.0-20190301062650-0462caccbb9c
	github.com/gin-contrib/size v0.0.0-20190301062339-6fb8220baadb
	github.com/gin-gonic/gin v1.3.0
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.2
	github.com/hashicorp/go-hclog v0.9.2
	github.com/hashicorp/go-plugin v1.0.1
	github.com/ipfs/go-cid v0.0.2
	github.com/ipfs/go-ipfs v0.4.22-0.20190718080458-55afc478ec02
	github.com/ipfs/go-ipfs-addr v0.0.1
	github.com/ipfs/go-ipfs-cmds v0.1.0
	github.com/ipfs/go-ipfs-config v0.0.6
	github.com/ipfs/go-ipfs-files v0.0.3
	github.com/ipfs/go-ipld-format v0.0.2
	github.com/ipfs/go-log v0.0.1
	github.com/ipfs/go-merkledag v0.2.0
	github.com/ipfs/go-metrics-interface v0.0.1
	github.com/ipfs/go-path v0.0.7
	github.com/ipfs/go-unixfs v0.2.0
	github.com/ipfs/interface-go-ipfs-core v0.1.0
	github.com/libp2p/go-libp2p-core v0.0.9
	github.com/libp2p/go-libp2p-record v0.1.0
	github.com/libp2p/go-msgio v0.0.4
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mr-tron/base58 v1.1.2
	github.com/multiformats/go-multiaddr v0.0.4
	github.com/multiformats/go-multihash v0.0.5
	github.com/mutecomm/go-sqlcipher v0.0.0-20190227152316-55dbde17881f
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/rs/cors v1.6.0
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd
	github.com/segmentio/ksuid v1.0.2
	github.com/stretchr/testify v1.3.0
	github.com/swaggo/gin-swagger v1.1.0
	github.com/swaggo/swag v1.6.2
	github.com/textileio/go-textile-bots v0.0.0-20190926211656-591f4fd421c4
	github.com/textileio/go-textile-core v0.0.1
	github.com/tyler-smith/go-bip39 v1.0.0
	github.com/whyrusleeping/go-logging v0.0.0-20170515211332-0457bb6b88fc
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonschema v1.1.0
	go.uber.org/fx v1.9.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/natefinch/lumberjack.v2 v2.0.0-20170531160350-a96e63847dc3
)
