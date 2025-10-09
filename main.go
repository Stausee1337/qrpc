package main

import (
	"fmt"
	"os"

	"github.com/stausee1337/qrpc/lexer"
)



func main() {
	dat, err := os.ReadFile("test.rpc");
	if err != nil {
		panic(err);
	}

	contents := string(dat);
	stream, err := lexer.LexToStream(contents);
	if err != nil {
		panic(err);
	}

	for _, tok := range stream {
		fmt.Printf("%v\n", tok.String());
	}


}
