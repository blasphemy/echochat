package main

import (
	"crypto/sha1"
	"fmt"
	"strings"
)

//TODO make this prettier
func CloakIP6(text string) string {
	stuff := strings.Split(text, ":")
	stuff[len(stuff)-1] = CloakString(stuff[len(stuff)-1], config.Salt)
	if stuff[len(stuff)-2] != "" {
		stuff[len(stuff)-2] = CloakString(stuff[len(stuff)-2], config.Salt)
	} else {
		stuff[len(stuff)-3] = CloakString(stuff[len(stuff)-3], config.Salt)
	}
	text = strings.Join(stuff, ":")
	return text
}

func CloakIP4(text string) string {
	stuff := strings.Split(text, ".")
	stuff[len(stuff)-1] = CloakString(stuff[len(stuff)-1], config.Salt)
	stuff[len(stuff)-2] = CloakString(stuff[len(stuff)-2], config.Salt)
	text = strings.Join(stuff, ".")
	return text
}

func CloakHost(text string) string {
	//hostname mode
	stuff := strings.Split(text, ".")
	stuff[0] = CloakString(stuff[0], config.Salt)
	text = strings.Join(stuff, ".")
	return text
}

func CloakString(text string, salt string) string {
	var r string
	r = Sha1String(text + salt)
	for len(r) < len(text) {
		r = r + Sha1String(r)
	}
	side := true
	for len(r) > len(text) {
		if side {
			r = r[1:]
			side = false
		} else {
			r = r[:len(r)-1]
			side = true
		}
	}
	return r
}

func Sha1String(text string) string {
	data := []byte(text)
	result := fmt.Sprintf("%x", sha1.Sum(data))
	return result
}
