package marshal

import "reflect"

func Marshal(obj interface{}) ([]string, []interface{}) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	values := make([]interface{}, v.NumField())
	names := make([]string, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
		names[i] = t.Field(i).Tag.Get("sql")
	}
	return names, values
}
