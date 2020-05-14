# mapper
Ingress, Pokémon Go, and Wizards Unite plugin for Snart.

## how to build
build [snart](https://github.com/go-snart/snart) normally, but add a `plugins.go` file like this:
```go
package main

import (
	_ "github.com/superloach/mapper"
)
```
