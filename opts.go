package minecrafter

import "github.com/Masterminds/semver/v3"

type buildOpts struct {
	Constraint *semver.Constraints
}

type BuildOpt func(o buildOpts) *buildOpts

// WithSemverConstraint will only build releases that match a given semver constraint, for example
//   // Builds only releases that are greater than 1.18.
//   con, _ := semver.NewConstraint("> 1.18")
//   minecrafter.BuildEdition(new(edition.Java), WithSemverConstraint(con))
// See here -> https://github.com/Masterminds/semver#checking-version-constraints for more constraint expressions
func WithSemverConstraint(constraint *semver.Constraints) BuildOpt {
	return func(o buildOpts) *buildOpts {
		o.Constraint = constraint
		return &o
	}
}

type walkReleasesOpts struct {
	Constraint *semver.Constraints
}

type WalkReleasesOpt func(o walkReleasesOpts) *walkReleasesOpts

// WithWalkSemverConstraint will only list releases that match a given semver constraint, for example
//   // Lists only releases that are greater than 1.18.
//   con, _ := semver.NewConstraint("> 1.18")
//   minecrafter.WalkReleases(new(edition.Java), ..., WithWalkSemverConstraint(con))
// See here -> https://github.com/Masterminds/semver#checking-version-constraints for more constraint expressions
func WithWalkSemverConstraint(constraint *semver.Constraints) WalkReleasesOpt {
	return func(o walkReleasesOpts) *walkReleasesOpts {
		o.Constraint = constraint
		return &o
	}
}
