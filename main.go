/*
Purpose: Used to add a pre-defined Tanium tag to a Windows computer wihtout pushing the tag from Tanium
Notes:

	Tags can be added to the tags array when more are needed
	A new executable will need to be built every time this is updated. Run the following command:
		go build taniumTag.go
*/
package main

import (
	"fmt"
	"tanium-tag/tanium-tag/registry"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	prompt      string
	name        string
	choices     []Tag            // items on the to-do list
	cursor      int              // which to-do list item our cursor is pointing at
	selected    map[int]struct{} // which to-do items are selected
	appliedTags []string
}

var appliedTagsModel = model{
	prompt: "The following tags are on this machine:\n\n",
	name:   "appliedTagsModel",
}

var initialModel = model{
	prompt:   "What kind of Tag would you like to apply?\n\n",
	name:     "initialModel",
	choices:  typeOfTag,
	selected: make(map[int]struct{}),
}

var coreModel = model{
	prompt:   "Please select the core tag that applys to this machine.\n\n",
	name:     "coreModel",
	choices:  coreTags,
	selected: make(map[int]struct{}),
}

var buModel = model{
	prompt:   "Please select the Business Unit (location) tag that applys to this machine.\n\n",
	name:     "buModel",
	choices:  buTags,
	selected: make(map[int]struct{}),
}

var appModel = model{
	prompt:   "Please select which apps you would like installed on this machine.\n\n",
	name:     "appModel",
	choices:  appTags,
	selected: make(map[int]struct{}),
}

func switchModel(m model, s string) model {
	newModel := m
	if m.name == "initialModel" {
		if s == "Core PC Tagging" {
			newModel = coreModel
		} else {
			newModel = appModel
		}
	}
	if m.name == "coreModel" {
		newModel = buModel
	}

	return newModel
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

var regError error //catch error message from registry changes

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// This key should enter the selected option.
		case "y":
			if len(m.selected) >= 1 {
				if m.name == "appModel" || m.name == "buModel" {
					for k := range m.selected {
						err := registry.AddString(m.choices[k].name)
						if err != nil {
							regError = err
						}
					}
					//Let the user know that the tags have been added. Maybe display the new tags
					return appliedTagsModel, nil
				}
				i := m.cursor
				if m.name == "coreModel" {
					err := registry.AddString(m.choices[i].name)
					if err != nil {
						regError = err
					}
					return switchModel(m, m.choices[i].name), nil
				}
				return switchModel(m, m.choices[i].name), nil
			}

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else if m.name == "appModel" { //Want to be able to select multiple application tags
				m.selected[m.cursor] = struct{}{}
			} else if m.name == "appliedTagsModel" {
				break
			} else if len(m.selected) == 0 { //Only want to be able to select one option
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	option := "| Press y to continue with selected.\n"
	s := ""
	if regError != nil {
		s += "\n\n" + regError.Error()
		s += "\n\nPress q to quit. "
		return s
	}
	s += "Tanium Tagging App\n\n"
	s += m.prompt

	if m.name == "appliedTagsModel" {
		option = ""
		var err error
		m.appliedTags, err = registry.GetKeyStrings()
		if err != nil {
			regError = err
			s += "\n\n" + regError.Error()
			s += "\n\nPress q to quit. "
			return s
		}
		for _, v := range m.appliedTags {
			s += fmt.Sprintf("%s\n", v)
		}
		s += "\n\nPress q to quit. "
		return s
	}
	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		if m.name == "buModel" || m.name == "appModel" {
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.name)
		} else {
			s += fmt.Sprintf("%s [%s] %s - %s\n", cursor, checked, choice.name, choice.description)
		}
	}

	// The footer
	s += "\nPress q to quit. " + option
	//	Press y to continue with selected.\n"

	// Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
