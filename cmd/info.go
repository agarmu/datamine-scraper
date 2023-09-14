/*
Copyright © 2023 Mukul Agarwal

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version   string = "dev"
	commit    string = "n/a"
	timestamp string = "n/a"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get information about the program",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("The Datamine Scraper (tdmscrape)")
		fmt.Println("\tCreated by: Mukul Agarwal <agarmukul23@gmail.com>")
		fmt.Println("\tLicense: GNU Affero General Public License, v3 or later")
		fmt.Println("\tFeedback/Issues: https://github.com/agarmu/datamine-scraper")
		fmt.Println("\tBuild Information:")
		fmt.Println("\t\tVersion:", version)
		fmt.Println("\t\tCommit:", commit)
		fmt.Println("\t\tBuild Timestamp:", timestamp)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
