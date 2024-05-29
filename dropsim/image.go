package dropsim

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"
	"os"
	"sort"
	"strings"
)

type ResponseImage struct {
	Title    string
	Images   []string
	Captions []string
	Filepath string
	Content  string
}

func (r *ResponseImage) GetScreenshot(filepath string) ([]byte, error) {
	browser := rod.New().MustConnect()
	page := browser.
		MustPage(fmt.Sprintf("file:////home/frankie/workspace/github.com/phantomwaves/proto-bot/%s", filepath)).
		MustWaitLoad()
	defer page.MustClose()

	img, err := page.MustWaitStable().Screenshot(true, &proto.PageCaptureScreenshot{
		Format: proto.PageCaptureScreenshotFormatPng,
		Clip: &proto.PageViewport{
			X:      0,
			Y:      0,
			Width:  488,
			Height: 511,
			Scale:  1,
		},
		FromSurface: true,
	})
	if err != nil {
		return nil, err
	}
	err = utils.OutputFile(fmt.Sprintf("responses/%s.png", filepath), img)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (r *ResponseImage) getTemplate(path string) {
	dat, _ := os.ReadFile(path)
	r.Content = string(dat)
}

func (r *ResponseImage) setTitle(title string) {
	r.Content = strings.ReplaceAll(r.Content, "{{ title }}", title)
}

func (r *ResponseImage) setCaptions() {
	captions := strings.Join(r.Captions, "\n")
	r.Content = strings.ReplaceAll(r.Content, "{{ captions }}", captions)
}
func (r *ResponseImage) setImages() {
	images := strings.Join(r.Images, "\n")
	r.Content = strings.ReplaceAll(r.Content, "{{ images }}", images)

}

func (r *ResponseImage) MakeResponse(itemCounts map[string]int) {
	keys := make([]string, 0, len(itemCounts))
	for k := range itemCounts {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return itemCounts[keys[i]] > itemCounts[keys[j]]
	})

	for _, k := range keys {
		if itemCounts[k] > 0 {
			r.Images = append(r.Images, fmt.Sprintf("<img src=\"../images/%s.png\" alt=\"\">", k))
			if itemCounts[k] == 1 {
				r.Captions = append(r.Captions, "<div></div>")
			} else if itemCounts[k] > 10000000 {
				n := itemCounts[k] / 1000000
				r.Captions = append(r.Captions, fmt.Sprintf("<div style=\"color: #00ff80\">%dM</div>", n))
			} else if itemCounts[k] >= 100000 {
				n := itemCounts[k] / 1000
				r.Captions = append(r.Captions, fmt.Sprintf("<div style=\"color: white\">%dK</div>", n))
			} else {
				r.Captions = append(r.Captions, fmt.Sprintf("<div>%d</div>", itemCounts[k]))
			}

		}
	}
	for _ = range 88 - len(keys) {
		r.Images = append(r.Images, "<div class=\"placeholder\"></div>")
		r.Captions = append(r.Captions, "<div></div>")
	}
	r.getTemplate("response/results_template.html")
	r.setTitle(r.Title)
	r.setCaptions()
	r.setImages()
	r.Filepath = fmt.Sprintf("responses/%s.html", r.Title)
	os.WriteFile(r.Filepath, []byte(r.Content), 0666)
}
