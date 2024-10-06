package main

import (
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
		m.generatedURL, // You can use your own content here
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
		Render(s.StatusHeader.Render("Qr code to the game") + "\n \n" + qrCodeString + "\n \n" + m.generatedURL)

	// Status (right side)
	var status string
	{

		var (
			color    string
			formType string
			flux     string
			inOrbit  string
			isExo    string
		)

		if m.form.GetString("color") != "" {
			color = m.form.GetString("color")
			color = "\n\n" + s.StatusHeader.Render("Planet Color") + "\n" + color
		}
		if m.form.GetString("formType") != "" {
			formType = m.form.GetString("formType")
			formType = "\n\n" + s.StatusHeader.Render("Planet Form") + "\n" + formType
		}
		if m.form.GetString("flux") != "" {
			flux = m.form.GetString("flux")
			flux = "\n\n" + s.StatusHeader.Render("Flux") + "\n" + flux
		}
		if m.form.GetString("inOrbit") != "" {
			inOrbit = m.form.GetString("inOrbit")
			inOrbit = "\n\n" + s.StatusHeader.Render("Is in orbit") + "\n" + inOrbit
		}
		if m.form.GetString("isExoplanet") != "" {
			isExo = m.form.GetString("isExoplanet")
			isExo = "\n\n" + s.StatusHeader.Render("Is exoplanet") + "\n" + isExo
		}

		const statusWidth = 28
		statusMarginLeft := m.width - statusWidth - lipgloss.Width(form) - s.Status.GetMarginRight()
		status = s.Status.
			Height(lipgloss.Height(form)).
			Width(statusWidth).
			MarginLeft(statusMarginLeft).
			Render(s.StatusHeader.Render("Your Answers") +
				color +
				formType +
				flux +
				inOrbit +
				isExo)
	}

	errors := m.form.Errors()
	header := m.appBoundaryView("Exoplanet Game Quiz")
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
