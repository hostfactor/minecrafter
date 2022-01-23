package crawler

import (
	"github.com/Masterminds/semver/v3"
	"github.com/gocolly/colly/v2"
)

type ReleaseParser interface {
	ParseRelease(e *colly.HTMLElement) *Release
}

type Release struct {
	Version     *semver.Version
	ArtifactURL string
}
