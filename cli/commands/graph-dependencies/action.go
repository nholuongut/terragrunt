package graphdependencies

import (
	"context"

	"github.com/nholuongut/terragrunt/configstack"
	"github.com/nholuongut/terragrunt/options"
)

// Run graph dependencies prints the dependency graph to stdout
func Run(ctx context.Context, opts *options.TerragruntOptions) error {
	stack, err := configstack.FindStackInSubfolders(ctx, opts)
	if err != nil {
		return err
	}

	// Exit early if the operation wanted is to get the graph
	stack.Graph(opts)

	return nil
}
