package crawler

import (
	"github.com/Masterminds/semver/v3"
	"github.com/gocolly/colly/v2"
)

type ReleaseParser interface {
	ParseRelease(e *colly.HTMLElement) *Release
}

type Release struct {
	// The value for the release.
	Version *semver.Version

	// The human-readable name for the release.
	Name string

	ArtifactURL string
}
