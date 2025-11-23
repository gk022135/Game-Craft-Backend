package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	db "gamecraft-backend/prisma_testing/prisma_testing_client"
	"reflect"
	"sort"
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

	defer func() {
		err = client.Prisma.QueryRaw(endSchema).Exec(context.Background(), nil)
		if err != nil {
			// You cannot return from inside a defer — simply log or modify values
			fmt.Println("error in starter data:", err)
		}
	}()


	err = client.Prisma.QueryRaw(starterData).Exec(context.Background(), nil)
	if err != nil {
		return "error in starter data", err
	}


	queryResult := []map[string]interface{}{}
	err = client.Prisma.QueryRaw(query).Exec(context.Background(), &queryResult)
	if err != nil {
		return "error in query", err
	}

	jsonResult, err := json.Marshal(queryResult)
	if err != nil {
		return "error marshaling query result", err
	}

	return string(jsonResult), nil
}


// sortByID as you requested (assumes JSON numbers -> float64)
func sortByID(rows []map[string]interface{}) {
	sort.Slice(rows, func(i, j int) bool {
		return rows[i]["id"].(float64) < rows[j]["id"].(float64)
	})
}

// CompareResults compares two JSON-stringified query results.
// - testingResult and userResult are JSON strings (e.g. `[{"id":1,"name":"a"}, ...]`).
// - orderMatters: if true, row order must match exactly.
// Returns true if results are considered equal.
func CompareResults(testingResult string, userResult string, orderMatters bool) bool {
	var expected []map[string]interface{}
	var actual []map[string]interface{}

	// Unmarshal both JSON strings into []map[string]interface{}
	if err := json.Unmarshal([]byte(testingResult), &expected); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(userResult), &actual); err != nil {
		return false
	}

	// If order matters, do a direct deep-equal (same order required)
	if orderMatters {
		return reflect.DeepEqual(expected, actual)
	}

	// Order does not matter:
	// Quick length check
	if len(expected) != len(actual) {
		return false
	}

	// Try to detect if every row has a numeric "id" field (JSON numbers are float64)
	allHaveNumericID := true
	for _, r := range expected {
		v, ok := r["id"]
		if !ok {
			allHaveNumericID = false
			break
		}
		_, okf := v.(float64)
		if !okf {
			allHaveNumericID = false
			break
		}
	}
	if allHaveNumericID {
		// Also ensure actual has numeric ids for safety
		for _, r := range actual {
			v, ok := r["id"]
			if !ok {
				allHaveNumericID = false
				break
			}
			_, okf := v.(float64)
			if !okf {
				allHaveNumericID = false
				break
			}
		}
	}

	if allHaveNumericID {
		// Sort both slices by id and compare
		sortByID(expected)
		sortByID(actual)
		return reflect.DeepEqual(expected, actual)
	}

	// Fallback: normalize each row to a JSON string, sort the string slices, then compare
	norm := func(rows []map[string]interface{}) []string {
		out := make([]string, 0, len(rows))
		for _, r := range rows {
			b, err := json.Marshal(r)
			if err != nil {
				// If marshaling fails for any row, return an empty slice to force false later
				return []string{}
			}
			out = append(out, string(b))
		}
		sort.Strings(out)
		return out
	}

	ne := norm(expected)
	na := norm(actual)
	if len(ne) == 0 || len(na) == 0 {
		return false
	}
	return reflect.DeepEqual(ne, na)
}