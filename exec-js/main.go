package main

import (
	"log"

	v8 "rogchap.com/v8go"
)

func main() {
	code := `
		var x

		function main() {
			x = 4
		}
		main()

		x
	`
	ctx := v8.NewContext()
	v, err := ctx.RunScript(code, "")
	log.Println(v, err)
}
