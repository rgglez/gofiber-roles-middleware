/*
Copyright 2024 Rodolfo González González

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package gofiberroles

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestZitadelMiddleware(t *testing.T) {
	// Get the test data
	token := os.Getenv("ZITADEL_TOKEN")
	requiredRoles := os.Getenv("ZITADEL_ROLES")
	validRoles := strings.Split(requiredRoles, ",")
	invalidRoles := strings.Split(requiredRoles+",ABCDEFGHIJK", ",")

	testCases := []struct {
		name           string
		expectedStatus int
		roles          []string
		token          string
	}{
		{"Valid roles", http.StatusOK, validRoles, "Bearer " + token},
		{"Invalid roles", http.StatusUnauthorized, invalidRoles, "Bearer " + token},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize Fiber app and middleware
			app := fiber.New()
			app.Use(New(Config{RequiredRoles: tc.roles, RequireAll: true}))
			// Protected route to test the middleware
			app.Get("/", func(c *fiber.Ctx) error {
				return c.SendString("Hello world")
			})

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", tc.token)
			resp, _ := app.Test(req)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}
