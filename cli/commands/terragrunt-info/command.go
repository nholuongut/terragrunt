// Package terragruntinfo provides the command to emit limited terragrunt state on stdout and exits.
package terragruntinfo

import (
	"github.com/nholuongut/terragrunt/options"
	"github.com/nholuongut/terragrunt/pkg/cli"
)

const (
	CommandName = "terragrunt-info"
)

func NewCommand(opts *options.TerragruntOptions) *cli.Command {
	return &cli.Command{
		Name:   CommandName,
		Usage:  "Emits limited terragrunt state on stdout and exits.",
		Action: func(ctx *cli.Context) error { return Run(ctx, opts.OptionsFromContext(ctx)) },
	}
}
