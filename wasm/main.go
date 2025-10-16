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

func goInit(this js.Value, args []js.Value) any {
	constructors = args[0]

	exports := objectConstructor.New()
	exports.Set("analyzeSourceFiles", js.FuncOf(analyzeSourceFiles))
	return exports
}

func analyzeSourceFiles(this js.Value, args []js.Value) any {
	files := args[0]

	items := make([]parser.Item, 0)
	for i := 0; i < files.Length(); i++ {
		fileDesc := files.Index(i)
		contents, filename := fileDesc.Get("contents").String(), fileDesc.Get("filename").String()

		stream, serr := lexer.LexToStream(contents, filename)
		if serr != nil {
			return convertAny(serr)
		}

		fileItems, serr := parser.ParseTokenStream(stream)
		if serr != nil {
			return convertAny(serr)
		}

		items = append(items, fileItems...)
	}

	result, serr := analysis.AnalyseSyntaxItems(items)
	if serr != nil {
		return convertAny(serr)
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
var constructors js.Value

func constructObject(ty string) js.Value {
	return constructors.Get(ty).New()
}

var udtCache = map[unsafe.Pointer]js.Value{}
var tyCache = map[string]js.Value{}
var typeIfaceType = reflect.TypeFor[analysis.Type]()

func getType(val reflect.Value) analysis.Type {
	v := val
	if !v.Type().Implements(typeIfaceType) && reflect.PointerTo(v.Type()).Implements(typeIfaceType) {
		v = v.Addr()
	} else if !v.Type().Implements(typeIfaceType) {
		return nil
	}

	return v.Interface().(analysis.Type)
}

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

func getCachedOrNew(val reflect.Value) (js.Value, bool) {
	obj := val.Interface()
	ty := getType(val)

	if ty != nil {
		if val.Type().Size() > 0 {
			iface := obtainIfacePointer(obj)
			result, cached := udtCache[iface]
			if cached {
				return result, true
			}
			
			result = constructObject(val.Type().Name())
			udtCache[iface] = result

			return result, false
		} else {
			result, cached := tyCache[val.Type().Name()]
			if cached {
				return result, true
			}
			
			result = constructObject(val.Type().Name())
			tyCache[val.Type().Name()] = result

			return result, false
		}
		// _, isTy := obj.(analysis.Type)
		// reflect.PointerTo()
		// isTy := reflect.PointerTo(
		// 	val.Type(),
		// ).Implements(reflect.TypeFor[analysis.Type]())
	}


	result := constructObject(val.Type().Name())
	return result, false
}

func convertObjectRecursively(obj any) js.Value {
	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		result, cached := getCachedOrNew(val)
		if cached {
			return result;
		}

		// Loop through fields
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := val.Type().Field(i)
			// fmt.Printf("Field Name: %s, Value: %v, Type: %s\n", fieldType.Name, field.Interface(), field.Type())
			result.Set(firstToLower(fieldType.Name), convertAny(field.Interface()))
		}

		// result.Set("rawType", val.Type().Name())

		return result
	case reflect.Slice:
		iface := obtainIfacePointer(obj)
		result, ok := udtCache[iface]
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

		udtCache[iface] = result
	 	return result
	default:
		return js.ValueOf(fmt.Sprintf("%v", val.Interface()))
	}	
}

func main() {
    c := make(chan struct{})
	js.Global().Set("$goInit", js.FuncOf(goInit))
	<-c
}

