package ngrok

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/hashicorp/go-version"
	"gopkg.in/yaml.v3"
)

const (
	ngrokVersion      = "2"
	apiKeyYamlName    = "api_key"
	authTokenYamlName = "authtoken"
	versionYamlName   = "version"
)

type fileProvisioner struct {
}

func ngrokProvisioner() sdk.Provisioner {
	cmd := exec.Command("ngrok", "--version")
	versionByte, err := cmd.Output()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Example: "ngrok version 3.1.1\n" to "3.1.1\n"
	versionString := strings.TrimPrefix(string(versionByte), "ngrok version ")

	// Example: "3.1.1\n" to "3.1.1"
	versionString = strings.Trim(versionString, "\n")

	// NGROK_API_KEY is supported only from ngrok 3.2.1
	// NGROK_AUTH_TOKEN has already been supported
	requiredVersion, err := version.NewVersion("3.2.1")
	if err != nil {
		log.Fatalf(err.Error())
	}

	currentVersion, err := version.NewVersion(versionString)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// If the current ngrok CLI version is 3.2.1 or higher,
	// use environment variables to provision the Shell Plugin credentials
	if currentVersion.GreaterThan(requiredVersion) || currentVersion.Equal(requiredVersion) {
		return provision.EnvVars(defaultEnvVarMapping)
	}

	// If the current ngrok CLI version less than 3.2.1,
	// use configuration files to provision the credentials
	return fileProvisioner{}
}

func (f fileProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	provisionedConfigFilePath := filepath.Join(in.TempDir, "config.yml")
	config := make(map[string]interface{})

	var existingConfigFilePath string
	// use default locations, depending on the OS
	switch runtime.GOOS {
	case "darwin":
		existingConfigFilePath = filepath.Join(in.HomeDir, "/Library/Application Support/ngrok/ngrok.yml")
	case "linux":
		existingConfigFilePath = filepath.Join(in.HomeDir, "/.config/ngrok/ngrok.yml")
	}

	flagConfigFilePath, newArgs := getConfigValueAndNewArgs(out.CommandLine, provisionedConfigFilePath)
	if flagConfigFilePath != "" {

		existingConfigFilePath = flagConfigFilePath
	}
	out.CommandLine = newArgs

	existingContents, err := os.ReadFile(existingConfigFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			out.AddError(err)
			return
		}
	} else {
		if err := importer.FileContents(existingContents).ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}
	}

	config[authTokenYamlName] = in.ItemFields[fieldname.Authtoken]
	config[apiKeyYamlName] = in.ItemFields[fieldname.APIKey]
	config[versionYamlName] = ngrokVersion

	newContents, err := yaml.Marshal(&config)
	if err != nil {
		out.AddError(err)
		return
	}

	out.AddSecretFile(provisionedConfigFilePath, newContents)
}

// getConfigValueAndNewArgs returns the value of the original config flag if specified, and the arguments with the file path replaced by the newFilePath.
func getConfigValueAndNewArgs(args []string, newFilePath string) (string, []string) {
	for i, arg := range args {
		if arg == "--config" {
			if i+1 < len(args) {
				existingFilePath := args[i+1]
				args[i+1] = newFilePath
				return existingFilePath, args
			} else {
				return "", append(args, newFilePath)
			}

		} else if strings.HasPrefix(arg, "--config=") {
			existingFilePath := strings.TrimPrefix(arg, "--config=")
			args[i] = fmt.Sprintf("--config=%s", newFilePath)
			return existingFilePath, args
		}
	}
	return "", append(args, fmt.Sprintf("--config=%s", newFilePath))
}

func (f fileProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// nothing to do here: files get deleted automatically by 1Password CLI
}

func (f fileProvisioner) Description() string {
	return "Config file aware provisioner. It will first check if an already existing config file is present."
}
