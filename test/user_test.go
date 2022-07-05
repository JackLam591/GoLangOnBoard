package test

import (
	"encoding/json"
	"fmt"
	"helloworld/controllers"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestUserCountAPI(t *testing.T) {
	tests := []struct {
		description  string // description of the test case
		route        string // route path to test
		expectedCode int    // expected HTTP status code
	}{
		// First test case
		{
			description:  "get HTTP status 200",
			route:        "/api/users/count",
			expectedCode: 200,
		},
		// second test case
		{
			description: "get HTTP response body",
			route:       "/api/users/count",
		},
	}

	app := fiber.New()

	app.Get(*&tests[0].route, controllers.CountUser)

	req, _ := http.NewRequest(http.MethodGet, *&tests[0].route, nil)

	resp, _ := app.Test(req, -1)
	body, _ := ioutil.ReadAll(resp.Body)
	var resBody map[string]interface{}
	if err := json.Unmarshal(body, &resBody); err != nil {
		panic(err)
	}
	fmt.Println(resBody["data"])

	assert.Equalf(t, *&tests[0].expectedCode, resp.StatusCode, *&tests[0].description)    // First test case
	assert.NotEmptyf(t, resBody["data"], "Expected %s not empty", *&tests[1].description) // Second test case

}
