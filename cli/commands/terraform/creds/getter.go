// Package creds provides a way to obtain credentials through different providers and set them to `opts.Env`.
package creds

import (
	"context"

	"github.com/nholuongut/terragrunt/cli/commands/terraform/creds/providers"
	"github.com/nholuongut/terragrunt/options"
	"golang.org/x/exp/maps"
)

type Getter struct {
	obtainedCreds map[string]*providers.Credentials
}

func NewGetter() *Getter {
	return &Getter{
		obtainedCreds: make(map[string]*providers.Credentials),
	}
}

// ObtainAndUpdateEnvIfNecessary obtains credentials through different providers and sets them to `opts.Env`.
func (getter *Getter) ObtainAndUpdateEnvIfNecessary(ctx context.Context, opts *options.TerragruntOptions, authProviders ...providers.Provider) error {
	for _, provider := range authProviders {
		creds, err := provider.GetCredentials(ctx)
		if err != nil {
			return err
		}

		if creds == nil {
			continue
		}

		for providerName, prevCreds := range getter.obtainedCreds {
			if prevCreds.Name == creds.Name {
				opts.Logger.Warnf("%s credentials obtained using %s are overwritten by credentials obtained using %s.", creds.Name, providerName, provider.Name())
			}
		}

		getter.obtainedCreds[provider.Name()] = creds

		maps.Copy(opts.Env, creds.Envs)
	}

	return nil
}
