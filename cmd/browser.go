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
	"github.com/glassechidna/awsweb/pkg/awsweb"
	"github.com/glassechidna/awsweb/pkg/awsweb/browser"
	"github.com/spf13/cobra"
)

func init() {
	var browserCmd = &cobra.Command{
		Use:   "browser",
		Short: "Open browser window at AWS web console",
		Long: `Assumes the given role and logs you into the AWS web console
		in the role's default region.`,
		Run: func(cmd *cobra.Command, args []string) {
			browserName := ""
			profile := ""
			noOpen, _ := cmd.PersistentFlags().GetBool("no-open")

			if len(args) == 1 {
				profile = args[0]
			} else if len(args) == 2 {
				browserName = args[0]
				profile = args[1]
			}

			provider := awsweb.GetProvider(profile)
			loginUrl := awsweb.GetLoginUrl(provider)

			if noOpen {
				fmt.Println(loginUrl)
			} else {
				b, _ := browserByName(browserName)
				b.Launch(loginUrl, profile)
			}
		},
	}

	browserCmd.PersistentFlags().Bool("no-open", false, "Disable opening a browser")
	RootCmd.AddCommand(browserCmd)
}

func browserByName(name string) (browser.Browser, error) {
	switch name {
	case browser.MozillaFirefoxName:
		return &browser.MozillaFirefox{}, nil
	case browser.GoogleChromeName:
		return &browser.GoogleChrome{}, nil
	case "":
		return browser.DefaultBrowser()
	default:
		return nil, nil
	}
}
