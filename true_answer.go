package main

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewTrueAnswer() string {
	s := m.styles
	message := "Oops! Some of your answers were incorrect. \n\n Press 'CTRL + C and write 'ssh localhost -p 23236 ' to console to restart the game."

	return s.Base.Render(
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.AdaptiveColor(indigo)).
			Padding(1, 2).
			Render(message + "\n\n" + AsciiArt),
	)
}
