package monkeymock

import (
	"log"
	"reflect"
)

// returns the object type name of the held mockedObjectRef
func (m *mockStruct) gRefObjectTypeName() string {
	return getHumanTypeName(m.mockedObjectRef)
}

// returns human readable string representing the type of the given object
// supports custom object types (custom struct and interface names)
func getHumanTypeName(object interface{}) string {
	t := reflect.TypeOf(object)
	var tname string
	if t.Kind() == reflect.Ptr {
		tname = t.Elem().Name()
		if tname == "" {
			return "*" + t.Elem().Kind().String()
		}
		if tname == "rtype" { // handle the occasional reflex.Type ptr
			return t.Elem().Elem().Name()
		}
		return "*" + tname
	}
	tname = t.Name()
	if tname == "" {
		return t.Kind().String()
	}
	return tname
}

// returns (false, <type-as-string>) if the given refObject is not mockable
// otherwise will return (true, "")
func isMockableObjectRef(refObject interface{}) (bool, string) {
	// return true, ""
	t := reflect.TypeOf(refObject)
	switch t.Kind() {
	case reflect.Struct:
		return true, ""
	case reflect.Ptr:
		switch t.Elem().Kind() {
		case reflect.Struct:
			return true, ""
		}
	}
	return false, getHumanTypeName(refObject)
}

func objectRefHasMethod(object interface{}, methodName string) bool {
	v := reflect.ValueOf(object)
	switch v.Kind() {
	case reflect.Interface:
		log.Print("DEBUG: reflect.Interface")
		fallthrough
	case reflect.Ptr:
		methodPtr := v.MethodByName(methodName)
		if methodPtr.IsValid() {
			return true
		}
	default:
		t := reflect.PtrTo(v.Type())
		_, found := t.MethodByName(methodName)
		if found {
			return true
		}
	}
	return false
}
