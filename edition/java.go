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

type Java struct {
}

var JavaEditionBasePath = "https://minecraft.fandom.com"

func (j *Java) ParseRelease(e *colly.HTMLElement) *crawler.Release {
	if e.Text == "Server" && path.Ext(e.Attr("href")) == ".jar" {
		rawVer := utils.SemverRegex.FindString(e.Request.URL.Path)
		v, err := semver.NewVersion(rawVer)
		if err != nil {
			return nil
		}
		return &crawler.Release{
			Version:     v,
			Name:        rawVer,
			Value:       rawVer,
			ArtifactURL: e.Attr("href"),
		}
	}
	return nil
}

func (j *Java) GetVersionListURL() string {
	return JavaEditionBasePath + "/wiki/Java_Edition_version_history"
}

func (j *Java) GenVersionURL(version string) string {
	return fmt.Sprintf("%s_%s", JavaEditionBasePath+"/wiki/Java_Edition", version)
}

func (j *Java) GetTagVariations() []docker.TagVariation {
	return []docker.TagVariation{
		{Tag: "17-alpine", DisplayName: "java-17", Skip: j.java17Skip, IsDefault: j.java17IsDefault},
		{Tag: "16-alpine", DisplayName: "java-16", Skip: j.java16Skip, IsDefault: j.java16IsDefault},
		{Tag: "11-jre-slim", DisplayName: "java-11", Skip: j.java11Skip, IsDefault: j.java11IsDefault},
		{Tag: "8-jre-slim", DisplayName: "java-8", Skip: j.java8Skip, IsDefault: j.java8IsDefault},
	}
}

var ver118Constraint, _ = semver.NewConstraint(">= 1.18")

// Requires 16
var ver117Constraint, _ = semver.NewConstraint("< 1.18, >= 1.17")

// can use 11
var ver116Constraint, _ = semver.NewConstraint("< 1.17, >= 1.12")

// can use 8
var ver112Constraint, _ = semver.NewConstraint("< 1.12")

func (j *Java) java17Skip(_ *semver.Version, _ string) bool {
	return false
}

func (j *Java) java16Skip(version *semver.Version, _ string) bool {
	return ver118Constraint.Check(version)
}

func (j *Java) java11Skip(version *semver.Version, _ string) bool {
	return ver117Constraint.Check(version) || ver118Constraint.Check(version)
}

func (j *Java) java8Skip(version *semver.Version, _ string) bool {
	return !ver112Constraint.Check(version) && !ver116Constraint.Check(version)
}

func (j *Java) java17IsDefault(version *semver.Version, _ string) bool {
	return ver118Constraint.Check(version)
}

func (j *Java) java16IsDefault(version *semver.Version, _ string) bool {
	return ver117Constraint.Check(version)
}

func (j *Java) java11IsDefault(version *semver.Version, _ string) bool {
	return ver116Constraint.Check(version)
}

func (j *Java) java8IsDefault(version *semver.Version, _ string) bool {
	return ver112Constraint.Check(version)
}
