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

var HandlerDefault = New()

// New returns custom handler
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
