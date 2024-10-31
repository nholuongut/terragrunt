package graph

import (
	"context"
	"errors"

	runall "github.com/nholuongut/terragrunt/cli/commands/run-all"

	"github.com/nholuongut/terragrunt/cli/commands/terraform"
	"github.com/nholuongut/terragrunt/config"
	"github.com/nholuongut/terragrunt/util"

	"github.com/nholuongut/terragrunt/configstack"
	"github.com/nholuongut/terragrunt/options"
	"github.com/nholuongut/terragrunt/shell"
)

func Run(ctx context.Context, opts *options.TerragruntOptions) error {
	target := terraform.NewTarget(terraform.TargetPointParseConfig, graph)

	return terraform.RunWithTarget(ctx, opts, target)
}

func graph(ctx context.Context, opts *options.TerragruntOptions, cfg *config.TerragruntConfig) error {
	if cfg == nil {
		return errors.New("terragrunt was not able to render the config as json because it received no config. This is almost certainly a bug in Terragrunt. Please open an issue on github.com/nholuongut/terragrunt with this message and the contents of your terragrunt.hcl")
	}
	// consider root for graph identification passed destroy-graph-root argument
	rootDir := opts.GraphRoot

	// if destroy-graph-root is empty, use git to find top level dir.
	// may cause issues if in the same repo exist unrelated modules which will generate errors when scanning.
	if rootDir == "" {
		gitRoot, err := shell.GitTopLevelDir(ctx, opts, opts.WorkingDir)
		if err != nil {
			return err
		}

		rootDir = gitRoot
	}

	rootOptions, err := opts.Clone(rootDir)
	if err != nil {
		return err
	}

	rootOptions.WorkingDir = rootDir

	stack, err := configstack.FindStackInSubfolders(ctx, rootOptions)
	if err != nil {
		return err
	}

	dependentModules := stack.ListStackDependentModules()

	workDir := opts.WorkingDir
	modulesToInclude := dependentModules[workDir]
	// workdir to list too
	modulesToInclude = append(modulesToInclude, workDir)

	// include from stack only elements from modulesToInclude
	for _, module := range stack.Modules {
		module.FlagExcluded = true
		if util.ListContainsElement(modulesToInclude, module.Path) {
			module.FlagExcluded = false
		}
	}

	return runall.RunAllOnStack(ctx, opts, stack)
}