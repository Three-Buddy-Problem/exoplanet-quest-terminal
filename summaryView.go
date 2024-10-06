package main

func (m Model) viewSummaryPage() string {
	// Create the next page content
	return m.styles.Base.Render("Thank you for completing the form!\n\nHere is your next step...")
}
