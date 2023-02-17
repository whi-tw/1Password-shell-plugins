package sourcehut

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func sourcehutCLI() schema.Executable {
	return schema.Executable{
		Name:    "sourcehut CLI (~emersion/hut)", // TODO: Check if this is correct
		Runs:    []string{"hut"},
		DocsURL: sdk.URL("https://sr.ht/~emersion/hut/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.PersonalAccessToken,
			},
		},
	}
}
