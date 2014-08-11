package log

import "fmt"

// print to log with level INFO
func I(msg string) {
	fmt.Printf("INFO: %s\n", msg)
}


//print to log with level ERROR
func E(msg string) {
	fmt.Printf("ERROR: %s\n", msg)
}
