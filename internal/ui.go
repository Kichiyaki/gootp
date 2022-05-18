package internal

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	entry       Entry
	title       string
	description string
}

func (i item) Title() string {
	return i.title
}
func (i item) Description() string {
	return i.description
}
func (i item) FilterValue() string {
	return i.title
}

type UI struct {
	list list.Model
}

func NewUI(entries []Entry) UI {
	ui := UI{list: list.New(entriesToItems(entries), list.NewDefaultDelegate(), 0, 0)}
	ui.list.Title = "GoOTP"
	return ui
}

func (ui UI) Init() tea.Cmd {
	return nil
}

func (ui UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return ui, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		ui.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	ui.list, cmd = ui.list.Update(msg)
	return ui, cmd
}

func (ui UI) View() string {
	return docStyle.Render(ui.list.View())
}

func entriesToItems(entries []Entry) []list.Item {
	items := make([]list.Item, len(entries))
	for i, e := range entries {
		title := e.Label
		if e.Issuer != "" {
			title = e.Issuer + " - " + e.Label
		}
		items[i] = item{entry: e, title: title}
	}
	return items
}
