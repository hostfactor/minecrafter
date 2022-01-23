# Minecrafter

Small Go library for building and pushing Minecraft Docker images.

## Install

```
go get github.com/hostfactor/minecrafter
```

## Example

### Build and push every version

#### Java Edition

```go
package main

import (
	"github.com/hostfactor/minecrafter"
	"github.com/hostfactor/minecrafter/edition"
	"os"
)

func main() {
	builder := minecrafter.New([]string{os.Getenv("GITHUB_REGISTRY")})

	err := builder.BuildEdition(new(edition.Java))
	if err != nil {
		panic(err.Error())
	}
}
```

#### Bedrock Edition

```go
package main

import (
	"github.com/hostfactor/minecrafter"
	"github.com/hostfactor/minecrafter/edition"
	"os"
)

func main() {
	builder := minecrafter.New([]string{os.Getenv("GITHUB_REGISTRY")})

	err := builder.BuildEdition(new(edition.BedrockEdition))
	if err != nil {
		panic(err.Error())
	}
}
```

### Build and push a specific version

#### Java Edition

```go
package main

import (
	"github.com/hostfactor/minecrafter"
	"github.com/hostfactor/minecrafter/edition"
	"os"
)

func main() {
	builder := minecrafter.New([]string{os.Getenv("GITHUB_REGISTRY")})

	err := builder.BuildRelease(new(edition.Java), "1.18.1")
	if err != nil {
		panic(err.Error())
	}
}
```

#### Bedrock Edition

```go
package main

import (
	"github.com/hostfactor/minecrafter"
	"github.com/hostfactor/minecrafter/edition"
	"os"
)

func main() {
	builder := minecrafter.New([]string{os.Getenv("GITHUB_REGISTRY")})

	err := builder.BuildRelease(new(edition.BedrockEdition), "1.18.2")
	if err != nil {
		panic(err.Error())
	}
}
```
