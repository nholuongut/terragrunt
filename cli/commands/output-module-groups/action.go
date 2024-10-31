package outputmodulegroups

import (
	"context"
	"fmt"

	"github.com/nholuongut/terragrunt/configstack"
	"github.com/nholuongut/terragrunt/options"
)

func Run(ctx context.Context, opts *options.TerragruntOptions) error {
	stack, err := configstack.FindStackInSubfolders(ctx, opts)
	if err != nil {
		return err
	}

	js, err := stack.JSONModuleDeployOrder(opts.TerraformCommand)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(opts.Writer, "%s\n", js)
	if err != nil {
		return err
	}

	return nil
}
