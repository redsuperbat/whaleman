package manifests

import (
	"github.com/kataras/iris/v12"
	"github.com/redsuperbat/whaleman/data"
)

type Msg struct {
	Message string `json:"message"`
}

func createManifestResource(ctx iris.Context) {
	var Body struct {
		Url string `json:"url"`
	}
	if err := ctx.ReadJSON(&Body); err != nil {
		ctx.StatusCode(400)
		ctx.JSON(Msg{Message: "Invalid request body"})
		return
	}

	data.WriteManifestResource(Body.Url)
	ctx.JSON(Msg{Message: "Added resource"})
}

func RegisterManifests(app *iris.Application) {
	manifestResourcesApi := app.Party("/manifests")
	manifestResourcesApi.Post("/", createManifestResource)
}
