package helpers

import "fmt"

// Handle Errors
func HandleGeneralErr(err error, outStr string) {
	if err != nil {
		panic(fmt.Errorf("%s: %s", outStr, err))
	}
}
