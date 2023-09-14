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
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	rom "github.com/brandenc40/romannumeral"
)

func toChar(i int) rune {
	return rune('A' + i)
}

func generateFile() error {
	// make a tempfile
	file, err := os.CreateTemp("", "*.md")
	defer os.Remove(file.Name())
	if err != nil {
		fmt.Println("Failed to create temporary file.")
		log.Fatal(err)
	}
	// check for pandoc
	cmd := exec.Command("pandoc", "--version")
	err = cmd.Run()
	if err != nil {
		fmt.Println("Unable to execute pandoc")
		log.Fatal(err)
	}
	w := bufio.NewWriter(file)
	fmt.Fprintln(w, `---
title: My notebook
jupyter:
nbformat: 4
nbformat_minor: 5
---`)
	// title
	fmt.Fprintf(w, `:::::: {.cell .markdown}
# Project %d -- %s

_This skeleton for this file was generated by [the TDM Scraper made by Mukul Agarwal](https://github.com/agarmu/datamine-scraper) from the contents of [this url](%s)._
::::::
:::::: {.cell .markdown}
**TA Help:** John Smith, Alice Jones

- Help with figuring out how to write a function.
    
**Collaboration:** Friend1, Friend2

- Helped figuring out how to load the dataset.
- Helped debug error with my plot.
::::::

:::::: {.cell .code}		
::::::
`, globalConfig.projectNumber, globalConfig.name, globalConfig.url.String())

	// generate questions
	for _, q := range questions {
		fmt.Fprintf(w, `
:::::: {.cell .markdown}
## %s
`, q.Header)
		if q.Desc != "" {
			fmt.Fprintf(w, `

**%s**
`, q.Desc)
		}
		fmt.Fprintf(w, "::::::")
		fmt.Fprintln(w, "")
		for i, sq := range q.Subquestions {
			fmt.Fprintf(w, `
:::::: {.cell .markdown}
**%c. %s**
`, toChar(i), sq.Header)
			if !globalConfig.subsubquestionsOwnCodeBlocks {
				fmt.Fprintf(w, "\n\n")
				for j, ssq := range sq.Subsubquestions {
					roman, err := rom.IntToString(j + 1)
					roman = strings.ToLower(roman)
					if err != nil {
						fmt.Printf("Failed to convert %d to roman.\n", j)
						log.Fatal(err)
					}
					fmt.Fprintf(w, `*%s. %s*<br/>`, roman, ssq)
				}
				fmt.Fprint(w, `
::::::

:::::: {.cell .code}		

::::::

:::::: {.cell .markdown}
Markdown notes and sentences and analysis written here.
::::::

`)
			} else {
				fmt.Fprintln(w, "::::::")
				for j, ssq := range sq.Subsubquestions {
					roman, err := rom.IntToString(j + 1)
					roman = strings.ToLower(roman)
					if err != nil {
						fmt.Printf("Failed to convert %d to roman.\n", j)
						log.Fatal(err)
					}
					fmt.Fprintf(w, `:::::: {.cell .markdown}
*%s. %s*
::::::

:::::: {.cell .code}		

::::::

:::::: {.cell .markdown}
Markdown notes and sentences and analysis written here.
::::::
`, roman, ssq)
				}
			}
		}
	}
	fmt.Fprint(w, `:::::: {.cell .markdown}
## Pledge

By submitting this work I hereby pledge that this is my own, personal work. I've acknowledged in the designated place at the top of this file all sources that I used to complete said work, including but not limited to: online resources, books, and electronic communications. I've noted all collaboration with fellow students and/or TA's. I did not copy or plagiarize another's work.

> As a Boilermaker pursuing academic excellence, I pledge to be honest and true in all that I do. Accountable together – We are Purdue.
::::::
`)
	err = w.Flush()
	if err != nil {
		fmt.Println("Unable to flush writer.")
		log.Fatal(err)
	}
	cmd = exec.Command("pandoc", file.Name(), "--to", "ipynb", "--output", globalConfig.path)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}