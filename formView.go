package main

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/mdp/qrterminal/v3"
	"strings"
)

func (m Model) viewFormPage() string {
	s := m.styles

	// Form (left side)
	v := strings.TrimSuffix(m.form.View(), "\n\n")
	form := m.lg.NewStyle().Margin(1, 0).Render(v)

	// Generate the QR code string
	qrBuilder := &strings.Builder{}
	qrterminal.GenerateWithConfig(
		"https://github.com/mdp/qrterminal", // You can use your own content here
		qrterminal.Config{
			Level:     qrterminal.L, // Low error correction for smaller QR code
			Writer:    qrBuilder,    // Write QR code output to strings.Builder
			BlackChar: qrterminal.WHITE,
			WhiteChar: qrterminal.BLACK,
			QuietZone: 1,
		},
	)
	qrCodeString := qrBuilder.String()

	qrStyled := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).           // Add double-line border around the QR code
		BorderForeground(lipgloss.AdaptiveColor(indigo)). // Add some space below the QR code
		Render(s.StatusHeader.Render("Qr code to the game") + "\n \n" + qrCodeString)

	// Status (right side)
	var status string
	{

		var (
			answer1 string
			answer2 string
			answer3 string
		)

		if m.form.GetString("answer1") != "" {
			answer1 = m.form.GetString("answer1")
			answer1 = "\n\n" + s.StatusHeader.Render("Answer 1") + "\n" + answer1
		}
		if m.form.GetString("answer2") != "" {
			answer2 = m.form.GetString("answer2")
			answer2 = "\n\n" + s.StatusHeader.Render("Answer 2") + "\n" + answer2
		}
		if m.form.GetString("answer3") != "" {
			answer3 = m.form.GetString("answer3")
			answer3 = "\n\n" + s.StatusHeader.Render("Answer 3") + "\n" + answer3
		}

		const statusWidth = 28
		statusMarginLeft := m.width - statusWidth - lipgloss.Width(form) - s.Status.GetMarginRight()
		status = s.Status.
			Height(lipgloss.Height(form)).
			Width(statusWidth).
			MarginLeft(statusMarginLeft).
			Render(s.StatusHeader.Render("Your Answers") +
				answer1 +
				answer2 +
				answer3)
	}

	errors := m.form.Errors()
	header := m.appBoundaryView("Charm Employment Application")
	if len(errors) > 0 {
		header = m.appErrorBoundaryView(m.errorView())
	}
	body := lipgloss.JoinHorizontal(lipgloss.Top, form, status)

	footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))
	if len(errors) > 0 {
		footer = m.appErrorBoundaryView("")
	}

	return s.Base.Render(qrStyled + "\n" + header + "\n" + body + "\n\n" + footer)
}

func FormField() Model {
	m := Model{width: maxWidth, state: stateFormPage}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("answer1").
				Title("What’s the paramater1").
				Prompt("?").
				Value(&m.answer1),

			huh.NewInput().
				Key("answer2").
				Title("What’s for lunch?").
				Prompt("?").
				Value(&m.answer2),

			huh.NewInput().
				Key("answer3").
				Title("What’s for lunch?").
				Prompt("?").
				Value(&m.answer3),

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
	return m
}
