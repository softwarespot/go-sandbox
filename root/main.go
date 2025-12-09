package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	r, err := os.OpenRoot("./data")
	if err != nil {
		panic(err)
	}

	f, err := r.OpenFile("text.txt", os.O_RDONLY, 0o644)
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	// This doesn't work
	if _, err := r.OpenFile("../README.md", os.O_RDONLY, 0o644); err == nil {
		panic("expected error")
	}
}
