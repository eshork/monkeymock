package monkeymock

/*

This part of the module is STALLED until the language is updated with support
for creating new named & exported methods onto existing struct types.

Should be doable once this is complete: https://github.com/golang/go/issues/16522

Alternatively, this is loosely related, but is unlikely as it asks
for functionality that doesn't exist within the language: https://github.com/golang/go/issues/31688

*/

import (
	// 	"fmt"
	"reflect"
	// 	"strings"
	// 	"testing"
	"log"
)

type mockDoubleInterface interface {
	AsDouble() interface{}
}

func (m *mockStruct) AsDouble() interface{} {

	panicAsDoubleNotImplemented()

	// normalize the type
	objPtrType, objType := getNormalizedObjectTypes(reflect.TypeOf(m.mockedObjectRef))

	// make a new struct type that exposes the necessary exported fields
	dupeObjTypeTemp := dupeStructType(objType)
	var _ = dupeStructType(objPtrType)

	// make some normalized handles for our duplicate
	dupePtrType, dupeType := getNormalizedObjectTypes(dupeObjTypeTemp)

	dupeStructMethods(objType, dupeType)       // dupe concrete methods
	dupeStructMethods(objPtrType, dupePtrType) // dupe ptr methods
	return reflect.New(dupeType).Interface()   // hand back our double
}

func dupeStructType(structType reflect.Type) reflect.Type {
	// normalize the input
	_, structTypeConcrete := getNormalizedObjectTypes(structType)

	log.Printf("INTYPE: %+v\n", structTypeConcrete)
	var fields [](reflect.StructField)

	numFields := structTypeConcrete.NumField()
	fields = make([](reflect.StructField), 0, numFields)
	log.Printf("%+v\n", fields)

	for i := 0; i < numFields; i++ {
		field := structTypeConcrete.Field(i)
		log.Printf("%+v\n", field)

		if field.PkgPath == "" {
			fields = append(fields, field)
		}
	}

	newStructType := reflect.StructOf(fields)
	log.Printf("OUTTYPE: %+v\n", newStructType)
	return newStructType
}

// create an empty replica of every Exported struct method
func dupeStructMethods(srcType reflect.Type, destType reflect.Type) {
	log.Print("*** LETS DUPE SOME METHODS ***")
	log.Printf("INTYPE: %+v\n", srcType)

	numMethods := srcType.NumMethod()
	for i := 0; i < numMethods; i++ {
		log.Printf("METHOD #: %d\n", i)
		methodHndl := srcType.Method(i)
		log.Printf("METHOD : %+v\n", methodHndl)
		log.Print(methodHndl.Type.Kind())
		log.Print(methodHndl.Type)
		// Method(int) Method

		log.Printf("METHOD 1in : %+v\n", methodHndl.Type.In(0))

		retvals := []reflect.Value{reflect.ValueOf(7)}

		dupedMethodSig := dupeMethodSignature(methodHndl, destType)
		log.Printf("DUPED METHOD : %+v\n", dupedMethodSig)

		newFuncval := reflect.MakeFunc(dupedMethodSig, GenerateZeroFunctionHandler(retvals))
		log.Printf("METHOD 1new : %+v\n", newFuncval.Type())

	}
	log.Printf("OUTTYPE: %+v\n", destType)

}

// GenerateZeroFunctionHandler returns a dummy function that will return all zero'd value results for the given expected returns types
func GenerateZeroFunctionHandler(returns []reflect.Value) func([]reflect.Value) []reflect.Value {
	zerodValues := make([]reflect.Value, len(returns))
	for i, v := range returns {
		zerodValues[i] = reflect.Zero(v.Type())
	}
	return GenerateEmptyFunctionHandler(zerodValues)
}

// GenerateEmptyFunctionHandler creates an empty function handler for the given signature pattern.
// The return will be a handle to an anonymous function that is willing to accept any input args pattern,
// and will return back the provided return values.
func GenerateEmptyFunctionHandler(returns []reflect.Value) func([]reflect.Value) []reflect.Value {
	// func handle that literally returns the privided returns pattern
	anonFunc := func(args []reflect.Value) (results []reflect.Value) {
		// return []reflect.Value{reflect.ValueOf(7)}
		return returns
	}
	return anonFunc
}

// generates a duplicate object method signature, based on the given srcMethodType
func dupeMethodSignature(srcMethodSig reflect.Method, newTargetObjectType reflect.Type) reflect.Type {
	// if srcMethodSig.Kind() != reflect.Func {
	// 	panic("OH NOES!")
	// }

	methodType := srcMethodSig.Type

	// inTypes := make([]reflect.Type, len(methodType.NumIn()))

	methodNumIn := methodType.NumIn()
	inTypes := make([]reflect.Type, methodNumIn)
	for i := 0; i < methodNumIn; i++ {
		if i == 0 {
			inTypes[i] = newTargetObjectType
		} else {
			inTypes[i] = methodType.In(i)
		}
	}

	methodNumOut := methodType.NumOut()
	outTypes := make([]reflect.Type, methodNumOut)
	for i := 0; i < methodNumOut; i++ {
		outTypes[i] = methodType.Out(i)
	}

	log.Printf("METHOD SIG IN TYPE: %+v\n", inTypes)
	log.Printf("METHOD SIG OUT TYPE: %+v\n", outTypes)

	isVariadic := methodType.IsVariadic()

	newMethodSig := reflect.FuncOf(inTypes, outTypes, isVariadic)

	log.Printf("INPUT SIG: %+v\n", methodType)
	log.Printf("OUTPUT SIG: %+v\n", newMethodSig)

	return newMethodSig
}
