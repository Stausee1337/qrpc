package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/urfave/cli/v2"

	"github.com/stausee1337/qrpc/analysis"
	"github.com/stausee1337/qrpc/codegen"
	"github.com/stausee1337/qrpc/lexer"
	"github.com/stausee1337/qrpc/parser"
	"github.com/stausee1337/qrpc/source"
)

type fileDesc struct {
	filename string
	contents string
}

func genSchemasFromFiles(files []string, outputDir string) int {
	fileDescs := make([]fileDesc, 0)
	for _, filename := range files {
		dat, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read file: %v\n", err.Error())
			return 1
		}
		contents := string(dat)
		fileDescs = append(fileDescs, fileDesc{ filename: filename, contents: contents })
	}

	analysis, ok := analyzeFilesHandleError(fileDescs)
	if !ok {
		return 1
	}

	packageName := path.Base(outputDir)
	err := codegen.CodegenFromAnalysis(
		analysis,
		codegen.Options{
			OutDir: outputDir,
			PackageName: packageName,
		},
	)
	if err != nil {
		return 1
	}
	
	return 0
}

func main() {
	app := cli.NewApp()
	app.Name = "qrpc-gen"
	app.Usage = "qrpc-gen --output-dir value [...input-files]"
	app.Description = "Generate Go Models and Resolvers from .qrpc Files"
	app.HideVersion = true
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name: "output-dir",
			Usage: "Specifies the Output Directory",
			Required: true,
		},
	};
	app.Version = "0.1.4"

	app.Action = func(ctx *cli.Context) error {
		args := ctx.Args() 
		stringArgs := make([]string, 0)
		for i := 0; i < args.Len(); i++ {
			stringArgs = append(stringArgs, args.Get(i))
		}

		outputDir := ctx.String("output-dir")

		result := genSchemasFromFiles(stringArgs, outputDir);
		if result != 0 {
			return errors.New("some error")
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func analyzeFilesHandleError(fileDescs []fileDesc) (analysis.AnalysisResult, bool) {
	items := make([]parser.Item, 0)

	for _, desc := range fileDescs {
		stream, serr := lexer.LexToStream(desc.contents, desc.filename)
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
