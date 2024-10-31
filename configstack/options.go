package configstack

import (
	"github.com/nholuongut/terragrunt/config"
	"github.com/nholuongut/terragrunt/config/hclparse"
)

type Option func(Stack) Stack

func WithChildTerragruntConfig(config *config.TerragruntConfig) Option {
	return func(stack Stack) Stack {
		stack.childTerragruntConfig = config
		return stack
	}
}

func WithParseOptions(parserOptions []hclparse.Option) Option {
	return func(stack Stack) Stack {
		stack.parserOptions = parserOptions
		return stack
	}
}
