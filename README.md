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
* User-defined templating

## TODO

* Manpage support
* Restructured Text support

## Install

```
go get -u github.com/jimschubert/venom
```

## Usage

Venom can be used either to add a documentation command as a child to another command, or to perform ad hoc writes of documentation directly.

The options, constructed via `venom.NewOptions()`, are the same for both use cases. These options follow the builder pattern to make available options easily discoverable.

All defined formats will be generated into the output directory, which defaults to `docs` and is configurable.

### As a child command

First, initialize `venom.NewOptions()`. 
Next, initialize venom by passing your root command the above options. For example:

```go
func init() {
	opts := venom.NewOptions().
		WithFormats(venom.Yaml | venom.Json | venom.Markdown).
		WithShowHiddenCommands()
	cobra.CheckErr(venom.Initialize(rootCmd, opts))
}
```

Here is a full example of generating only markdown: 

```go
package main

import (
	"github.com/jimschubert/venom"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "example",
	Short: "root command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	opts := venom.NewOptions().WithFormats(venom.Markdown)
	cobra.CheckErr(venom.Initialize(rootCmd, opts))
}

func main() {
	rootCmd.Execute()
}
```

If you compile and run this application, you won't see a `docs` command because it's hidden. Invoke the hidden command (e.g. `example docs`) to output the documentation. 

The `docs` subcommand which venom creates will allow the user to specify a subset of the allowed formats, to define a new output directory, and to show all hidden commands.

```
Usage:
  example docs [flags]

Flags:
      --formats strings   A comma-separated list of formats to output. Allowed: [yaml,markdown] (default [yaml,markdown])
  -h, --help              help for docs
      --out-dir string    The target output directory (default "docs")
      --show-hidden       Also show hidden commands

```

### Via Write

Suppose you want to wire this functionality up into an existing documentation command, or maybe you want to generate on 
build via a go generator utility command. You can do this by invoking the `Write` command directly.

For example:

```go
docs := venom.NewDocumentation(rootCmd, venom.NewOptions().WithFormats(venom.Markdown))
if err := venom.Write(docs)); err != nil {
	// do something with err
}
```

## Custom Templates

You can provide your own templates if the built-in templates don't suit your needs. The built-in templates are intended 
to match as closely as possible with those output by Cobra's built-in command. But, there are cases where these aren't desirable. For instance:

* markdown-driven doc sites like Docusaurus generate first-level headers if missing in markdown
* you want front-matter for an extended markdown system like Jekyll
* you want to do something unexpected like output Asciidoc by tweaking the Markdown templates
* you simply don't like the formatting

To provide custom templates, you just need to make sure your files are named exactly the same as they are under [./templates](./templates) in this repository.
Then, pass an implementation of `fs.FS` to our options. You can use go embed to read an entire directory called `your_directory` like this:

```
//go:embed your_directory/*.tmpl
var templates embed.FS

func init() {
	opts := venom.NewOptions().WithCustomTemplates(templates)
	cobra.CheckErr(venom.Initialize(rootCmd, opts))
}
```

And as long as you have both `markdown_command.tmpl` and `markdown_index.tmpl` defined under `your_directory`, you're all set!

The types definitions which are bound to these templates can be found in [./types.go](./types.go).

Command templates will be bound to a data structure matching:

```go
struct {
    Command
    Doc Documentation
}
```

The embedded `Command` allows you to interact with the command's fields directly at the top level of the template. The `Doc` 
field is the full documentation, providing you access to the root command and all child commands.

**NOTE** Not all output formats are template driven. Be sure to review [./templates](./templates).

## Build/Test

```shell
go test -v -race -cover ./...
```

## License

This project is [licensed](./LICENSE) under Apache 2.0.
