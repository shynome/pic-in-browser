package bilibili

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/jellydator/ttlcache/v3"
	"github.com/labstack/echo/v4"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

type DynamicParams struct {
	ID string `param:"id" json:"id" form:"id" query:"id"`
}

func GetDynamicPicHandler(c echo.Context) (err error) {
	defer err2.Handle(&err)

	var params DynamicParams
	try.To(c.Bind(&params))

	f, err := GetDynamicPicWithCache(c.Request().Context(), params.ID)
	if err != nil {
		return
	}

	return c.File(f)
}

type cacheItem struct {
	rw   *sync.RWMutex
	file string
}

var dynamicPicCache = func() *ttlcache.Cache[string, *cacheItem] {
	cache := ttlcache.New(
		ttlcache.WithTTL[string, *cacheItem](5 * time.Minute),
	)
	cache.OnEviction(func(ctx context.Context, er ttlcache.EvictionReason, i *ttlcache.Item[string, *cacheItem]) {
		os.Remove(i.Value().file)
	})
	go cache.Start()
	return cache
}()

func GetDynamicPicWithCache(ctx context.Context, id string) (string, error) {
	item := dynamicPicCache.Get(id)
	if item == nil {
		rw := &sync.RWMutex{}
		citem := &cacheItem{rw: rw, file: ""}
		item = dynamicPicCache.Set(id, citem, ttlcache.DefaultTTL)
	}
	if f := dynamicExistPic(item); f != "" {
		// 更新缓存过期时间
		dynamicPicCache.Set(item.Key(), item.Value(), ttlcache.DefaultTTL)
		return f, nil
	}

	citem := item.Value()
	citem.rw.Lock()
	defer citem.rw.Unlock()

	f := fmt.Sprintf("/tmp/bilibili-dynamic-%s", id)

	ctx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	img := try.To1(GetDynamicPic(ctx, id))
	if err := os.WriteFile(f, img, 0644); err != nil {
		return "", nil
	}
	citem.file = f

	return f, nil
}

func dynamicExistPic(item *ttlcache.Item[string, *cacheItem]) string {
	if item == nil {
		return ""
	}
	i := item.Value()
	i.rw.RLock()
	defer i.rw.RUnlock()
	return i.file
}

func GetDynamicPic(ctx context.Context, id string) (img []byte, err error) {
	defer err2.Handle(&err)

	link := fmt.Sprintf("https://m.bilibili.com/opus/%s", id)

	var currentLink = ""
	tasks := chromedp.Tasks{
		emulation.SetTimezoneOverride("Asia/Shanghai"),
		emulation.SetUserAgentOverride("Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1").WithPlatform("iPhone"),
		emulation.SetDeviceMetricsOverride(400, 800, 2.5, true),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return chromedp.Run(ctx,
				page.Enable(),
				page.SetLifecycleEventsEnabled(true),
			)
		}),
		chromedp.Navigate(link),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// copy from https://github.com/chromedp/chromedp/issues/431#issuecomment-592950397
			cctx, cancel := context.WithCancel(ctx)
			chromedp.ListenTarget(cctx, func(ev interface{}) {
				switch e := ev.(type) {
				case *page.EventLifecycleEvent:
					if e.Name == "networkIdle" {
						cancel()
					}
				}
			})
			select {
			case <-cctx.Done():
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}),
		// chromedp.WaitReady(".dyn-card *"),
		chromedp.Evaluate(getClearElemJs, nil),
		chromedp.Location(&currentLink),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var elem = ".opus-modules"
			switch {
			case strings.HasPrefix(currentLink, "https://m.bilibili.com/dynamic/"):
				elem = ".dyn-card"
			}
			q := chromedp.Screenshot(elem, &img, chromedp.NodeVisible)
			return q.Do(ctx)
		}),
	}
	try.To(chromedp.Run(ctx, tasks...))

	return
}

const getClearElemJs = `
{
let hidden = selector => {
	let e = document.querySelector(selector)
	if(e) e.style.display = "none";
}

// dynamic
hidden(".dyn-header__right");
hidden(".launch-app-btn.dynamic-float-openapp");

// opus
hidden(".easy-follow-btn");
hidden(".launch-app-btn.float-openapp");
hidden(".openapp-dialog");
hidden("..fixed-openapp");

// 隐藏验证弹窗
hidden('.geetest_panel')
}
`
