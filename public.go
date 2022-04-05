package minecrafter

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/gocolly/colly/v2"
	"github.com/hostfactor/minecrafter/docker"
	"github.com/hostfactor/minecrafter/edition"
	"github.com/hostfactor/minecrafter/utils"
)

type Interface interface {
	BuildEdition(ed edition.Edition, opts ...BuildOpt) error
	BuildRelease(ed edition.Edition, version string, opts ...BuildOpt) error
	WalkReleases(ed edition.Edition, walker VersionWalker, opts ...WalkReleasesOpt) error
}

func New(registries []string) Interface {
	return &impl{
		Collector:  colly.NewCollector(colly.UserAgent("Mozilla/5.0")),
		Docker:     docker.New(),
		Registries: registries,
	}
}

type VersionWalker func(version string, element *colly.HTMLElement) error

type impl struct {
	Collector  *colly.Collector
	Docker     docker.Interface
	Registries []string
	BuiltFirst bool
}

func (i *impl) WalkReleases(ed edition.Edition, walker VersionWalker, opts ...WalkReleasesOpt) error {
	o := &walkReleasesOpts{}
	for _, v := range opts {
		o = v(*o)
	}

	i.Collector.OnHTML("table.wikitable td a[href]", func(element *colly.HTMLElement) {
		version := utils.SemverRegex.FindString(element.Text)
		if version == "" {
			return
		}

		if o.Constraint != nil {
			parsed, err := semver.NewVersion(version)
			if err != nil {
				return
			}

			if !o.Constraint.Check(parsed) {
				return
			}
		}

		err := walker(version, element)
		if err != nil {
			fmt.Println("Err:", err.Error(), "Failed to list version:", version, element.Attr("href"))
		}
	})

	return i.Collector.Visit(ed.GetVersionListURL())
}

func (i *impl) BuildEdition(ed edition.Edition, opts ...BuildOpt) error {
	o := &buildOpts{}
	for _, v := range opts {
		o = v(*o)
	}

	return i.WalkReleases(ed, func(version string, element *colly.HTMLElement) error {
		err := i.buildRelease(ed, version, false)
		if err != nil {
			fmt.Println("Err:", err.Error(), "Failed to build release:", version, element.Attr("href"))
		}
		return err
	}, WithWalkSemverConstraint(o.Constraint))
}

func (i *impl) BuildRelease(ed edition.Edition, version string, opts ...BuildOpt) error {
	o := &buildOpts{}
	for _, v := range opts {
		o = v(*o)
	}
	return i.buildRelease(ed, version, true)
}

func (i *impl) buildRelease(ed edition.Edition, version string, specificRelease bool) error {
	col := colly.NewCollector()
	col.OnHTML(`a[href].external.text`, func(e *colly.HTMLElement) {
		i.buildReleaseForElement(e, ed, specificRelease)
	})
	return col.Visit(ed.GenVersionURL(version))
}

func (i *impl) buildReleaseForElement(e *colly.HTMLElement, ed edition.Edition, specificRelease bool) {
	release := ed.ParseRelease(e)
	if release == nil {
		return
	}

	toPush := make([]string, 0, len(ed.GetTagVariations()))
	for _, v := range ed.GetTagVariations() {
		if v.Skip(release.Version, v.Tag) {
			fmt.Println("Skipping tag", v.Tag, "for version", release.Version, ": Incompatible.")
			continue
		}
		fmt.Println("Building image for version", release.Version, "download url", release.ArtifactURL, "tag", v.Tag)

		tags := make([]string, 0, len(i.Registries))
		for _, registry := range i.Registries {
			tags = append(tags, fmt.Sprintf("%s:%s-%s", registry, release.Version, v.DisplayName))
			if v.IsDefault(release.Version, v.Tag) {
				tags = append(tags, fmt.Sprintf("%s:%s", registry, release.Version))
			}
			if !specificRelease && !i.BuiltFirst {
				tags = append(tags, fmt.Sprintf("%s:latest", registry))
			}
		}
		i.BuiltFirst = true

		buildArgs := map[string]string{
			"ARTIFACT_URL": release.ArtifactURL,
			"VERSION":      release.Version.String(),
			"VERSION_URL":  e.Request.URL.String(),
			"TAG":          v.Tag,
		}

		err := i.Docker.Build(".", docker.BuildSpec{
			Tags:      tags,
			BuildArgs: buildArgs,
		})
		if err != nil {
			fmt.Println("Failed to run", err.Error())
		}
		toPush = append(toPush, tags...)
	}

	if len(toPush) == 0 {
		fmt.Println("Nothing to push. Skipping.")
		return
	}

	for _, v := range toPush {
		fmt.Println("Pushing", v)
		err := i.Docker.Push(v)
		if err != nil {
			fmt.Println("Failed to push", err.Error())
		}
	}
}
