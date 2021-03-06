package fswatch

import (
	"fmt"
	"io/ioutil"

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

	readedDirs := map[string]bool{}
	var readDirs func([]string)
	readDirs = func(dirs []string) {
		var nextDirs []string
		for _, p := range dirs {
			err = watcher.Add(p)
			if err != nil {
				fmt.Println("[error] watcher.Add:", err)
			}
			files, err := ioutil.ReadDir(p)
			if err != nil {
				fmt.Println("[error] watcher.Add:", err)
			}
			for _, v := range files {
				if readedDirs[v.Name()] {
					continue
				}
				if v.IsDir() {
					readedDirs[p+"/"+v.Name()] = true
					nextDirs = append(nextDirs, p+"/"+v.Name())
				}
			}
		}
		if len(nextDirs) > 0 {
			readDirs(nextDirs)
		}
	}
	readDirs(path)

	<-done
}
