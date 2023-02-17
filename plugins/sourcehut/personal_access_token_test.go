package sourcehut

import (
	"testing"
	
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)
	
func TestPersonalAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.Token: "wZgkmEhxDuuogQoN6fgpWlmGGDNt4abNM6sIRVktdIiPgogskJIfLqWF3dByEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SOURCEHUT_TOKEN": "wZgkmEhxDuuogQoN6fgpWlmGGDNt4abNM6sIRVktdIiPgogskJIfLqWF3dByEXAMPLE",
				},
			},
		},
	})
}

func TestPersonalAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"SOURCEHUT_TOKEN": "wZgkmEhxDuuogQoN6fgpWlmGGDNt4abNM6sIRVktdIiPgogskJIfLqWF3dByEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "wZgkmEhxDuuogQoN6fgpWlmGGDNt4abNM6sIRVktdIiPgogskJIfLqWF3dByEXAMPLE",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in sourcehut/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
			// 	{
			// 		Fields: map[sdk.FieldName]string{
			// 			fieldname.Token: "wZgkmEhxDuuogQoN6fgpWlmGGDNt4abNM6sIRVktdIiPgogskJIfLqWF3dByEXAMPLE",
			// 		},
			// 	},
			},
		},
	})
}
