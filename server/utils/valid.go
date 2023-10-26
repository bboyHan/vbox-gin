package utils

import "reflect"

func Contains(slice interface{}, target interface{}) bool {
	sliceValue := reflect.ValueOf(slice)
	itemValue := reflect.ValueOf(target)

	if sliceValue.Kind() != reflect.Slice {
		panic("slice parameter is not a slice")
	}

	for i := 0; i < sliceValue.Len(); i++ {
		if reflect.DeepEqual(sliceValue.Index(i).Interface(), itemValue.Interface()) {
			return true
		}
	}

	return false
}
