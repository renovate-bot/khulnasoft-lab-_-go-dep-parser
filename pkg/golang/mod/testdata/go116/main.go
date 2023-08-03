package main

import (
	"log"

	"github.com/khulnasoft-labs/go-dep-parser/pkg/golang/mod"
)

func main() {
	if _, err := mod.Parse(nil); err != nil {
		log.Fatal(err)
	}
}
