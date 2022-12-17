package bilibili

import (
	"github.com/labstack/echo/v4"
)

func Register(g *echo.Group) {
	g.Any("/dynamic-pic", GetDynamicPicHandler)
	g.GET("/dynamic-pic/:id", GetDynamicPicHandler)

	g.Any("/auth", GetAuthHandler)
	g.Any("/auth/:uid", GetAuthHandler)
}
