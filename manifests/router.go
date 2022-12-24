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
		ctx.JSON(Msg{Message: "Invalid request body: " + err.Error()})
		return
	}

	if err := data.WriteManifestResource(Body.Url); err != nil {
		ctx.StatusCode(500)
		ctx.JSON(Msg{Message: "Unable to add manifest resource"})
		ctx.Application().Logger().Error(err)
		return
	}
	ctx.JSON(Msg{Message: "Added resource: " + Body.Url})
}

func getManifestResources(ctx iris.Context) {
	if err, s := data.ReadManifestResources(); err != nil {
		ctx.StatusCode(500)
		ctx.JSON(Msg{Message: "Unable to read manifest resources"})
		ctx.Application().Logger().Error(err)
	} else {
		ctx.StatusCode(200)
		ctx.JSON(s)
	}
}

func RegisterManifests(app *iris.Application) {
	manifestResourcesApi := app.Party("/manifests")
	manifestResourcesApi.Post("/", createManifestResource)
	manifestResourcesApi.Get("/", getManifestResources)
}
