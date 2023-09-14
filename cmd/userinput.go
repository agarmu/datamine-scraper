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
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Config struct {
	subsubquestionsOwnCodeBlocks bool
	name                         string
	projectNumber                int
	overwrite                    bool
	path                         string
	url                          *url.URL
}

var globalConfig = Config{
	subsubquestionsOwnCodeBlocks: false,
	name:                         "",
	projectNumber:                -1,
	overwrite:                    false,
	path:                         "",
	url:                          nil,
}

func isSet() bool {
	return globalConfig.projectNumber > 0 && globalConfig.name != ""
}

func getInitialUserInput() error {
	for globalConfig.name == "" {
		name, err := getValue("What is your name?", "First Last")
		if err != nil {
			return err
		} else if name != "" {
			globalConfig.name = name
			break
		}
		fmt.Println("Error: Name is required.")
	}
	for globalConfig.projectNumber <= 0 {
		number, err := getValue("What is the project number?", "0")
		if err != nil {
			return err
		} else {
			number, err := strconv.Atoi(number)
			if err != nil {
				fmt.Println("Error. Input was not a number")
				continue
			}
			if number > 0 {
				globalConfig.projectNumber = number
				break
			}
		}
		fmt.Println("Error: Project Number is required.")
	}
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	dashConnectedName := strings.Join(strings.Split(strings.ToLower(globalConfig.name), " "), "-")
	filename := fmt.Sprintf("%s-project%02d.ipynb", dashConnectedName, globalConfig.projectNumber)
	defaultPath := filepath.Join(dir, filename)
	for globalConfig.path == "" {
		resp, err := getValue("Where would you like to store this file?", defaultPath)
		if err != nil {
			return err
		}
		var path = strings.TrimSpace(resp)
		if path == "" {
			path = defaultPath
		}
		path, err = filepath.Abs(path)
		if err != nil {
			fmt.Println("There was an error in your path.")
			continue
		}
		if filepath.Ext(path) != ".ipynb" {
			fmt.Println("Warning: .ipynb extension not used.")
		}
		_, err = os.Open(path)
		if errors.Is(err, fs.ErrNotExist) {
			// this is a path we can use!
			globalConfig.path = path
		} else {
			fmt.Println("That path already exists! Pick another one.")
			continue
		}
	}
	return nil
}

func getValue(prompt string, placeholder string) (string, error) {
	p := tea.NewProgram(getStringModel(prompt, placeholder))
	resp, err := p.Run()
	p.Kill()
	if err != nil {
		return "", err
	}
	field, ok := resp.(StringModel)
	if !ok {
		log.Fatal("Response did not cast back.")
	}
	if field.forceExit {
		log.Fatal("Ctrl-C Exit used. Aborting...")
	}
	return strings.TrimSpace(field.textInput.Value()), nil
}

type StringModel struct {
	prompt    string
	textInput textinput.Model
	err       error
	forceExit bool
}

func getStringModel(prompt string, placeholder string) StringModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 256
	return StringModel{
		textInput: ti,
		err:       nil,
		prompt:    prompt,
		forceExit: false,
	}
}

func (m StringModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m StringModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyCtrlC:
			m.forceExit = true
			return m, tea.Quit
		case tea.KeyEsc:
			m.textInput.SetValue("")
			return m, nil
		case tea.KeyTab:
			if m.textInput.Value() == "" {
				m.textInput.SetValue(m.textInput.Placeholder)
			}
			return m, nil
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m StringModel) View() string {
	var escQuit string = "\n(esc to reset)"
	var tabDefault string = " (press tab for default value)"
	if m.textInput.Value() == "" {
		escQuit = ""
	} else {
		tabDefault = ""
	}
	return fmt.Sprintf(
		"%s%s\n\n%s\n%s",
		m.prompt,
		tabDefault,
		m.textInput.View(),
		escQuit,
	) + "\n"
}
