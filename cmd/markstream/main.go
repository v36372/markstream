package main

import (
	"markstream"
  // "net/http"
  // "golang.org/x/net/websocket"
  "github.com/kataras/iris"
  "github.com/iris-contrib/middleware/logger"
)

type page struct {
	Title string
}

func main() {
	ms := markstream.NewMarkStream()
	// go m.InitBackgroundTask()
	go ms.Process()
	go ms.Input()

  iris.Config.Render.Template.Directory = "templates/web/default"
  iris.Config.Render.Template.Gzip = true
  iris.OnError(iris.StatusForbidden, func(ctx *iris.Context) {
		ctx.HTML(iris.StatusForbidden, "<h1> You are not allowed here </h1>")
	})
  iris.Static("/css", "./resources/css", 1)
	iris.Static("/js", "./resources/js", 1)
  iris.Use(logger.New(iris.Logger))
  iris.Get("/", func(ctx *iris.Context) {
		ctx.MustRender("index.html", page{"Hello world"})
	})

	iris.Listen(":8080")
}
