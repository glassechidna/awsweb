// Copyright © 2017 Aidan Steele <aidan.steele@glassechidna.com.au>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/glassechidna/awscredcache/sneakyvendor/aws-shared-defaults"
	"github.com/glassechidna/awsweb/pkg/awsweb"
	"github.com/go-ini/ini"
	"github.com/spf13/cobra"
)

func init() {
	var setCmd = &cobra.Command{
		Use:   "set",
		Short: "Update ~/.aws/config to use temporary credentials",
		Long: `
Modifies the "default" entry in ~/.aws/config and ~/.aws/credentials
to use a profile's temporary credentials.
`,
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			set(profile)
		},
	}

	RootCmd.AddCommand(setCmd)
}

func set(profile string) {
	provider := awsweb.GetProvider(profile)
	creds, _ := provider.Retrieve()

	cfgPath := shareddefaults.SharedConfigFilename()
	cfgIni, _ := ini.Load(cfgPath)

	cfgSect := cfgIni.Section("default")
	cfgSect.NewKey("region", provider.Region())

	credPath := shareddefaults.SharedCredentialsFilename()
	credIni, _ := ini.Load(credPath)

	credSect := credIni.Section("default")
	credSect.NewKey("aws_access_key_id", creds.AccessKeyID)
	credSect.NewKey("aws_secret_access_key", creds.SecretAccessKey)
	credSect.NewKey("aws_session_token", creds.SessionToken)

	cfgIni.SaveTo(cfgPath)
	credIni.SaveTo(credPath)
}
