package main

import (
	"flag"
	"fmt"
)

func main() {
	name := flag.String("name", "world", "the name to greet")
	flag.Parse()
	fmt.Printf("Hello, %s!\n", *name)
}
