package sourcehut

import (
	"context"
	"os"
	"path/filepath"

	"git.sr.ht/~emersion/go-scfg"
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

// TODO: Make this work (load from config file)
func sourcehutConfigFile(in sdk.ProvisionInput) ([]byte, error) {
	return []byte{}, nil
}

func getSourcehutConfigFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "hut", "config")
}

// TODO: Check if the platform stores the Personal Access Token in a local config file, and if so,
// implement the function below to add support for importing it.
func TrysourcehutConfigFile() sdk.Importer {
	configFilePath, err := getSourcehutConfigFilePath()
	if err != nil {
		return nil
	}
	return importer.TryFile(configFilePath, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		rootBlock, err := scfg.Load(configFilePath)
		if err != nil {
			return
		}
		var config Config
		instanceNames := make(map[string]struct{})
		for _, instanceDir := range rootBlock.GetAll("instance") {
			instance := &InstanceConfig{
				Origins: make(map[string]string),
			}

			if err := instanceDir.ParseParams(&instance.Name); err != nil {
				return
			}

			if _, ok := instanceNames[instance.Name]; ok {
				return
			}
			instanceNames[instance.Name] = struct{}{}

			if dir := instanceDir.Children.Get("access-token"); dir != nil {
				if err := dir.ParseParams(&instance.AccessToken); err != nil {
					return
				}
			}
			if dir := instanceDir.Children.Get("access-token-cmd"); dir != nil {
				if len(dir.Params) == 0 {
					return
				}
				instance.AccessTokenCmd = dir.Params
			}
			if instance.AccessToken == "" && len(instance.AccessTokenCmd) == 0 {
				return
			}
			if instance.AccessToken != "" && len(instance.AccessTokenCmd) > 0 {
				return
			}

			for _, service := range []string{"builds", "git", "hg", "lists", "meta", "pages", "paste", "todo"} {
				serviceDir := instanceDir.Children.Get(service)
				if serviceDir == nil {
					continue
				}

				originDir := serviceDir.Children.Get("origin")
				if originDir == nil {
					continue
				}

				var origin string
				if err := originDir.ParseParams(&origin); err != nil {
					return
				}

				instance.Origins[service] = origin
			}

			cfg.Instances = append(cfg.Instances, instance)
		}

		return cfg, nil

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
	Token string
	URL   string
}

type InstanceConfig struct {
	Name string

	AccessToken    string
	AccessTokenCmd []string

	Origins map[string]string
}
