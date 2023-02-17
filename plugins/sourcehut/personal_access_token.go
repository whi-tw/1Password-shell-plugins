package sourcehut

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func PersonalAccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.PersonalAccessToken,
		DocsURL:       sdk.URL("https://man.sr.ht/meta.sr.ht/oauth.md"),
		ManagementURL: sdk.URL("https://meta.sr.ht/oauth2"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to sourcehut.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 67,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.URL,
				MarkdownDescription: "sourcehut instance URL",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.TempFile(sourcehutConfigFile, provision.AddArgs("--config", "{{ .Path }}")),
		Importer: importer.TryAll(
			TrysourcehutConfigFile(),
		)}
}

func sourcehutConfigFile(in sdk.ProvisionInput) ([]byte, error) {
	config := Config{
		Token:       in.ItemFields[fieldname.Token],
		InstanceURL: in.ItemFields[fieldname.URL],
	}

	contents, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return []byte(contents), nil
}

// TODO: Check if the platform stores the Personal Access Token in a local config file, and if so,
// implement the function below to add support for importing it.
func TrysourcehutConfigFile() sdk.Importer {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil
	}
	return importer.TryFile(filepath.Join(configDir, "hut", "config"), func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		// var config Config
		// if err := contents.ToYAML(&config); err != nil {
		// 	out.AddError(err)
		// 	return
		// }

		// if config.Token == "" {
		// 	return
		// }

		// out.AddCandidate(sdk.ImportCandidate{
		// 	Fields: map[sdk.FieldName]string{
		// 		fieldname.Token: config.Token,
		// 	},
		// })
	})
}

type Config struct {
	Token       string
	InstanceURL string
}
