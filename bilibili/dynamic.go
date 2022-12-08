package bilibili

import (
	"fmt"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/labstack/echo/v4"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

type DynamicParams struct {
	ID string `param:"id" json:"id" form:"id" query:"id"`
}

func GetDynamicPic(c echo.Context) (err error) {
	defer err2.Handle(&err)

	ctx, cancel := chromedp.NewContext(c.Request().Context())
	defer cancel()

	var params DynamicParams
	try.To(c.Bind(&params))

	link := fmt.Sprintf("https://m.bilibili.com/dynamic/%s", params.ID)

	var img []byte
	tasks := chromedp.Tasks{
		emulation.SetUserAgentOverride("Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1").WithPlatform("iPhone"),
		emulation.SetDeviceMetricsOverride(400, 800, 2.5, true),
		chromedp.Navigate(link),
		chromedp.WaitReady(".dyn-card *"),
		chromedp.Evaluate(getClearElemJs, nil),
		chromedp.Screenshot(".dyn-card", &img, chromedp.NodeVisible),
	}
	try.To(chromedp.Run(ctx, tasks...))

	return c.Blob(200, "image/png", img)
}

const getClearElemJs = `
let c = document.querySelector(".dyn-header__right");
if(c){ c.hidden = true };
document.querySelector(".launch-app-btn.dynamic-float-openapp").style.display = "none";
`
