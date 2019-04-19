package monkeymock

import (
	"reflect"
	"strings"
)

func stringifyMethodName(m *mockMethodStruct) string {
	objectName := m.parentMockStruct.gRefObjectTypeName()
	methodName := "<" + objectName + ">." + m.methodName
	return methodName
}

func stringifyMethodReturns(m *mockMethodStruct) string {
	return "value <type>"
}

func stringifyMethodArgs(m *mockMethodStruct) string {
	return "value <type>, value <type>, value <type>"
}

// func stringifyValuesList(typeList []reflect.Value) string {
// 	listTypeSig := ""
// 	for _, v := range typeList {
// 		listTypeSig += "<" + v.Type().String() + ">, "
// 	}
// 	return strings.Trim(listTypeSig, " ,")
// }

func stringifyTypesList(typeList []reflect.Type) string {
	listTypeSig := ""
	for _, v := range typeList {
		listTypeSig += "<" + v.String() + ">, "
	}
	return strings.Trim(listTypeSig, " ,")
}

func typeListToString(typeList []interface{}) string {
	listTypeSig := ""
	for _, v := range typeList {
		listTypeSig += "<" + getHumanTypeName(v) + ">, "
	}
	return strings.Trim(listTypeSig, " ,")
}

func getObjectMethodByName(object interface{}, methodName string) *reflect.Method {
	objType := reflect.TypeOf(object)
	methodHandle, found := objType.MethodByName(methodName)
	if found {
		return &methodHandle
	}
	return nil
}

func (m *mockMethodStruct) getObjectMethodArgTypes() []reflect.Type {
	methodPtr := getObjectMethodByName(m.parentMockStruct.mockedObjectRef, m.methodName)
	if methodPtr != nil {
		methodPtrType := reflect.TypeOf(methodPtr.Func.Interface())
		num := methodPtrType.NumIn()
		retVal := make([]reflect.Type, num-1)
		for i := 0; i < num-1; i++ {
			argI := i + 1
			retVal[i] = methodPtrType.In(argI)
		}
		return retVal
	}
	return []reflect.Type{}
}

func getArgsListTypes(args methodArgumentsList) []reflect.Type {
	retVal := make([]reflect.Type, len(args))
	for i, v := range args {
		retVal[i] = reflect.TypeOf(v)
	}
	return retVal
}

// throw a panic if the given args list does not match the method signature
func (m *mockMethodStruct) ensureMethodArgs(args methodArgumentsList) {
	expectedArgs := m.getObjectMethodArgTypes()
	// log.Print(expectedArgs)
	givenArgs := getArgsListTypes(args)
	// log.Print(givenArgs)

	if len(expectedArgs) != len(givenArgs) {
		// panic
		panicMockWithArgsMismatch(m, stringifyTypesList(expectedArgs), "receivedArgs string")
		// mockMethod *mockMethodStruct, expectedArgs string, receivedArgs string
	}
	for i, v := range expectedArgs {
		if givenArgs[i] != v {
			// panic
			// panicMockWithArgsMismatch
			panicMockWithArgsMismatch(m, stringifyTypesList(expectedArgs), stringifyTypesList(givenArgs))
		}
	}
	// seems good, carry on
}

func callObjectMethodByName(methodHandle *reflect.Method, object interface{}, args methodArgumentsList) methodReturnsList {
	// build the args list
	in := make([]reflect.Value, len(args)+1)
	in[0] = reflect.ValueOf(object) // first argument is the reference object itself
	// the rest are converted in order over to reflect.Value types
	for i, v := range args {
		in[i+1] = reflect.ValueOf(v)
	}

	// do the call
	returnedValueArray := methodHandle.Func.Call(in)

	// convert the returned Value array into something more generic
	returnsList := make(methodReturnsList, len(returnedValueArray))
	for i, v := range returnedValueArray {
		returnsList[i] = v.Interface()
	}

	return returnsList
}

// make a shallow copy of interface an list
func copyInterfaceList(interfaceList []interface{}) []interface{} {
	retsList := make([]interface{}, len(interfaceList))
	copy(retsList, interfaceList)
	return retsList
}
