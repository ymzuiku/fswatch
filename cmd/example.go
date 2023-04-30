package main

import (
	"log"

	"github.com/ymzuiku/fswatch"
)

func main() {
	fswatch.Watch([]string{"./"}, []string{}, func(file string) {
		log.Println("watch:", file)
	})
}
