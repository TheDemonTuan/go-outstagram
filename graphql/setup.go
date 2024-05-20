package graphql

import (
	"github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql"
)

type Input struct {
	Query         string                 `query:"query"`
	OperationName string                 `query:"operationName"`
	Variables     map[string]interface{} `query:"variables"`
}

func generateSchema() graphql.Schema {
	fields := graphql.Fields{
		"posts":            postGetAll,
		"posts_by_user_id": postGetAllByUserID,
		//"user":  userGetByUserID,
	}
	rootQuery := graphql.ObjectConfig{Name: "Graphql", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		panic(err)
	}

	return schema
}

func Query(ctx *fiber.Ctx) error {
	schema := generateSchema()
	var input Input
	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Cannot parse body: "+err.Error())
	}

	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  input.Query,
		OperationName:  input.OperationName,
		VariableValues: input.Variables,
	})

	ctx.Set("Content-Type", "application/graphql-response+json")
	return ctx.JSON(result)
}

func Mutation(ctx *fiber.Ctx) error {
	schema := generateSchema()
	var input Input
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.
			Status(fiber.StatusInternalServerError).
			SendString("Cannot parse body: " + err.Error())
	}

	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  input.Query,
		OperationName:  input.OperationName,
		VariableValues: input.Variables,
	})

	ctx.Set("Content-Type", "application/graphql-response+json")
	return ctx.JSON(result)
}
