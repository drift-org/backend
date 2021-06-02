package helpers

import (
	"log"
	"reflect"
)

/*
This function checks if there is an error.
If an error is present, it will print the message(optional), the error and panic.
Else, it will just return nothing and do nothing
*/
func AlertError(err error, msg ...string) {
	// paramter msg has ... so we could keep it optional.
	// GoLang doesn't support method overload
	if err != nil {
		if len(msg) > 0 {
			log.Println(msg[0])
		}
		log.Fatal(err)
	}
}

/*
This function maps a generic slice to another slice.
*/
func Map(t interface{}, f func(interface{}) interface{}) []interface{} {
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)
		arr := make([]interface{}, s.Len())
		for i := 0; i < s.Len(); i++ {
			arr[i] = f(s.Index(i).Interface())
		}
		return arr
	}
	return nil
}
