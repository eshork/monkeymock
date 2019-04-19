package monkeymock

import (
	"reflect"
)

// takes a relfect.Type, and always returns (*objectType, objectType), ie a pointer Type and a concrete Type
// - Works consistently whether the given objectType was originally a pointer type or a concrete type
// - If a ptr to a ptr is encountered (ie, **obj), or any depth thereof (ex: *****obj), the root obj and obj pointer will be resolved
func getNormalizedObjectTypes(objectType reflect.Type) (reflect.Type, reflect.Type) {
	if objectType.Kind() == reflect.Ptr { // have a ptr type
		rootObjType := objectType.Elem()
		if rootObjType.Kind() == reflect.Ptr { // is this a ptr to a ptr? If so, recurse
			return getNormalizedObjectTypes(rootObjType)
		}
		return objectType, rootObjType
	}
	// not a ptr, must be concrete type
	return reflect.PtrTo(objectType), objectType
}
