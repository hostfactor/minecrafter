package edition

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/gocolly/colly/v2"
	"github.com/hostfactor/minecrafter/crawler"
	"github.com/hostfactor/minecrafter/docker"
	"github.com/hostfactor/minecrafter/utils"
	"path"
)

var BedrockEditionBasePath = "https://minecraft.fandom.com"

type Bedrock struct {
}

func (b *Bedrock) ParseRelease(e *colly.HTMLElement) *crawler.Release {
	if e.Text == "Linux" && path.Ext(e.Attr("href")) == ".zip" {
		rawVer := utils.SemverRegex.FindString(e.Request.URL.Path)
		value := utils.SemverRegex.FindString(e.Attr("href"))
		v, err := semver.NewVersion(rawVer)
		if err != nil {
			return nil
		}
		return &crawler.Release{
			Version:     v,
			Value:       value,
			Name:        rawVer,
			ArtifactURL: e.Attr("href"),
		}
	}
	return nil
}

func (b *Bedrock) GetTagVariations() []docker.TagVariation {
	return []docker.TagVariation{
		{
			Tag: "latest",
			IsDefault: func(version *semver.Version, tag string) bool {
				return true
			},
			Skip: func(version *semver.Version, tag string) bool {
				return false
			},
		},
	}
}

func (b *Bedrock) GetVersionListURL() string {
	return BedrockEditionBasePath + "/wiki/Bedrock_Edition_version_history"
}

func (b *Bedrock) GenVersionURL(version string) string {
	return fmt.Sprintf("%s_%s", BedrockEditionBasePath+"/wiki/Bedrock_Edition", version)
}
