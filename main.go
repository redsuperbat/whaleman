package main

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/redsuperbat/whaleman/data"
	"github.com/redsuperbat/whaleman/manifests"
	"github.com/redsuperbat/whaleman/sync"
)

func main() {
	port := ":8090"
	app := iris.New()
	app.Use(iris.Compression)

	// Data
	logger := golog.New()
	data.EnsureDataDir(logger)

	// Manifest routes
	manifests.RegisterManifests(app)

	// Sync routes
	sync.RegisterSync(app)

	app.Listen(port)
}
