# ProHttpHandler

A teeny tiny http handler for Go that I wanted to reuse across a couple of projects.

It gives you two things: static asset handling (in a path of your choice) and exact matching of paths to http handler funcs

## Example 

An example that serves all files in "public" at / and serves up a super boring homepage on port 8080.

```go
package main

import (
  "fmt"
  "github.com/joho/prohttphandler"
  "log"
  "net/http"
)

func main() {
  handler := prohttphandler.New("public")

  handler.ExactMatchFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "This is the homepage at /")
  })

  log.Fatal(http.ListenAndServe(":8080", handler))
}
```

There are some very poor [docs here](http://godoc.org/github.com/joho/prohttphandler)

---
Copyright 2013 [John Barton](http://whoisjohnbarton.com) (but under MIT Licence)
