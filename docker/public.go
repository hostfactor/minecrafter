package docker

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"os"
	"os/exec"
)

type TagVariation struct {
	Tag         string
	DisplayName string
	IsDefault   IsDefaultFunc
	Skip        SkipFunc
}

type IsDefaultFunc func(version *semver.Version, tag string) bool

type SkipFunc func(version *semver.Version, tag string) bool

type Interface interface {
	Build(ctx string, spec BuildSpec) error
	Push(tag string) error
}

func New() Interface {
	return &impl{
		CommandRunner: new(DefaultCommandRunner),
	}
}

type impl struct {
	CommandRunner CommandRunner
}

func (i *impl) Build(ctx string, spec BuildSpec) error {
	args := []string{
		"build",
	}

	for _, v := range spec.Tags {

		args = append(args, "-t", v)
	}

	for k, v := range spec.BuildArgs {
		args = append(args, "--build-arg", fmt.Sprintf("%s=%s", k, v))
	}

	args = append(args, ctx)

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return i.CommandRunner.Run(cmd)
}

func (i *impl) Push(tag string) error {
	cmd := exec.Command("docker", "push", tag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return i.CommandRunner.Run(cmd)
}

type BuildSpec struct {
	Tags      []string
	BuildArgs map[string]string
}
