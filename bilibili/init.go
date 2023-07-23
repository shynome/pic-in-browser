package bilibili

import (
	"os"
	"path/filepath"

	"github.com/chromedp/chromedp"
	"github.com/labstack/echo/v4"
	"github.com/lainio/err2/try"
)

var opts []func(*chromedp.ExecAllocator)

func Register(g *echo.Group, headless bool) {
	homedir := try.To1(os.UserHomeDir())
	userdatadir := filepath.Join(homedir, "./.config/pic-in-browser")
	opts = append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", headless),
		chromedp.UserDataDir(userdatadir),
		chromedp.Flag("password-store", ""),
	)

	g.Any("/dynamic-pic", GetDynamicPicHandler)
	g.GET("/dynamic-pic/:id", GetDynamicPicHandler)
}
