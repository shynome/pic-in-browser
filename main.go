package main

import (
	"flag"

	"github.com/labstack/echo/v4"
	"github.com/shynome/pic-in-browser/bilibili"
)

var args struct {
	Addr string
}

func init() {
	flag.StringVar(&args.Addr, "addr", ":7070", "server listen addr")
}

func main() {
	flag.Parse()

	e := echo.New()
	e.HideBanner = true

	bilibili.Register(e.Group("/bilibili"))

	e.Start(args.Addr)
}
