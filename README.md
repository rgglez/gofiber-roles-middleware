[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![GitHub all releases](https://img.shields.io/github/downloads/rgglez/gofiber-roles-middleware/total)
![GitHub issues](https://img.shields.io/github/issues/rgglez/gofiber-roles-middleware)
![GitHub commit activity](https://img.shields.io/github/commit-activity/y/rgglez/gofiber-roles-middleware)
[![Go Report Card](https://goreportcard.com/badge/github.com/rgglez/gofiber-roles-middleware)](https://goreportcard.com/report/github.com/rgglez/gofiber-roles-middleware)
[![GitHub release](https://img.shields.io/github/release/rgglez/gofiber-roles-middleware.svg)](https://github.com/rgglez/gofiber-roles-middleware/releases/)


**gofiber-roles-middleware** is a [gofiber](https://gofiber.io/) [middleware](https://docs.gofiber.io/category/-middleware/) which verifies if the specified role(s) are present in a customizable key in the custom claims part of a given [JWT token](https://jwt.io/).

## Installation

```bash
go get github.com/rgglez/gofiber-roles-middleware
```

## Usage

```go
import gofiberroles "github.com/rgglez/gofiber-roles-middleware/gofiberroles"

// Initialize Fiber app and middleware
app := fiber.New()
app.Use(gofiberroles.New(gofiberroles.Config{RequiredRoles: []string{"admin", "user"}, RequireAll: true}))
```

## Configuration

There are some configuration options available in the ```Config``` struct:

* **```Next```** defines a function to skip this middleware when returned true. Optional. Default: nil
* **```RequiredRoles```** an array of strings which defines the required roles which the user must have (in the claims). Required.
* **```RequireAll```** a boolean which defines if all the required roles must be present in the claims. Optional. Default: true.
* **```ClaimsKey```** a string which will be used as the key to search for the roles in the claims. Optional. Default: "urn:zitadel:iam:org:project:roles". Notice that this is the default used by [Zitadel](https://zitadel.io).

## Notes

* This middleware **does not verify the signature of the token**. It assumes that your program does that verification with some other middleware.
* The middleware was written with [Zitadel](https://zitadel.io) in mind. You might need to make some adjustments so it works with other claims structure.

## Example

An example is provided in the [example/](example/) directory.

### Run it

```bash
cd example
go run main.go
```

### Try it

Then, if it started correctly, assuming that you filled the enviroment variables in ```test_data.sh``` . Or, just set the same variables explained in the Testing section.

```bash
# first step is optional
source /path/to/test_data.sh
```

You will need [pytest](https://en.wikipedia.org/wiki/Pytest) tu run the test:

```bash
cd tests
pytest
```

## Testing

A test is included. To run the test you must:

1. Get a valid JWT token, maybe from your Zitadel instance.
1. Set the test data in the enviroment. An example bash script is provided in [```test_data.sh```](test_data.sh) as a guide. You must fill in the values with your own data accordingly:

    ```bash
    # A token got from a valid login 
    export ZITADEL_TOKEN=
    # A string of comma-separated role names
    export ZITADEL_ROLES=
    ```
    If you use this script, you should need to [source](https://www.geeksforgeeks.org/source-command-in-linux-with-examples/) it.

1. Run
    ```bash
    go test
    ```
    inside the [```gofiberroles/```](gofiberroles/) directory.

## Dependencies

* [github.com/gofiber/fiber/v2](https://github.com/gofiber/fiber/v2)

## License

Copyright (c) 2024 Rodolfo González González

Licensed under the [Apache 2.0](LICENSE) license. Read the [LICENSE](LICENSE) file.
