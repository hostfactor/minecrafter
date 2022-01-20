package edition

import (
	"github.com/hostfactor/minecrafter/crawler"
	"github.com/hostfactor/minecrafter/docker"
)

type Edition interface {
	crawler.ReleaseParser
	GetTagVariations() []docker.TagVariation
	GetVersionListURL() string
	GenVersionURL(version string) string
}
