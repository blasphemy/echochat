package log

import "fmt"

func I(msg string) {
	fmt.Printf("INFO: %s\n", msg)
}

func E(msg string) {
	fmt.Printf("ERROR: %s\n", msg)
}
