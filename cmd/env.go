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
	"fmt"

	"github.com/glassechidna/awsweb/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// envCmd represents the env command

func init() {
	var envCmd = &cobra.Command{
		Use:   "env",
		Short: "Output temporary credentials for exporting to env vars",
		Long: `
		Generates environment variables that can be eval()ed in *nix terminals
		or Invoke-Expression'd in Powershell, or whatever you do in cmd.exe'`,
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			mfaSecret := viper.GetString("mfa-secret")
			shell := viper.GetString("shell")
			doEnv(profile, mfaSecret, shell)
		},
	}

	RootCmd.AddCommand(envCmd)

	envCmd.PersistentFlags().StringP("shell", "", "", "One of powershell, cmd, or bash")
	viper.BindPFlag("shell", envCmd.PersistentFlags().Lookup("shell"))
}

func doEnv(profile string, mfaSecret string, shell string) {
	creds, region := shared.GetCreds(profile, mfaSecret)
	printEnvVar("AWS_ACCESS_KEY_ID", creds.AccessKeyID, shell)
	printEnvVar("AWS_SECRET_ACCESS_KEY", creds.SecretAccessKey, shell)
	printEnvVar("AWS_SESSION_TOKEN", creds.SessionToken, shell)
	printEnvVar("AWS_DEFAULT_REGION", region, shell)
	printEnvVar("AWS_REGION", region, shell)
}

func printEnvVar(name string, value string, shell string) {
	switch shell {
	case "powershell":
		fmt.Printf("$Env:%s = \"%s\"\n", name, value)
	case "cmd":
		fmt.Printf("SET %s=%s\n", name, value)
	case "docker":
		fmt.Printf("-e %s=\"%s\" ", name, value)
	default:
		fmt.Printf("export %s=\"%s\"\n", name, value)
	}
}
