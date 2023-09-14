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
	_ "embed"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

//go:embed LICENSE.txt
var license string
var plain bool = false

// licenseCmd represents the license command
var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Prints the license",

	RunE: func(cmd *cobra.Command, args []string) error {
		if plain {
			fmt.Println(license)
			return nil
		}
		p := tea.NewProgram(
			ScrollViewModel{content: license, title: "License: GNU AGPL"},
			tea.WithAltScreen(),
			tea.WithMouseCellMotion(),
		)
		if _, err := p.Run(); err != nil {
			fmt.Println("Could not display license with fancy formatter.")
			fmt.Println("You may still be able to use --plain to print the license without formatting.")
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(licenseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// licenseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	licenseCmd.PersistentFlags().BoolVarP(&plain, "plain", "p", false, "Print plain message")
}
