package helpers

import (
	"context"
	"encoding/json"
	db "gamecraft-backend/prisma_testing/prisma_testing_client"
)

func QueryRunner(starterSchema string, starterData string, query string, endSchema string) (string, error) {
	// Connect Prisma
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return "error connecting db", err
	}
	defer client.Prisma.Disconnect()

	err := client.Prisma.QueryRaw(starterSchema).Exec(context.Background(), nil)
	if err != nil {
		return "error in starter schema", err
	}


	err = client.Prisma.QueryRaw(starterData).Exec(context.Background(), nil)
	if err != nil {
		return "error in starter data", err
	}


	queryResult := []map[string]interface{}{}
	err = client.Prisma.QueryRaw(query).Exec(context.Background(), &queryResult)
	if err != nil {
		return "error in query", err
	}

	err = client.Prisma.QueryRaw(endSchema).Exec(context.Background(), nil)
	if err != nil {
		return "error in starter data", err
	}

	jsonResult, err := json.Marshal(queryResult)
	if err != nil {
		return "error marshaling query result", err
	}

	return string(jsonResult), nil
}