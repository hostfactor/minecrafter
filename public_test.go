package minecrafter

import (
	"github.com/Masterminds/semver/v3"
	"github.com/gocolly/colly/v2"
	"github.com/hostfactor/minecrafter/docker"
	"github.com/hostfactor/minecrafter/edition"
	"github.com/hostfactor/minecrafter/mocks"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var javaRelease1181Page, _ = os.ReadFile("./testfiles/java_release_page_1_18_1.html")
var javaRelease1171Page, _ = os.ReadFile("./testfiles/java_release_page_1_17_1.html")
var javaReleaseListPage, _ = os.ReadFile("./testfiles/java_release_list_page.html")

type PublicTestSuite struct {
	suite.Suite
	Server *httptest.Server
	Docker *mocks.Interface

	Minecrafter Interface
}

func (p *PublicTestSuite) TearDownSuite() {
	p.Server.Close()
}

func (p *PublicTestSuite) SetupTest() {
	p.Server = newTestServer()

	p.Docker = new(mocks.Interface)
	p.Minecrafter = &impl{
		Collector:  colly.NewCollector(),
		Docker:     p.Docker,
		Registries: []string{"hfcr.io"},
	}
}

func (p *PublicTestSuite) TestBuildJavaEdition() {
	// -- Given
	//
	given := new(edition.Java)
	edition.JavaEditionBasePath = p.Server.URL
	p.Docker.On("Build", ".", docker.BuildSpec{
		Tags: []string{
			"hfcr.io:1.18.1-java-17",
			"hfcr.io:1.18.1",
			"hfcr.io:latest",
		},
		BuildArgs: map[string]string{
			"ARTIFACT_URL": "https://launcher.mojang.com/v1/objects/125e5adf40c659fd3bce3e66e67a16bb49ecc1b9/server.jar",
			"VERSION":      "1.18.1",
			"VERSION_URL":  p.Server.URL + "/wiki/Java_Edition" + "_1.18.1",
			"TAG":          "17-alpine",
		},
	}).Return(nil).Once()
	p.Docker.On("Build", ".", docker.BuildSpec{
		Tags: []string{
			"hfcr.io:1.17.1-java-17",
		},
		BuildArgs: map[string]string{
			"ARTIFACT_URL": "https://launcher.mojang.com/v1/objects/1-17-1/server.jar",
			"VERSION":      "1.17.1",
			"VERSION_URL":  p.Server.URL + "/wiki/Java_Edition" + "_1.17.1",
			"TAG":          "17-alpine",
		},
	}).Return(nil).Once()
	p.Docker.On("Build", ".", docker.BuildSpec{
		Tags: []string{
			"hfcr.io:1.17.1-java-16",
			"hfcr.io:1.17.1",
		},
		BuildArgs: map[string]string{
			"ARTIFACT_URL": "https://launcher.mojang.com/v1/objects/1-17-1/server.jar",
			"VERSION":      "1.17.1",
			"VERSION_URL":  p.Server.URL + "/wiki/Java_Edition" + "_1.17.1",
			"TAG":          "16-alpine",
		},
	}).Return(nil).Once()

	p.Docker.On("Push", "hfcr.io:1.17.1-java-17").Return(nil).Once()
	p.Docker.On("Push", "hfcr.io:1.17.1-java-16").Return(nil).Once()
	p.Docker.On("Push", "hfcr.io:1.17.1").Return(nil).Once()
	p.Docker.On("Push", "hfcr.io:1.18.1-java-17").Return(nil).Once()
	p.Docker.On("Push", "hfcr.io:1.18.1").Return(nil).Once()
	p.Docker.On("Push", "hfcr.io:latest").Return(nil).Once()

	// -- When
	//
	err := p.Minecrafter.BuildEdition(given)

	// -- Then
	//
	if p.NoError(err) {
		p.Docker.AssertExpectations(p.T())
	}
}

func (p *PublicTestSuite) TestBuildJavaEditionWithConstraint() {
	// -- Given
	//
	given := new(edition.Java)
	edition.JavaEditionBasePath = p.Server.URL
	con, _ := semver.NewConstraint("> 1.18.1")

	// -- When
	//
	err := p.Minecrafter.BuildEdition(given, WithSemverConstraint(con))

	// -- Then
	//
	if p.NoError(err) {
		p.Docker.AssertExpectations(p.T())
	}
}

func (p *PublicTestSuite) TestBuildRelease() {
	// -- Given
	//
	given := new(edition.Java)
	edition.JavaEditionBasePath = p.Server.URL
	p.Docker.On("Build", ".", docker.BuildSpec{
		Tags: []string{
			"hfcr.io:1.18.1-java-17",
			"hfcr.io:1.18.1",
		},
		BuildArgs: map[string]string{
			"ARTIFACT_URL": "https://launcher.mojang.com/v1/objects/125e5adf40c659fd3bce3e66e67a16bb49ecc1b9/server.jar",
			"VERSION":      "1.18.1",
			"VERSION_URL":  p.Server.URL + "/wiki/Java_Edition" + "_1.18.1",
			"TAG":          "17-alpine",
		},
	}).Return(nil)

	p.Docker.On("Push", "hfcr.io:1.18.1-java-17").Return(nil)
	p.Docker.On("Push", "hfcr.io:1.18.1").Return(nil)

	// -- When
	//
	err := p.Minecrafter.BuildRelease(given, "1.18.1")

	// -- Then
	//
	if p.NoError(err) {
		p.Docker.AssertExpectations(p.T())
	}
}

func (p *PublicTestSuite) TestWalkVersionsJavaEditionWithConstraint() {
	// -- Given
	//
	given := new(edition.Java)
	edition.JavaEditionBasePath = p.Server.URL
	expected := []string{
		"1.19",
	}
	con, _ := semver.NewConstraint("> 1.18.1")
	actual := make([]string, 0, len(expected))

	// -- When
	//
	err := p.Minecrafter.WalkReleases(given, func(version string, _ *colly.HTMLElement) error {
		actual = append(actual, version)
		return nil
	}, WithWalkSemverConstraint(con))

	// -- Then
	//
	if p.NoError(err) {
		p.Equal(expected, actual)
	}
}

func (p *PublicTestSuite) TestWalkVersionsJavaEdition() {
	// -- Given
	//
	given := new(edition.Java)
	edition.JavaEditionBasePath = p.Server.URL
	expected := []string{
		"1.19",
		"1.18.1",
		"1.17.1",
	}
	actual := make([]string, 0, len(expected))

	// -- When
	//
	err := p.Minecrafter.WalkReleases(given, func(version string, _ *colly.HTMLElement) error {
		actual = append(actual, version)
		return nil
	})

	// -- Then
	//
	if p.NoError(err) {
		p.Equal(expected, actual)
	}
}

var serverIndexResponse = []byte("hello world\n")

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write(serverIndexResponse)
	})

	mux.HandleFunc("/wiki/Java_Edition_1.18.1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(javaRelease1181Page)
	})

	mux.HandleFunc("/wiki/Java_Edition_1.17.1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(javaRelease1171Page)
	})

	mux.HandleFunc("/wiki/Java_Edition_version_history", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(javaReleaseListPage)
	})

	return httptest.NewServer(mux)
}

func TestPublicTestSuite(t *testing.T) {
	suite.Run(t, new(PublicTestSuite))
}
