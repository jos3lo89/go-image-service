package routes

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func SetupPageRoutes(app *fiber.App) {
	renderer := Renderer()

	app.Get("/:page", dynamicPageHandler(renderer))
}

func Renderer() *template.Template {
	return template.Must(template.ParseGlob("views/*.html"))
}

func dynamicPageHandler(_ *template.Template) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.Params("page")

		page = strings.TrimSuffix(page, ".html")

		if _, err := os.Stat("views/" + page + ".html"); err == nil {
			return c.Render(page, nil)
		}

		return c.Status(http.StatusNotFound).SendString("Página no encontrada")
	}
}
