package crawler

import (
	"github.com/gocolly/colly/v2"
)

type ReleaseParser interface {
	ParseRelease(e *colly.HTMLElement) *Release
}

type Release struct {
	Version     string
	ArtifactURL string
}
