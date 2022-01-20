package minecrafter

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/hostfactor/minecrafter/crawler"
	"github.com/hostfactor/minecrafter/docker"
	"github.com/hostfactor/minecrafter/edition"
	"github.com/hostfactor/minecrafter/utils"
)

type Interface interface {
	BuildEdition(ed edition.Edition) error
	BuildRelease(ed edition.Edition, version string) error
}

func New(registries []string) Interface {
	return &impl{
		Collector:  colly.NewCollector(colly.UserAgent("Mozilla/5.0")),
		Docker:     docker.New(),
		Registries: registries,
	}
}

type impl struct {
	Collector  *colly.Collector
	Docker     docker.Interface
	Registries []string
}

func (i *impl) BuildEdition(ed edition.Edition) error {
	i.Collector.OnHTML("table.wikitable td a[href]", func(element *colly.HTMLElement) {
		version := utils.SemverRegex.FindString(element.Text)
		if version == "" {
			return
		}

		err := i.buildRelease(element.Request, ed, version)
		if err != nil {
			panic(err.Error())
		}
	})

	return i.Collector.Visit(ed.GetVersionListURL())
}

func (i *impl) BuildRelease(ed edition.Edition, version string) error {
	return i.buildRelease(i.Collector, ed, version)
}

func (i *impl) buildRelease(v crawler.Visitor, ed edition.Edition, version string) error {
	i.Collector.OnHTML(`a[href].external.text`, func(e *colly.HTMLElement) {
		i.buildReleaseForElement(e, ed)
	})
	return v.Visit(ed.GenVersionURL(version))
}

func (i *impl) buildReleaseForElement(e *colly.HTMLElement, ed edition.Edition) {
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
		}

		buildArgs := map[string]string{
			"ARTIFACT_URL": release.ArtifactURL,
			"VERSION":      release.Version,
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
