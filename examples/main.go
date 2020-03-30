package main

import (
	"github.com/atony2099/gos"
)

func main() {

	g := gos.Default()

	g.Router.GET("/hello", func(c *gos.Context) {
		c.Json(200, map[string]interface{}{
			"mesage": "success",
			"code":   1000,
		})
	})

	g.Run(":4000")

}
