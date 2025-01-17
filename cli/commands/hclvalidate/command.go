// Package hclvalidate provides the `hclvalidate` command for Terragrunt.
//
// `hclvalidate` command recursively looks for hcl files in the directory tree starting at workingDir, and validates them
// based on the language style guides provided by Hashicorp. This is done using the official hcl2 library.
package hclvalidate

import (
	"github.com/nholuongut/terragrunt/options"
	"github.com/nholuongut/terragrunt/pkg/cli"
)

const (
	CommandName = "hclvalidate"

	ShowConfigPathFlagName = "terragrunt-hclvalidate-show-config-path"
	ShowConfigPathEnvName  = "TERRAGRUNT_HCLVALIDATE_SHOW_CONFIG_PATH"

	JSONOutputFlagName = "terragrunt-hclvalidate-json"
	JSONOutputEnvName  = "TERRAGRUNT_HCLVALIDATE_JSON"
)

func NewFlags(opts *Options) cli.Flags {
	return cli.Flags{
		&cli.BoolFlag{
			Name:        ShowConfigPathFlagName,
			EnvVar:      ShowConfigPathEnvName,
			Usage:       "Show a list of files with invalid configuration.",
			Destination: &opts.ShowConfigPath,
		},
		&cli.BoolFlag{
			Name:        JSONOutputFlagName,
			EnvVar:      JSONOutputEnvName,
			Destination: &opts.JSONOutput,
			Usage:       "Output the result in JSON format.",
		},
	}
}

func NewCommand(generalOpts *options.TerragruntOptions) *cli.Command {
	opts := NewOptions(generalOpts)

	return &cli.Command{
		Name:   CommandName,
		Usage:  "Find all hcl files from the config stack and validate them.",
		Flags:  NewFlags(opts).Sort(),
		Action: func(ctx *cli.Context) error { return Run(ctx, opts) },
	}
}
