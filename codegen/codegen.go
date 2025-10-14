package codegen

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"text/template"

	"github.com/stausee1337/qrpc/analysis"
	"github.com/stausee1337/qrpc/parser"
	"github.com/stausee1337/qrpc/source"
)

//go:embed *.gotpl
var codegenTemplates embed.FS

func generateTypeName(ty any) string {
	switch v := ty.(type) {
	case analysis.Operation:
		return "input"
	case analysis.RecordType:
		return source.ToPascalCase(v.Name.String())
	case analysis.EnumType:
		return source.ToPascalCase(v.Name.String())
	default:
		panic("unknown thing in generateTypeName")
	}
}

func getTypeFields(ty any) []analysis.NamedType {
	switch v := ty.(type) {
	case analysis.Operation:
		return v.InputTypes
	case analysis.RecordType:
		return v.Fields
	default:
		panic("unknown thing in getTypeFields")
	}
}

func generateServiceOperationName(service *analysis.Service, op parser.OperationKind) string {
	switch op {
	case parser.OperationQuery:
		return fmt.Sprintf("%vServiceQueries", service.GoName())
	case parser.OperationMutation:
		return fmt.Sprintf("%vServiceMutation", service.GoName())
	default:
		panic("unreachable")
	}
}

func generateResolveOperationName(service *analysis.Service, op parser.OperationKind) string {
	switch op {
	case parser.OperationQuery:
		return fmt.Sprintf("еееResolve%vServiceQuery", service.GoName())
	case parser.OperationMutation:
		return fmt.Sprintf("еееResolve%vServiceMutation", service.GoName())
	default:
		panic("unreachable")
	}
}

func generateNameFunc(
	f func(service *analysis.Service, op parser.OperationKind) string,
	op parser.OperationKind,
) func(service *analysis.Service) string {
	return func(service *analysis.Service) string {
		return f(service, op);
	}
}

func CodegenFromAnalysis(analysis analysis.AnalysisResult, outDir string) {
	t := template.New("").Funcs(template.FuncMap{
		"go": source.ToPascalCase,
		"typeName": generateTypeName,
		"typeFields": getTypeFields,
		"serviceQueries": generateNameFunc(generateServiceOperationName, parser.OperationQuery),
		"resolveQueries": generateNameFunc(generateResolveOperationName, parser.OperationQuery),
		"serviceMutations": generateNameFunc(generateServiceOperationName, parser.OperationMutation),
		"resolveMutations": generateNameFunc(generateResolveOperationName, parser.OperationMutation),
	})
	t, err := t.ParseFS(codegenTemplates, "*.gotpl")
	if err != nil {
		panic(err)
	}

	root := t.Lookup("generated.gotpl")
	var buf bytes.Buffer
	err = root.Execute(&buf, analysis)
	if err != nil {
		panic(err)
	}

	// buf.WriteTo(os.Stdout)

	file, err := os.Create("generated.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf.WriteTo(file)
}

