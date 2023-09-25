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
	"log"
	"net/url"
	"os"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/spf13/cobra"
)

var existSubSubQuestions = false
var questions []Question

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tdmscrape [URL]",
	Short: "Scrapes Project information from the Data Mine Website into a Jupyter Notebook",
	Example: `
To use this program, simply pass the requisite url as an argument.

E.g., to import the url "https://the-examples-book.com/projects/current-projects/10100-2023-project01", run:

$ tdmscrape "https://the-examples-book.com/projects/current-projects/10100-2023-project01"
	
The appropriate .ipynb file will be created in your current directory.`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Optionally run one of the validators provided by cobra
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("Exactly 1 argument needed.")
		}
		if _, err := url.ParseRequestURI(args[0]); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		globalConfig.url, err = url.ParseRequestURI(args[0])
		if err != nil {
			l := log.New(os.Stderr, "", 0)
			l.Println("UNREACHABLE: URL DETECTION from Args Validator BYPASSED.")
			l.Printf("Args: %s", args)
			l.Printf("Please report this as a bug. (use tdmscrape info to get more info about reporting)")
			os.Exit(255)
		}
		questions, err = scrapeURL()
		if err != nil {
			return err
		}
		// post-processing
		postProcess()
		// user interaction
		err = getInitialUserInput()
		if err != nil {
			return err
		}
		err = generateFile()
		if err != nil {
			return err
		}
		return nil
	},
}

// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&globalConfig.overwrite, "overwrite", "o", false, "Overwrite existing notebook")
	rootCmd.Flags().BoolVarP(&globalConfig.subsubquestionsOwnCodeBlocks, "sub-sub-questions-own-blocks", "s", false, "sub-sub-questions get their own code blocks and response area")
	rootCmd.Flags().StringVarP(&globalConfig.name, "name", "n", "", "name to use for document")
	rootCmd.Flags().IntVarP(&globalConfig.projectNumber, "number", "i", -1, "project number")
}

func scrapeURL() ([]Question, error) {
	c := colly.NewCollector()
	questions := []Question{}
	// Find and visit all question sections
	c.OnHTML(".sect2", func(question *colly.HTMLElement) {
		// inside question area
		q := Question{
			Header:       "",
			Desc:         "",
			Subquestions: []SubQuestion{},
		}
		var err error
		q.Header = strings.TrimSpace(question.DOM.Find("h3").First().Text())
		if !strings.Contains(q.Header, "Question") {
			// not a question, exit
			return
		}
		// get the description, if it exists
		desc := question.DOM.Find(".paragraph strong")
		if desc.Length() > 0 {
			// there is a description, extract it
			q.Desc, err = desc.First().Html()
			if err != nil {
				log.Panic("Malformed Question Description: ", err)
			}
			q.Desc = strings.TrimSpace(q.Desc)
		}
		// get the subquestions, if they exist
		subquestions := question.DOM.Find(".olist ol")
		if subquestions.Length() < 1 {
			subquestions = question.DOM.Find(".ulist ul")
			if subquestions.Length() < 1 {
				log.Panic("Malformed Subquestions: @", q.Header)
			}
		}
		subquestions = subquestions.First().ChildrenFiltered("li")
		if subquestions.Length() < 1 {
			log.Panic("No subquestions: ", subquestions)
		}
		subquestions.Each(func(i int, s *goquery.Selection) {
			sq := SubQuestion{
				Header:          "",
				Subsubquestions: []string{},
			}
			headerParagraph := s.ChildrenFiltered("p")
			if headerParagraph.Length() != 1 {
				log.Panicln("Wrong number of subquestions in one container.", s.Text())
			}
			header, err := headerParagraph.First().Html()
			if err != nil {
				log.Panic(err)
			}
			sq.Header = strings.TrimSpace(header)
			// get subquestions, if any
			subsubquestions := s.Find(".olist ol")
			if subsubquestions.Length() < 1 {
				subsubquestions = s.Find(".ulist ul")
			}
			if subsubquestions.Length() == 1 {
				// there are subquestions. great!
				subsubquestions := subsubquestions.First().ChildrenFiltered("li")
				subsubquestions.Each(func(_ int, s *goquery.Selection) {
					ssqParagraph := s.ChildrenFiltered("p")
					if ssqParagraph.Length() != 1 {
						log.Panicln("Wrong number of subquestions in one container.", s.Text())
					}
					ssq, err := ssqParagraph.First().Html()
					if err != nil {
						log.Panic(err)
					}
					sq.Subsubquestions = append(sq.Subsubquestions, strings.TrimSpace(ssq))
				})
			}
			q.Subquestions = append(q.Subquestions, sq)
		})
		questions = append(questions, q)
	})
	err := c.Visit(globalConfig.url.String())
	return questions, err
}

type Question struct {
	Header       string
	Desc         string
	Subquestions []SubQuestion
}

type SubQuestion struct {
	Header          string
	Subsubquestions []string
}

func sanitize(question Question, domain string) (Question, error) {
	converter := md.NewConverter(domain, true, nil)
	var err error
	q := question
	q.Header, err = converter.ConvertString(q.Header)
	if err != nil {
		return q, err
	}
	for i := range q.Subquestions {
		q.Subquestions[i].Header, err = converter.ConvertString(q.Subquestions[i].Header)
		if err != nil {
			return q, err
		}
		for j := range q.Subquestions[i].Subsubquestions {
			q.Subquestions[i].Subsubquestions[j], err = converter.ConvertString(q.Subquestions[i].Subsubquestions[j])
			if err != nil {
				return q, err
			}
		}
	}
	return q, nil
}

func postProcess() error {
	domain := globalConfig.url.Scheme + "://" + globalConfig.url.Hostname()
	var err error
	for _, q := range questions {
		q, err = sanitize(q, domain)
		if err != nil {
			return err
		}
		if !existSubSubQuestions {
			for _, q := range q.Subquestions {
				if len(q.Subsubquestions) > 0 {
					existSubSubQuestions = true
					break
				}
			}
		}
	}
	return nil
}
