package helpers

import (
	"celify/pkg/models"
	"fmt"
	"reflect"
)

func CompareInterfaces(a, b interface{}) bool {
	// Convert maps to a common type for comparison
	a = ConvertMapInterfaceToMapString(a)
	b = ConvertMapInterfaceToMapString(b)
	return reflect.DeepEqual(a, b)
}

func ConvertMapInterfaceToMapString(i interface{}) interface{} {
	t := reflect.TypeOf(i)
	fmt.Printf("Type: %s\n", t.String())
	switch x := i.(type) {
	case map[string]interface{}:
		m1 := map[string]interface{}{}
		for k, v := range x {
			m1[k] = ConvertMapInterfaceToMapString(v)
		}
		return m1
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = ConvertMapInterfaceToMapString(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = ConvertMapInterfaceToMapString(v)
		}
	case *models.TargetData:
		return ConvertMapInterfaceToMapString(x.Data)
	}
	return i
}
