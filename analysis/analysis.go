package analysis

import (
	"github.com/stausee1337/qrpc/parser"
	"github.com/stausee1337/qrpc/source"
)

func AnalyseSyntaxItems(items []parser.Item) (AnalysisResult, *source.SourceError) {	
	if err := runWellFormedChecks(items); err != nil {
		return AnalysisResult{}, err;
	}
	
	resolver := makeResolver()	
	if err := resolver.resolve(items); err != nil {
		return AnalysisResult{}, err;
	}

	return AnalysisResult{}, nil;
}

