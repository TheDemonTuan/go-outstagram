package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"net/http"
)

var Handler = handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: &Resolver{}}))

func WrapHandler(f func(http.ResponseWriter, *http.Request)) func(ctx *fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		fasthttpadaptor.NewFastHTTPHandler(http.HandlerFunc(f))(ctx.Context())
	}
}
