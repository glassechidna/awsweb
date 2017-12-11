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
	"github.com/aws/aws-sdk-go/aws/credentials"
	"strings"
	"os"
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
			shell := viper.GetString("shell")
			unset, _ := cmd.Flags().GetBool("unset")

			creds := credentials.Value{}
			region := ""
			profile := ""

			if !unset {
				if len(args) > 0 { profile = args[0] }
				mfaSecret := viper.GetString("mfa-secret")
				creds, region = shared.GetCreds(profile, mfaSecret)
			}

			doEnv(creds, profile, region, shell, unset)
		},
	}

	RootCmd.AddCommand(envCmd)

	envCmd.Flags().Bool("unset", false, "Generate output to unset env vars")
	envCmd.PersistentFlags().StringP("shell", "", "", "One of powershell, cmd, or bash")
	viper.BindPFlag("shell", envCmd.PersistentFlags().Lookup("shell"))
}

func doEnv(creds credentials.Value, profile, region, shell string, unset bool) {
	printExplanation(shell)
	printEnvVar("AWS_ACCESS_KEY_ID", creds.AccessKeyID, shell, unset)
	printEnvVar("AWS_SECRET_ACCESS_KEY", creds.SecretAccessKey, shell, unset)
	printEnvVar("AWS_SESSION_TOKEN", creds.SessionToken, shell, unset)
	printEnvVar("AWS_DEFAULT_REGION", region, shell, unset)
	printEnvVar("AWS_REGION", region, shell, unset)
	printEnvVar("AWSWEB_PROFILE", profile, shell, unset)
}

func printExplanation(shell string) {
	switch shell {
	case "powershell":
		fmt.Printf(`
# The output of this command is meant to be eval'd, i.e. re-run this command:
#
# $Cmd = (awsweb env --shell powershell mycompany-prod) | Out-String
# Invoke-Expression $Cmd
`)
	case "cmd":
	case "docker":
	default:
		command := strings.Join(os.Args, " ")
		fmt.Printf(`
# The output of this command is meant to be eval'd, i.e. re-run this command:
#
# eval $(%s)
`, command)
	}
}

func printEnvVar(name, value, shell string, unset bool) {
	switch shell {
	case "powershell":
		if unset {
			fmt.Printf("Remove-Item Env:\\%s\n", name)
		} else {
			fmt.Printf("$Env:%s = \"%s\"\n", name, value)
		}
	case "cmd":
		if unset {
			fmt.Printf("SET %s=\n", name)
		} else {
			fmt.Printf("SET %s=%s\n", name, value)
		}
	case "docker":
		fmt.Printf("-e %s=\"%s\" ", name, value)
	default:
		if unset {
			fmt.Printf("unset %s\n", name)
		} else {
			fmt.Printf("export %s=\"%s\"\n", name, value)
		}
	}
}
