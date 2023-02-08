package main

import (
	"github.com/cwheart/flavour"
)

func main() {
	f := flavour.New()
	f.Get("/", func(ctx *flavour.Context) error {
		return ctx.JSON(200, map[string]interface{}{
			"code": "200",
			"msg":  "success",
		})
	})

	f.Get("/hello", func(ctx *flavour.Context) error {
		return ctx.JSON(200, map[string]interface{}{
			"code": "200",
			"msg":  "success",
		})
	})
	f.Start()
}
