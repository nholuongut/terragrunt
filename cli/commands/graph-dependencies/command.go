// Package graphdependencies provides the command to print the terragrunt dependency graph to stdout.
package graphdependencies

import (
	"github.com/nholuongut/terragrunt/options"
	"github.com/nholuongut/terragrunt/pkg/cli"
)

const (
	CommandName = "graph-dependencies"
)

func NewCommand(opts *options.TerragruntOptions) *cli.Command {
	return &cli.Command{
		Name:   CommandName,
		Usage:  "Prints the terragrunt dependency graph to stdout.",
		Action: func(ctx *cli.Context) error { return Run(ctx, opts.OptionsFromContext(ctx)) },
	}
}
