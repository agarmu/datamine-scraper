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
package main

import (
	"fmt"

	"github.com/agarmu/datamine-scraper/cmd"
	latest "github.com/tcnksm/go-latest"
)

var myVersion = "1.1.1"

func main() {
	// check for version
	githubTag := &latest.GithubTag{
		Owner:             "agarmu",
		Repository:        "datamine-scraper",
		FixVersionStrFunc: latest.DeleteFrontV(),
	}
	res, err := latest.Check(githubTag, myVersion)
	if err != nil {
		fmt.Printf("ERROR: Could not check whether I am up-to-date at %s\n", err)
	} else if res.Outdated {
		fmt.Printf("ERROR: You have version %s installed, but the current version is %s. Please update your version of this program.\n", myVersion, res.Current)
	} else {
		cmd.Execute()
	}
}
