package main

import (
	"fmt"
	"reflect"
	"syscall/js"
	"unicode"
	"unicode/utf8"
	"unsafe"

	"github.com/stausee1337/qrpc/analysis"
	"github.com/stausee1337/qrpc/lexer"
	"github.com/stausee1337/qrpc/parser"
)


func test(this js.Value, args []js.Value) any {
	contents := args[0].String()
	filename := args[1].String()

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

	result, serr := analysis.AnalyseSyntaxItems(items)
	if serr != nil {
		panic(serr)
	}	

	return convertAny(result);
}


func convertAny(x any) js.Value {
	switch x := x.(type) {
	case js.Value:
		return x
	case js.Func:
		return x.Value
	case nil:
		return js.ValueOf(nil)
	case bool:
		return js.ValueOf(x)
	case int:
		return js.ValueOf(float64(x))
	case int8:
		return js.ValueOf(float64(x))
	case int16:
		return js.ValueOf(float64(x))
	case int32:
		return js.ValueOf(float64(x))
	case int64:
		return js.ValueOf(float64(x))
	case uint:
		return js.ValueOf(float64(x))
	case uint8:
		return js.ValueOf(float64(x))
	case uint16:
		return js.ValueOf(float64(x))
	case uint32:
		return js.ValueOf(float64(x))
	case uint64:
		return js.ValueOf(float64(x))
	case uintptr:
		return js.ValueOf(float64(x))
	case float32:
		return js.ValueOf(float64(x))
	case float64:
		return js.ValueOf(x)
	case string:
		return js.ValueOf(x)
	case map[string]any:
		return js.ValueOf(x)
	default:
		return convertObjectRecursively(x)
	}
}

var globalObject = js.Global()
var objectConstructor = globalObject.Get("Object")
var arrayConstructor = globalObject.Get("Array")

var typeCache = map[unsafe.Pointer]js.Value{}

type eface struct {
    _type unsafe.Pointer
    data  unsafe.Pointer
}

func obtainIfacePointer(x any) unsafe.Pointer {
    return (*eface)(unsafe.Pointer(&x)).data
}


func firstToLower(s string) string {
    r, size := utf8.DecodeRuneInString(s)
    if r == utf8.RuneError && size <= 1 {
        return s
    }
    lc := unicode.ToLower(r)
    if r == lc {
        return s
    }
    return string(lc) + s[size:]
}

func convertObjectRecursively(obj any) js.Value {
	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		iface := obtainIfacePointer(obj)
		result, ok := typeCache[iface]
		if ok {
			return result;
		}
		result = objectConstructor.New()

		// Loop through fields
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := val.Type().Field(i)
			// fmt.Printf("Field Name: %s, Value: %v, Type: %s\n", fieldType.Name, field.Interface(), field.Type())
			result.Set(firstToLower(fieldType.Name), convertAny(field.Interface()))
		}

		result.Set("rawType", val.Type().Name())

		typeCache[iface] = result
		return result
	case reflect.Slice:
		iface := obtainIfacePointer(obj)
		result, ok := typeCache[iface]
		if ok {
			return result;
		}
	 	result = arrayConstructor.New(val.Len())

	 	// Loop through fields
	 	for i := 0; i < val.Len(); i++ {
	 		field := val.Index(i)
	 		// fmt.Printf("Field Name: %s, Value: %v, Type: %s\n", fieldType.Name, field.Interface(), field.Type())
	 		result.SetIndex(i, convertAny(field.Interface()))
	 	}

		typeCache[iface] = result
	 	return result
	default:
		return js.ValueOf(fmt.Sprintf("%v", val.Interface()))
	}	
}

func main() {
    c := make(chan struct{})
	js.Global().Set("test", js.FuncOf(test))
	<-c
}

