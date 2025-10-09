package main

import (
	"os"

	"github.com/stausee1337/qrpc/lexer"
	"github.com/stausee1337/qrpc/parser"
)

func main() {
	filename := "test.rpc"

	dat, err := os.ReadFile(filename);
	if err != nil {
		panic(err);
	}

	contents := string(dat);
	stream, serr := lexer.LexToStream(contents, filename);
	if serr != nil {
		panic(serr);
	}

	// for _, tok := range stream {
	// 	fmt.Printf("%v\n", tok.String());
	// }

	serr = parser.ParseTokenStream(stream)
	if serr != nil {
		panic(serr);
	}

}
