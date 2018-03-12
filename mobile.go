package main

import "github.com/textileio/textile-go/mobile"

func main()  {
	textile := mobile.NewTextile("/Users/sander/go/src/github.com/textileio/textile-go/.ipfs")
	textile.Start()
}
