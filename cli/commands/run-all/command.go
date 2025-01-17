// Package runall provides the `run-all` command that runs a terraform command against a 'stack' by running the specified command in each subfolder.
package runall

import (
	"context"
	"sort"

	"github.com/nholuongut/terragrunt/cli/commands"

	awsproviderpatch "github.com/nholuongut/terragrunt/cli/commands/aws-provider-patch"
	graphdependencies "github.com/nholuongut/terragrunt/cli/commands/graph-dependencies"
	"github.com/nholuongut/terragrunt/cli/commands/hclfmt"
	renderjson "github.com/nholuongut/terragrunt/cli/commands/render-json"
	"github.com/nholuongut/terragrunt/cli/commands/terraform"
	terragruntinfo "github.com/nholuongut/terragrunt/cli/commands/terragrunt-info"
	validateinputs "github.com/nholuongut/terragrunt/cli/commands/validate-inputs"
	"github.com/nholuongut/terragrunt/options"
	"github.com/nholuongut/terragrunt/pkg/cli"
)

const (
	CommandName = "run-all"
)

func NewCommand(opts *options.TerragruntOptions) *cli.Command {
	return &cli.Command{
		Name:        CommandName,
		Usage:       "Run a terraform command against a 'stack' by running the specified command in each subfolder.",
		Description: "The command will recursively find terragrunt modules in the current directory tree and run the terraform command in dependency order (unless the command is destroy, in which case the command is run in reverse dependency order).",
		Subcommands: subCommands(opts).SkipRunning(),
		Action:      action(opts),
		Flags:       NewFlags(opts).Sort(),
	}
}

func NewFlags(opts *options.TerragruntOptions) cli.Flags {
	return cli.Flags{
		&cli.GenericFlag[string]{
			Name:        commands.TerragruntOutDirFlagName,
			EnvVar:      commands.TerragruntOutDirFlagEnvName,
			Destination: &opts.OutputFolder,
			Usage:       "Directory to store plan files.",
		},
		&cli.GenericFlag[string]{
			Name:        commands.TerragruntJSONOutDirFlagName,
			EnvVar:      commands.TerragruntJSONOutDirFlagEnvName,
			Destination: &opts.JSONOutputFolder,
			Usage:       "Directory to store json plan files.",
		},
	}
}

func action(opts *options.TerragruntOptions) cli.ActionFunc {
	return func(cliCtx *cli.Context) error {
		opts.RunTerragrunt = func(ctx context.Context, opts *options.TerragruntOptions) error {
			if cmd := cliCtx.Command.Subcommand(opts.TerraformCommand); cmd != nil {
				cliCtx := cliCtx.WithValue(options.ContextKey, opts)
				return cmd.Action(cliCtx)
			}

			return terraform.Run(ctx, opts)
		}

		return Run(cliCtx.Context, opts.OptionsFromContext(cliCtx))
	}
}

func subCommands(opts *options.TerragruntOptions) cli.Commands {
	cmds := cli.Commands{
		terragruntinfo.NewCommand(opts),    // terragrunt-info
		validateinputs.NewCommand(opts),    // validate-inputs
		graphdependencies.NewCommand(opts), // graph-dependencies
		hclfmt.NewCommand(opts),            // hclfmt
		renderjson.NewCommand(opts),        // render-json
		awsproviderpatch.NewCommand(opts),  // aws-provider-patch
	}

	sort.Sort(cmds)

	// add terraform command `*` after sorting to put the command at the end of the list in the help.
	cmds.Add(terraform.NewCommand(opts))

	return cmds
}
