package fswatch

import (
"fmt"

"github.com/fsnotify/fsnotify"
)

func Watch(path []string, fn func(file string)) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("watch error:", err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// fmt.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					fn(event.Name)
					// log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("[error] watcher.Errors:", err)
			}
		}
	}()

	for _, p := range path {
		err = watcher.Add(p)
		if err != nil {
			fmt.Println("[error] watcher.Add:", err)
		}
	}

	<-done
}
