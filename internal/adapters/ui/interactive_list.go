package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zadewu/focus/internal/domain"
)

// focusItem wraps domain.Focus to satisfy the bubbles list.DefaultItem interface.
type focusItem struct {
	focus   domain.Focus
	current bool
}

func (i focusItem) Title() string {
	name := domain.ExtractShortName(i.focus.Name)
	if i.current {
		return CurrentMark.Render("▶") + " " + ActiveStyle.Render(name)
	}
	if i.focus.Archived {
		return ArchivedStyle.Render(name) + DimStyle.Render("  archived")
	}
	return name
}

func (i focusItem) Description() string { return "" }

func (i focusItem) FilterValue() string { return i.focus.Name }

// interactiveListModel is the Bubble Tea model for the focus list TUI.
type interactiveListModel struct {
	list     list.Model
	selected string // chosen focus name; empty means cancelled
	quitting bool
}

func (m interactiveListModel) Init() tea.Cmd { return nil }

func (m interactiveListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't intercept enter/esc while the filter input is active.
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if item, ok := m.list.SelectedItem().(focusItem); ok {
				m.selected = item.focus.Name
			}
			m.quitting = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 2)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m interactiveListModel) View() string {
	if m.quitting {
		return ""
	}
	return lipgloss.NewStyle().Padding(1, 2).Render(m.list.View())
}

// RunInteractiveList launches a scrollable, filterable TUI listing all focuses.
// Returns the name of the selected focus, or "" if the user cancelled.
func RunInteractiveList(focuses []domain.Focus, current string) (string, error) {
	items := make([]list.Item, len(focuses))
	for i, f := range focuses {
		items[i] = focusItem{focus: f, current: f.Name == current}
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	l := list.New(items, delegate, 0, 0)
	l.Title = "focus sessions  (/ filter · Enter switch · q quit)"
	l.Styles.Title = HeaderStyle

	m := interactiveListModel{list: l}
	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return "", err
	}
	return result.(interactiveListModel).selected, nil
}
