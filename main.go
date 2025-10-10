package main

import (
	"os"

	"github.com/stausee1337/qrpc/analysis"
	"github.com/stausee1337/qrpc/lexer"
	"github.com/stausee1337/qrpc/parser"
)

func main() {
	filename := "test.rpc"

	dat, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	contents := string(dat)
	stream, serr := lexer.LexToStream(contents, filename)
	if serr != nil {
		panic(serr)
	}

	// for _, tok := range stream {
	// 	fmt.Printf("%v\n", tok.String());
	// }

	items, serr := parser.ParseTokenStream(stream)
	if serr != nil {
		panic(serr)
	}

	_, serr = analysis.AnalyseSyntaxItems(items)
	if serr != nil {
		panic(serr)
	}

}
