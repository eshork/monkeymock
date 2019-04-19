package monkeymock

/*

AsPartial layers Mock processing overtop real object instances. It does this
using Monkey Patching, a practice natively supported within certain dynamic
languages (like Ruby and Python) that replaces method/function definitions at
runtime. The Go language is not dynamic nor does it natively support the
redifinition of functions or methods at runtime. So, to make this work, we cheat.

The heavy lifting is implemented by the "bou.ke/monkey" package. The rest is
just syntax sugar to make it read well within your test definitions.
ref: https://github.com/bouk/monkey

Beware that if you're using "bou.ke/monkey" in other areas of your tests or
your application code, you may encounter strange issues if you patch/unpatch
methods that are also being patched/unpatched by this code.

*/

import (
	"reflect"

	"bou.ke/monkey"
)

type mockPartialInterface interface {
	AsPartial() interface{}
}

type mockPartialInterceptRecordKey struct {
	objectType reflect.Type
	methodName string
}

type mockPartialInterceptRecord struct {
	objectType reflect.Type
	methodName string
	patchGuard *monkey.PatchGuard
}

type mockPartialInterceptRecordsMap map[mockPartialInterceptRecordKey]*mockPartialInterceptRecord

var interceptRecords mockPartialInterceptRecordsMap

func (m *mockStruct) AsPartial() interface{} {
	// make sure we have an intercept set up for every method we're currently tracking
	for _, v := range m.mockMethodPtrs {
		createPartialObjectMethodIntercept(reflect.TypeOf(m.mockedObjectRef), v.methodName)
	}
	return m.mockedObjectRef
}

// ensures the framework is ready to intercept calls to the given object type and method name
// - overrides methods of both the given type and pointers to the given type
// - overrides are for all object instances of the given type, the general intercept handler
//     will determine if the particular object itself is worthy of interference
// - multiple calls have no cumulative effect (safe to call multiple times)
// - once the interception is configured, it remains in effect until cleared)
func createPartialObjectMethodIntercept(objectType reflect.Type, methodName string) {
	objectPtrType, objectConcreteType := getNormalizedObjectTypes(objectType)

	// check concrete and ptr (in that order) for existing patch
	if existingPartialObjectMethodIntercept(objectConcreteType, methodName) ||
		existingPartialObjectMethodIntercept(objectPtrType, methodName) {
		return // don't repatch if already existing
	}

	// if concrete type exports the named method, patch it
	if methodHandle, found := objectConcreteType.MethodByName(methodName); found == true {
		patchHandlerFunc := reflect.MakeFunc(methodHandle.Type, func(args []reflect.Value) (results []reflect.Value) {
			return handlePartialObjectMethodIntercept(objectConcreteType, methodName, args)
		}).Interface()
		patchGuard := monkey.PatchInstanceMethod(objectConcreteType, methodName, patchHandlerFunc)
		savePartialObjectMethodIntercept(&mockPartialInterceptRecord{
			objectType: objectConcreteType,
			methodName: methodName,
			patchGuard: patchGuard,
		})

	} else // dont patch ptr if object was patched
	// if ptr type exports the named method, patch it
	if methodHandle, found := objectPtrType.MethodByName(methodName); found == true {
		patchHandlerFunc := reflect.MakeFunc(methodHandle.Type, func(args []reflect.Value) (results []reflect.Value) {
			return handlePartialObjectMethodIntercept(objectPtrType, methodName, args)
		}).Interface()
		patchGuard := monkey.PatchInstanceMethod(objectPtrType, methodName, patchHandlerFunc)
		savePartialObjectMethodIntercept(&mockPartialInterceptRecord{
			objectType: objectPtrType,
			methodName: methodName,
			patchGuard: patchGuard,
		})
	}

}

func savePartialObjectMethodIntercept(record *mockPartialInterceptRecord) {
	if interceptRecords == nil {
		interceptRecords = make(mockPartialInterceptRecordsMap)
	}
	interceptRecords[mockPartialInterceptRecordKey{
		objectType: record.objectType,
		methodName: record.methodName,
	}] = record
}

func existingPartialObjectMethodIntercept(objectType reflect.Type, methodName string) bool {
	record := getPartialObjectMethodIntercept(objectType, methodName)
	if record != nil {
		return true
	}
	return false
}

func getPartialObjectMethodIntercept(objectType reflect.Type, methodName string) *mockPartialInterceptRecord {
	key := mockPartialInterceptRecordKey{objectType, methodName}
	return interceptRecords[key]
}

// clears a single object method intercept, specified by the given type and method name
func clearPartialObjectMethodIntercept(objectType reflect.Type, methodName string) {
	key := mockPartialInterceptRecordKey{objectType, methodName}
	record := interceptRecords[key]
	if record != nil {
		delete(interceptRecords, key)
		if record.patchGuard != nil {
			record.patchGuard.Unpatch()
		}
	}
}

// clears all current known method intercepts
func clearPartialObjectMethodIntercepts() {
	for i := range interceptRecords {
		clearPartialObjectMethodIntercept(i.objectType, i.methodName)
	}
}

// handles a single method intercept in a generic fashion
// - if the object/method combination has no mock defition, normal method execution resumes as quickly as possible
// - if the object/method combination has a mock definition, execution is redirected onto the mock method handler
func handlePartialObjectMethodIntercept(objectType reflect.Type, methodName string, args []reflect.Value) (results []reflect.Value) {
	// if the object is in the known Mock List, then we need to route all calls through the Mock.Call functionality
	for _, v := range gTheMockList {
		if len(args) > 0 && areSameObject(v.(*mockStruct).mockedObjectRef, args[0].Interface()) {
			// temporarily unpatch, in case there is an expectation to call the original
			patchRecord := getPartialObjectMethodIntercept(objectType, methodName)
			patchRecord.patchGuard.Unpatch()
			defer patchRecord.patchGuard.Restore()
			// convert inputs
			interfaceArgs := make([]interface{}, len(args))
			for i, v := range args {
				interfaceArgs[i] = v.Interface()
			}
			// send this off to the normal Mock Method handler
			interfaceRets := v.Call(methodName, interfaceArgs[1:]...)
			// convert outputs
			retList := make([]reflect.Value, len(interfaceRets))
			for i, v := range interfaceRets {
				retList[i] = reflect.ValueOf(v)
			}
			return retList
		}
	}

	// no Mock registered for this object...
	patchRecord := getPartialObjectMethodIntercept(objectType, methodName)
	patchRecord.patchGuard.Unpatch()
	defer patchRecord.patchGuard.Restore()
	methodHndl := args[0].MethodByName(methodName)
	if !methodHndl.IsValid() {
		panic("could not find the expected method on object: " + methodName)
	}
	return methodHndl.Call(args[1:])
}

func areSameObject(leftObj interface{}, rightObj interface{}) bool {
	expectedPtr, actualPtr := reflect.ValueOf(leftObj), reflect.ValueOf(rightObj)
	if expectedPtr.Kind() != reflect.Ptr || actualPtr.Kind() != reflect.Ptr {
		return false
	}
	expectedType, actualType := reflect.TypeOf(leftObj), reflect.TypeOf(rightObj)
	if expectedType != actualType {
		return false
	}
	if leftObj != rightObj {
		return false
	}
	return true
}
