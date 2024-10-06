package main

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"math/rand"
)

const maxWidth = 80

type state int

const (
	stateFormPage state = iota
	stateCorrectAnswer
	stateWrongAnswer
)

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// appBoundaryView creates a styled boundary for the provided text.
func (m Model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}

// appErrorBoundaryView creates a styled error boundary for the provided text.
func (m Model) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars("/"),
		lipgloss.WithWhitespaceForeground(red),
	)
}

func (m Model) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += err.Error()
	}
	return s
}

const AsciiArt = `         
         ,MMM8&&&.
    _...MMMMM88&&&&..._
 .::'''MMMMM88&&&&&&'''::.
::     MMMMM88&&&&&&     ::
'::....MMMMM88&&&&&&....::'
   ''''MMMMM88&&&&''''
         'MMM8&&&'
`

func FormField() *Model {
	m := Model{width: maxWidth, state: stateFormPage}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	m.form = huh.NewForm(
		huh.NewGroup(
			// Question 1: Color of the planet
			huh.NewSelect[string]().
				Key("color").
				Title("What is the color of the planet?").
				Options(
					huh.NewOption("Blue", "Blue"),
					huh.NewOption("Red", "Red"),
					huh.NewOption("Green", "Green"),
					huh.NewOption("Yellow", "Yellow"),
				).
				Value(&m.color),

			// Question 2: Form of the planet, dynamically determined
			huh.NewSelect[string]().
				Key("formType").
				Title("What is the form of the planet?").
				Options(
					huh.NewOption("Stone", "Stone"),      // Default: placeholder to be updated
					huh.NewOption("Gas", "Gas"),          // Default: placeholder to be updated
					huh.NewOption("Can Be Both", "Both"), // Default: placeholder to be updated
				).
				Value(&m.formType),

			// Question 3: Flux value (open field input)
			huh.NewInput().
				Key("flux").
				Title("Flux of the planet?").
				Prompt("?").
				Value(&m.flux),

			// Question 4: Is the planet in orbit?
			huh.NewSelect[string]().
				Key("inOrbit").
				Title("Is the planet in orbit?").
				Options(
					huh.NewOption("Yes", "Yes"),
					huh.NewOption("No", "No"),
				).
				Value(&m.inOrbit),

			// Question 5: Is the planet is exoplanet?
			huh.NewSelect[string]().
				Key("isExoplanet").
				Title("Is the planet is Exoplanet?").
				Options(
					huh.NewOption("Yes", "Yes"),
					huh.NewOption("No", "No"),
				).
				Value(&m.inOrbit),

			huh.NewConfirm().
				Key("done").
				Title("All done?").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("Welp, finish up then")
					}
					return nil
				}).
				Affirmative("Yep").
				Negative("Wait, no"),
		),
	).
		WithWidth(45).
		WithShowHelp(false).
		WithShowErrors(false)

	return &m
}

func (m *Model) ValidateAnswers() bool {
	// Define the correct answers
	correctColorFormMapping := map[string]string{
		"Blue":   "Stone",
		"Red":    "Gas",
		"Green":  "Can Be Both",
		"Yellow": "Can Be Both",
	}

	// Validate color and form of the planet
	if correctForm, ok := correctColorFormMapping[m.color]; ok {
		if m.formType != correctForm {
			return false
		}
	} else {
		return false
	}

	// Validate exoplanet status based on orbit
	if (m.inOrbit == "Yes" && m.isExo != "Yes") || (m.inOrbit == "No" && m.isExo != "No") {
		return false
	}

	// Validate exoplanet status based on random ID range
	if m.randomID >= 0 && m.randomID <= 16 {
		if m.isExo != "Yes" {
			return false // If ID is between 0 and 16, it must be an exoplanet
		}
	} else if m.randomID > 16 && m.randomID <= 62 {
		if m.isExo != "No" {
			return false // If ID is between 17 and 62, it cannot be an exoplanet
		}
	}

	// Additional validations (e.g., flux value) can be added here if needed

	return true
}
func (m *Model) GenerateUrl() string {

	// Generate a random ID between 0 and 122
	randomID := rand.Intn(62)

	// Construct the URL with the random ID
	generatedURL := fmt.Sprintf("https://project.exoplanet.quest/game/?id=%d", randomID)
	m.generatedURL = generatedURL

	return generatedURL
}
