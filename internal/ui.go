package internal

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	tickDuration = 1 * time.Second
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
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

type refreshOTPsMsg struct {
	t time.Time
}

type UI struct {
	list    list.Model
	entries []Entry
}

func NewUI(entries []Entry) UI {
	ui := UI{
		list:    list.New(entriesToItems(entries, time.Now()), list.NewDefaultDelegate(), 0, 0),
		entries: entries,
	}
	ui.list.Title = "GoOTP"
	ui.list.SetShowPagination(false)
	return ui
}

func (ui UI) Init() tea.Cmd {
	return ui.tick()
}

func (ui UI) Update(teaMsg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := teaMsg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return ui, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		ui.list.SetSize(msg.Width-h, msg.Height-v)
	case refreshOTPsMsg:
		cmd = tea.Batch(
			ui.list.SetItems(entriesToItems(ui.entries, msg.t)),
			ui.tick(),
		)
	}

	var listCmd tea.Cmd
	ui.list, listCmd = ui.list.Update(teaMsg)
	return ui, tea.Batch(cmd, listCmd)
}

func (ui UI) View() string {
	return docStyle.Render(ui.list.View())
}

func (ui UI) tick() tea.Cmd {
	return tea.Tick(tickDuration, func(t time.Time) tea.Msg {
		return refreshOTPsMsg{t: t}
	})
}

func entriesToItems(entries []Entry, t time.Time) []list.Item {
	items := make([]list.Item, len(entries))
	for i, e := range entries {
		items[i] = item{
			title:       buildItemTitle(e.Issuer, e.Label),
			description: buildItemDescription(e, t),
		}
	}
	return items
}

func buildItemTitle(issuer, label string) string {
	if issuer != "" {
		return issuer + " - " + label
	}
	return label
}

func buildItemDescription(e Entry, t time.Time) string {
	otp, remaining, err := GenerateOTP(e, t)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%s - %d", otp, remaining)
}
