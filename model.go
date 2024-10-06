package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	state        state
	lg           *lipgloss.Renderer
	styles       *Styles
	form         *huh.Form
	width        int
	color        string
	formType     string
	flux         string
	inOrbit      string
	isExo        string
	randomID     int
	generatedURL string
}

func (m *Model) Init() tea.Cmd {

	m.generatedURL = m.GenerateUrl()
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return &m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// Process the main form
	if m.state == stateFormPage {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			cmds = append(cmds, cmd)
		}

		m.color = m.form.GetString("color")
		m.formType = m.form.GetString("formType")
		m.flux = m.form.GetString("flux")
		m.inOrbit = m.form.GetString("inOrbit")
		m.isExo = m.form.GetString("isExoplanet")

	}

	// If the form is completed, validate answers
	if m.form.State == huh.StateCompleted && m.state == stateFormPage {
		if m.ValidateAnswers() {
			m.state = stateCorrectAnswer
		} else {
			m.state = stateWrongAnswer
		}
	}

	return &m, tea.Batch(cmds...)
}
