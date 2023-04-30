# fswatch

```go
fswatch.Watch([]string{"./"}, []string{}, func(file string) {
  log.Println("watch:", file)
})
```
