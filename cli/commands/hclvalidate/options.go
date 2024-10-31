package hclvalidate

import "github.com/nholuongut/terragrunt/options"

type Options struct {
	*options.TerragruntOptions

	ShowConfigPath bool
	JSONOutput     bool
}

func NewOptions(general *options.TerragruntOptions) *Options {
	return &Options{
		TerragruntOptions: general,
	}
}
