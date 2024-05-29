package dropsim

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
)

func ScreenshotTest() {
	browser := rod.New().MustConnect()
	page := browser.
		MustPage("file:////home/frankie/workspace/github.com/phantomwaves/proto-bot/response/results_template.html").
		MustWaitLoad()
	defer page.MustClose()

	img, _ := page.MustWaitStable().Screenshot(true, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
		Clip: &proto.PageViewport{
			X:      0,
			Y:      0,
			Width:  488,
			Height: 331,
			Scale:  1,
		},
		FromSurface: true,
	})
	_ = utils.OutputFile("test.png", img)
}
