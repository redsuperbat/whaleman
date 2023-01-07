package main

import (
	"github.com/kataras/iris/v12"
	"github.com/redsuperbat/whaleman/data"
	"github.com/redsuperbat/whaleman/sync"
)

func main() {
	port := ":8090"
	app := iris.New()
	app.Use(iris.Compression)

	// Data
	data.EnsureDataDir(app.Logger())

	// Sync routes
	sync.RegisterSync(app)

	app.Listen(port)
}
