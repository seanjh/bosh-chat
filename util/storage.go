package util

import (
	"log"
	"os"
)

// EnsurePath TODO
func EnsurePath(path string) {
	_, err := os.Stat(path)
	if err != nil {
		log.Println("Creating session storage path:", path)
		if os.Mkdir(path, 0700) != nil {
			log.Panicln("Failed to create storage path:", path)
		}
	}
}
