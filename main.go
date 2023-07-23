package main

import (
	"flag"

	"github.com/labstack/echo/v4"
	"github.com/lainio/err2/try"
	"github.com/shynome/pic-in-browser/bilibili"
)

var args struct {
	Addr     string
	headless bool
}

func init() {
	flag.StringVar(&args.Addr, "addr", ":7070", "server listen addr")
	flag.BoolVar(&args.headless, "headless", true, "headless")
}

func main() {
	flag.Parse()

	e := echo.New()
	e.HideBanner = true

	bilibili.Register(e.Group("/bilibili"), args.headless)

	try.To(e.Start(args.Addr))
}
