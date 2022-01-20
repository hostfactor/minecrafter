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

type JavaEdition struct {
}

var JavaEditionBasePath = "https://minecraft.fandom.com"

func (j *JavaEdition) ParseRelease(e *colly.HTMLElement) *crawler.Release {
	if e.Text == "Server" && path.Ext(e.Attr("href")) == ".jar" {
		return &crawler.Release{
			Version:     utils.SemverRegex.FindString(e.Request.URL.Path),
			ArtifactURL: e.Attr("href"),
		}
	}
	return nil
}

func (j *JavaEdition) GetVersionListURL() string {
	return JavaEditionBasePath + "/wiki/Java_Edition_version_history"
}

func (j *JavaEdition) GenVersionURL(version string) string {
	return fmt.Sprintf("%s_%s", JavaEditionBasePath+"/wiki/Java_Edition", version)
}

func (j *JavaEdition) GetTagVariations() []docker.TagVariation {
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

func (j *JavaEdition) java17Skip(_, _ string) bool {
	return false
}

func (j *JavaEdition) java16Skip(version, _ string) bool {
	v, err := semver.NewVersion(version)
	if err != nil {
		return false
	}

	return ver118Constraint.Check(v)
}

func (j *JavaEdition) java11Skip(version, _ string) bool {
	v, err := semver.NewVersion(version)
	if err != nil {
		return false
	}

	return ver117Constraint.Check(v) || ver118Constraint.Check(v)
}

func (j *JavaEdition) java8Skip(version, _ string) bool {
	v, err := semver.NewVersion(version)
	if err != nil {
		return false
	}

	return !ver112Constraint.Check(v) && !ver116Constraint.Check(v)
}

func (j *JavaEdition) java17IsDefault(version, _ string) bool {
	v, err := semver.NewVersion(version)
	if err != nil {
		return false
	}
	return ver118Constraint.Check(v)
}

func (j *JavaEdition) java16IsDefault(version, _ string) bool {
	v, err := semver.NewVersion(version)
	if err != nil {
		return false
	}
	return ver117Constraint.Check(v)
}

func (j *JavaEdition) java11IsDefault(version, _ string) bool {
	v, err := semver.NewVersion(version)
	if err != nil {
		return false
	}
	return ver116Constraint.Check(v)
}

func (j *JavaEdition) java8IsDefault(version, _ string) bool {
	v, err := semver.NewVersion(version)
	if err != nil {
		return false
	}
	return ver112Constraint.Check(v)
}
