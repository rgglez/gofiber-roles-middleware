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
	"strings"

	"encoding/base64"
	"encoding/json"

	fiber "github.com/gofiber/fiber/v2"
)

type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// RequiredRoles defines a  of strings.
	//
	// Required.
	RequiredRoles []string

	// RequireAll defines wether all the RequiredRoles must be present in
	// the claims. If false, then the roles are accepted even if just one
	// of the required is present in the claims.
	//
	// Required. Default: true
	RequireAll bool

	// ClaimsKey defines the name of the key where the roles are stored in
	// the claims. Mutually exclusive with Key.
	//
	// Optional. Default: "urn:zitadel:iam:org:project:roles"
	ClaimsKey string
}

// Set the default configuration values.
var ConfigDefault = Config{
	Next:          nil,
	RequiredRoles: []string{},
	RequireAll:    true,
	ClaimsKey:     "urn:zitadel:iam:org:project:roles",
}

//-----------------------------------------------------------------------------

// ExtractRoles from the roles map in the claims
func ExtractRoles(rolesMap map[string]interface{}) []string {
	roles := []string{}

	// Loop over the keys in the map, which represent the roles
	for role := range rolesMap {
		roles = append(roles, role)
	}

	return roles
}

// CheckRequiredRoles checks roles against requiredRoles based on the requireAll parameter.
func CheckRequiredRoles(roles, requiredRoles []string, requireAll bool) bool {
	// Create a map to store required roles for fast lookup
	requiredRolesMap := make(map[string]bool, len(requiredRoles))
	for _, role := range requiredRoles {
		requiredRolesMap[role] = true
	}

	// Counter for tracking matches if requireAll is true
	matches := 0

	// Check roles based on requireAll flag
	for _, role := range roles {
		if requiredRolesMap[role] {
			if !requireAll {
				// If requireAll is false, return true as soon as one match is found
				return true
			}
			matches++
		}
	}

	// If requireAll is true, ensure all requiredRoles are matched
	if requireAll {
		return matches == len(requiredRoles)
	}

	// If requireAll is false and no match was found
	return false
}

// Middleware
func New(config ...Config) fiber.Handler {
	cfg := ConfigDefault

	if len(config) > 0 {
		cfg = config[0]
		if cfg.ClaimsKey == "" {
			cfg.ClaimsKey = ConfigDefault.ClaimsKey
		}
	}

	return func(c *fiber.Ctx) error {
		// Should we pass?
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Get the token from the Authorization header
		bearer := c.Get("Authorization")
		encodedToken := strings.TrimPrefix(bearer, "Bearer ")
		parts := strings.Split(encodedToken, ".")
		claims, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Forbidden: Invalid token"})
		}
		var customClaims map[string]interface{}
		err = json.Unmarshal(claims, &customClaims)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Forbidden: Invalid token"})
		}
		var rolesMap map[string]interface{}
		var ok bool
		if rolesMap, ok = customClaims[cfg.ClaimsKey].(map[string]interface{}); !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Forbidden: No valid roles found"})
		}
		roles := ExtractRoles(rolesMap)

		// Validate the roles
		if ok = CheckRequiredRoles(roles, cfg.RequiredRoles, cfg.RequireAll); ok {
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Forbidden: No valid roles found"})
	}
}
