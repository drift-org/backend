package helpers

import (
	"log"
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
