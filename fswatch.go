package fswatch

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var baseIgnore = []string{"node_modules", ".git", ".vscode", ".idea", ".next"}

func Watch(paths, ignore []string, fn func(file string)) {
	ignore = append(ignore, baseIgnore...)
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
			var jump bool
			for _, v := range ignore {
				if strings.Contains(p, v) {
					jump = true
					break
				}
			}
			if jump {
				continue
			}

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
					dir1 := path.Join(p, v.Name())
					readedDirs[dir1] = true
					var jump bool
					for _, v := range ignore {
						if strings.Contains(p, v) {
							jump = true
							break
						}
					}
					if jump {
						continue
					}
					nextDirs = append(nextDirs, dir1)
				}
			}
		}
		if len(nextDirs) > 0 {
			readDirs(nextDirs)
		}
	}
	readDirs(paths)

	<-done
}
