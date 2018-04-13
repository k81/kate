Golang logging package support context.

**Example:**

```go
package main

import (
	"context"

	"github.com/k81/kate/log"
)

func foo(ctx context.Context) {
	ctx = log.SetContext(ctx, "context", "foo")

	log.Info(ctx, "foo called")
}

func bar(ctx context.Context) {
	ctx = log.SetContext(ctx, "context", "bar")

	log.Info(ctx, "bar called")
}

func main() {
	ctx := log.SetContext(context.Background(), "module", "example")
	log.Info(ctx, "program started")
	foo(ctx)
	bar(ctx)
	log.Info(ctx, "program exited")
}
```

**Output**:

```
INFO    2017-04-09 10:30:16.648 [8875] module=[example] msg=[program started] fileline=[main.go:23]
INFO    2017-04-09 10:30:16.648 [8875] module=[example] context=[foo] msg=[foo called] fileline=[main.go:12]
INFO    2017-04-09 10:30:16.648 [8875] module=[example] context=[bar] msg=[bar called] fileline=[main.go:18]
INFO    2017-04-09 10:30:16.648 [8875] module=[example] msg=[program exited] fileline=[main.go:26]
```
