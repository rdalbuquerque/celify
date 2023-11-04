package helpers

import (
	"fmt"
	"reflect"
)

func compareInterfaces(a, b interface{}) bool {
	// Convert maps to a common type for comparison
	a = convertMapInterfaceToMapString(a)
	b = convertMapInterfaceToMapString(b)
	return reflect.DeepEqual(a, b)
}

func convertMapInterfaceToMapString(i interface{}) interface{} {
	t := reflect.TypeOf(i)
	fmt.Printf("Type: %s\n", t.String())
	switch x := i.(type) {
	case map[string]interface{}:
		m1 := map[string]interface{}{}
		for k, v := range x {
			m1[k] = convertMapInterfaceToMapString(v)
		}
		return m1
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convertMapInterfaceToMapString(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convertMapInterfaceToMapString(v)
		}
	}
	return i
}
