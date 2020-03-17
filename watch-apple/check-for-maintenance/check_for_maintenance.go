package maintenance

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	// AppleComURL is apple.com URL
	AppleComURL = "https://www.apple.com/jp/"
)

var checkTexts = []string{
	"お待ちください",
	"しばらくしてから",
	"アップデート中です",
}

// Check checks if apple.com is under maintenance.
func Check() (bool, string) {
	doc, _ := goquery.NewDocument(AppleComURL)
	text := doc.Text()
	for _, t := range checkTexts {
		if strings.Contains(text, t) {
			return true, t
		}
	}
	return false, ""
}
