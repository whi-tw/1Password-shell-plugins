package sourcehut

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "sourcehut",
		Platform: schema.PlatformInfo{
			Name:     "sourcehut",
			Homepage: sdk.URL("https://sourcehut.com"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
		Executables: []schema.Executable{
			sourcehutCLI(),
		},
	}
}
