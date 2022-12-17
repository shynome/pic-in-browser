package bilibili

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/labstack/echo/v4"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

type AuthParams struct {
	UID string `param:"uid" json:"uid" form:"uid" query:"uid"`
}

func GetAuthHandler(c echo.Context) (err error) {
	defer err2.Handle(&err)

	var params AuthParams
	try.To(c.Bind(&params))

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserDataDir(path.Join(os.Getenv("HOME"), "/.config/bilibili-browser-user/", params.UID)),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	auth := try.To1(GetAuth(ctx, params.UID))
	return c.JSON(200, auth)
}

type Auth struct {
	UID      string `json:"uid"`
	SESSDATA string `json:"SESSDATA"`
	BiliJct  string `json:"bili_jct"`
}

func GetAuth(ctx context.Context, uid string) (auth Auth, err error) {
	defer err2.Handle(&err)

	actions := []chromedp.Action{
		chromedp.Navigate("https://www.bilibili.com/account/history"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetCookies().Do(ctx)
			if err != nil {
				return err
			}
			auth = Auth{UID: uid}
			c := 0
			for _, cookie := range cookies {
				switch cookie.Name {
				case "SESSDATA":
					auth.SESSDATA = cookie.Value
					c++
				case "bili_jct":
					auth.BiliJct = cookie.Value
					c++
				}
				if c == 2 {
					break
				}
			}
			if c != 2 {
				return fmt.Errorf("got cookie failed")
			}
			return nil
		}),
	}

	err = chromedp.Run(ctx, actions...)

	return
}
