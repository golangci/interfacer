# interfacer

A code checker that suggests interface types.

If a function takes a parameter of type `*os.File` but all it does is
`Read` from it, this program will suggest that you use `io.Reader`
instead.

## TODOs

* Suggest more interface types

* Ignore functions that implement common interfaces or func types 

```go
// do not suggest io.Writer
func httpHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte{})
})
```

* Field usage - cannot be interface type

```go
// do not suggest io.Closer
func foo(a someType) {
	a.Close()
	a.field = nil
}
```
