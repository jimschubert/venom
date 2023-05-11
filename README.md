# venom

All you need to know about cobra commands. Generate documentation in a variety of formats, similar to what's available in Cobra but with some simplifications and fixes.

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/jimschubert/venom?color=blue&sort=semver)
![Go Version](https://img.shields.io/github/go-mod/go-version/jimschubert/venom)
[![Apache 2.0 License](https://img.shields.io/badge/License-Apache%202.0-blue)](./LICENSE)  
[![build](https://github.com/jimschubert/venom/actions/workflows/build.yml/badge.svg)](https://github.com/jimschubert/venom/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jimschubert/venom)](https://goreportcard.com/report/github.com/jimschubert/venom)

## Features

* Documentation output for Markdown, YAML, JSON
* Customizable YAML and JSON marshaling

## TODO

* Manpage support
* Restructured Text support
* User-defined templating

## Install

```
go get -u github.com/jimschubert/venom
```

## Usage

First, initialize `venom.NewOptions()`. These options follow the builder pattern to make available options easily discoverable.
Next, initialize venom by passing your root command the above options. For example:

```go
func init() {
	opts := venom.NewOptions().
		WithFormats(venom.Yaml | venom.Json | venom.Markdown).
		WithShowHiddenCommands()
	cobra.CheckErr(venom.Initialize(rootCmd, opts))
}
```

All defined formats will be generated into the output directory, which defaults to `docs` and is configurable.


Here is a fuller example of generating only markdown: 
```go
package cmd

import (
	"github.com/jimschubert/venom"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "example",
	Short: "root command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	opts := venom.NewOptions().WithFormats(venom.Markdown)
	cobra.CheckErr(venom.Initialize(rootCmd, opts))
}
```

## Build/Test

```shell
go test -v -race -cover ./...
```


## License

This project is [licensed](./LICENSE) under Apache 2.0.
