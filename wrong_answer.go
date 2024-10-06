package main

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) viewWrongAnswer() string {
	s := m.styles
	message := "Oops! Some of your answers were incorrect. \n\nPress 'CTRL + C and write 'ssh localhost -p 23236 ' to console to restart the game."

	return s.Base.Render(
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.AdaptiveColor(indigo)).
			Padding(1, 2).
			Render(message + "\n\n" + AsciiArt),
	)
}
