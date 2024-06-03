package routes

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
	"outstagram/graph"
)

func graphqlRouter(r fiber.Router) {
	r.All("", func(c *fiber.Ctx) error {
		graph.WrapHandler(graph.Handler.ServeHTTP)(c)
		return nil
	})

	r.Add("GET", "playground", func(c *fiber.Ctx) error {
		graph.WrapHandler(playground.Handler("GraphQL", "/graphql"))(c)
		return nil
	})
}
