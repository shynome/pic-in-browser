package bilibili

import (
	"github.com/labstack/echo/v4"
)

func Register(g *echo.Group) {
	g.Any("/dynamic-pic", GetDynamicPic)
}
