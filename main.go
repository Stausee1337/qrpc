package main

import (
	"fmt"
	"os"

	"github.com/stausee1337/qrpc/analysis"
	"github.com/stausee1337/qrpc/codegen"
	"github.com/stausee1337/qrpc/lexer"
	"github.com/stausee1337/qrpc/parser"
	"github.com/stausee1337/qrpc/source"
)

func main() {
	files := os.Args[1:]
	analysis, ok := analyzeFilesHandleError(files)
	if !ok {
		os.Exit(1)
	}
	codegen.CodegenFromAnalysis(analysis, "./qmodel")
}

func analyzeFilesHandleError(files []string) (analysis.AnalysisResult, bool) {
	items := make([]parser.Item, 0)

	for _, filename := range files {
		dat, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read file: %v\n", err.Error())
			return analysis.AnalysisResult{}, false
		}

		contents := string(dat)
		stream, serr := lexer.LexToStream(contents, filename)
		if serr != nil {
			source.RenderSourceError(serr)
			return analysis.AnalysisResult{}, false
		}

		fileItems, serr := parser.ParseTokenStream(stream)
		if serr != nil {
			source.RenderSourceError(serr)
			return analysis.AnalysisResult{}, false
		}
		items = append(items, fileItems...)
	}

	res, serr := analysis.AnalyseSyntaxItems(items)
	if serr != nil {
		source.RenderSourceError(serr)
		return analysis.AnalysisResult{}, false
	}

	return res, true
}
