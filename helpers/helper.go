package helpers

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/drift-org/backend/models"
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

/*
Setup the indexes of all the collections

As we create indexes for models, make sure to add them to this function
*/
func SetupIndexes() {
	const INDEX_CREATION_TIME_THRESHOLD = 5

	// Closes the channel and cancel the context when the time is up.
	// Prevents the operation from running forever.
	ctx, cancel := context.WithTimeout(
		context.Background(), INDEX_CREATION_TIME_THRESHOLD*time.Second)
	defer cancel()

	models.CreateChallengeIndex(ctx)
}
