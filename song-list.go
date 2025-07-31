package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SongListModel struct {
	Songs       []Song
	CurrentSong int
	focused     bool
	error       ErrorMsg
}

type Song struct {
	Default string
	Name   string
	active bool
}

func (m SongListModel) Init() tea.Cmd {
	return nil
}

func (m *SongListModel) Focus() tea.Cmd {
	m.focused = true
	return nil
}

func (m *SongListModel) Blur() {
	m.focused = false
}

func (m SongListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k":
			if len(m.Songs) == 0 || m.CurrentSong == 0 {
				return m, nil
			}
			m.CurrentSong--
			return m, nil

		case "j":
			if len(m.Songs) == 0 || m.CurrentSong == len(m.Songs)-1 {
				return m, nil
			}
			m.CurrentSong++
			return m, nil
		
		// case "enter":
		// 	m.
			
		}
		
	case SongAddMsg:

		Songs, err := getSongs()
		if err != nil {
			m.error = ErrorMsg{
				msg: err.Error(),
			}
			return m, nil
		}
		if msg.msg == "" {
			m.Songs = Songs
			return m, nil
		}
		m.Songs = filterAnyChar(Songs, msg.msg)
		return m, nil
	}

	return m, cmd
}

func (m SongListModel) View() string {
	var Songs string

	Songs += lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Width(24).Align(lipgloss.Center).Render("Songs")

	for i, Song := range m.Songs {
		if m.focused && i == m.CurrentSong {
			Songs += fmt.Sprintf("\n %s \n", lipgloss.NewStyle().Background(lipgloss.Color("60")).Render(Song.Name))
		} else {
			Songs += fmt.Sprintf("\n %s \n", lipgloss.NewStyle().Render(Song.Name))
		}
	}

	return Songs
}
