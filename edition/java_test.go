package edition

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type JavaTestSuite struct {
	suite.Suite
}

type tagVariationTest struct {
	Version       string
	ShouldSkip    bool
	ShouldDefault bool
}

func (j *JavaTestSuite) TestGetTagVariationsJava16() {
	// -- Given
	//
	given := new(JavaEdition)

	// -- When
	//
	tests := []tagVariationTest{
		{Version: "1.11.1"},
		{Version: "1.10.9"},
		{Version: "1.17", ShouldDefault: true},
		{Version: "1.16.9"},
		{Version: "1.18", ShouldSkip: true},
		{Version: "1.18.1", ShouldSkip: true},
		{Version: "1.17.9", ShouldDefault: true},
		{Version: "1.16"},
		{Version: "1.16.9"},
		{Version: "1.12"},
		{Version: "1.12.1"},
		{Version: "1.11"},
	}

	for _, v := range tests {
		j.Equal(v.ShouldSkip, given.java16Skip(v.Version, "16-alpine"), v.Version)
		j.Equal(v.ShouldDefault, given.java16IsDefault(v.Version, "16-alpine"), v.Version)
	}
}

func (j *JavaTestSuite) TestGetTagVariationsJava11() {
	// -- Given
	//
	given := new(JavaEdition)

	// -- When
	//
	tests := []tagVariationTest{
		{Version: "1.11.1"},
		{Version: "1.10.9"},
		{Version: "1.17", ShouldSkip: true},
		{Version: "1.16.9", ShouldDefault: true},
		{Version: "1.18", ShouldSkip: true},
		{Version: "1.18.1", ShouldSkip: true},
		{Version: "1.17.9", ShouldSkip: true},
		{Version: "1.16", ShouldDefault: true},
		{Version: "1.16.9", ShouldDefault: true},
		{Version: "1.12", ShouldDefault: true},
		{Version: "1.12.1", ShouldDefault: true},
		{Version: "1.11"},
	}

	for _, v := range tests {
		j.Equal(v.ShouldSkip, given.java11Skip(v.Version, "11-jre-slim"), v.Version)
		j.Equal(v.ShouldDefault, given.java11IsDefault(v.Version, "11-jre-slim"), v.Version)
	}
}

func (j *JavaTestSuite) TestGetTagVariationsJava8() {
	// -- Given
	//
	given := new(JavaEdition)

	// -- When
	//
	tests := []tagVariationTest{
		{Version: "1.11.1", ShouldDefault: true},
		{Version: "1.10.9", ShouldDefault: true},
		{Version: "1.17", ShouldSkip: true},
		{Version: "1.16.9"},
		{Version: "1.18", ShouldSkip: true},
		{Version: "1.18.1", ShouldSkip: true},
		{Version: "1.17.9", ShouldSkip: true},
		{Version: "1.16"},
		{Version: "1.16.9"},
		{Version: "1.12"},
		{Version: "1.12.1"},
		{Version: "1.11", ShouldDefault: true},
	}

	for _, v := range tests {
		j.Equal(v.ShouldSkip, given.java8Skip(v.Version, "8-jre-slim"), v.Version)
		j.Equal(v.ShouldDefault, given.java8IsDefault(v.Version, "8-jre-slim"), v.Version)
	}
}

func TestJavaTestSuite(t *testing.T) {
	suite.Run(t, new(JavaTestSuite))
}
