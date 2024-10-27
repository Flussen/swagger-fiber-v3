// Package swagger provides middleware to integrate Swagger UI with Fiber v3,
// allowing API documentation generation from code comments and JSON files.
package swagger

import (
	"fmt"
	"html/template"
	"path"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v3"
	"github.com/swaggo/swag"
)

const (
	defaultDocURL = "doc.json"
	defaultIndex  = "index.html"
)

// HandlerDefault is the default Swagger handler generated by New().
var HandlerDefault = New()

// New returns a custom Fiber handler to serve Swagger UI. It takes optional
// configuration parameters to adjust settings like the documentation URL,
// custom plugins, and UI layout.
//
// The returned handler serves the Swagger documentation and the UI based on
// the specified configuration. It initializes a template for the Swagger UI
// index page and handles requests for the Swagger JSON documentation.
//
// Usage:
//
//	app := fiber.New()
//	app.Get("/docs/*", swagger.HandlerDefault) // example
func New(config ...Config) fiber.Handler {
	cfg := configDefault(config...)

	index, err := template.New("swagger_index.html").Parse(indexTmpl)
	if err != nil {
		panic(fmt.Errorf("fiber: swagger middleware error -> %w", err))
	}

	var (
		prefix string
		once   sync.Once
	)

	return func(c fiber.Ctx) error {
		once.Do(func() {
			prefix = strings.ReplaceAll(c.Route().Path, "*", "")
			forwardedPrefix := getForwardedPrefix(c)
			if forwardedPrefix != "" {
				prefix = forwardedPrefix + prefix
			}

			if len(cfg.URL) == 0 {
				cfg.URL = path.Join(prefix, defaultDocURL)
			}
		})

		p := c.Path(c.Params("*"))

		switch p {
		case defaultIndex:
			c.Type("html")
			return index.Execute(c, cfg)
		case defaultDocURL:
			doc, err := swag.ReadDoc(cfg.InstanceName)
			if err != nil {
				return err
			}
			return c.Type("json").SendString(doc)
		case "", "/":
			c.Set("Location", path.Join(prefix, defaultIndex))
			return c.Status(fiber.StatusMovedPermanently).Send(nil)
		default:
			return c.SendStatus(fiber.StatusNotFound)
		}
	}
}

// getForwardedPrefix extracts the "X-Forwarded-Prefix" header value from the request
// and normalizes it by removing any trailing slashes. This prefix is useful when
// the application is served behind a proxy or load balancer that modifies the route path.
func getForwardedPrefix(c fiber.Ctx) string {
	header := c.Get("X-Forwarded-Prefix")
	if len(header) == 0 {
		return ""
	}

	endIndex := len(header)
	for endIndex > 1 && header[endIndex-1] == '/' {
		endIndex--
	}
	return header[:endIndex]
}
