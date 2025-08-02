package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RootModel struct {
	textinput                 textinput.Model
	SongList                  SongListModel
	width, height, focusIndex int
	Player                    *PlayerModel
	Error                     ErrorMsg
}

func (m RootModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.Player.focused = true
		updated, _ := m.Player.Update(msg)
		m.Player = updated.(*PlayerModel)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		case "enter":
			// show song information in player
			if m.focusIndex == 1 {
				m.focusIndex = 2
				m.SongList.Blur()
				songlistupdate, cmd := m.SongList.Update(msg)

				m.SongList = songlistupdate.(SongListModel)

				m.Player.content = m.SongList.Songs[m.SongList.CurrentSong]
				m.Player.height = m.height
				updated, cmd := m.Player.Update(msg)
				m.Player.focused = true
				m.Player = updated.(*PlayerModel)
				return m, cmd
			}
			return m, m.commit(m.textinput.Value())
		case "j", "k":
			if m.focusIndex == 1 { // Song-list has focus
				var cmd tea.Cmd
				updated, cmd := m.SongList.Update(msg)
				m.SongList = updated.(SongListModel)
				return m, cmd
			}
		case "p":
			updated, cmd := m.Player.Update(msg)
			m.Player.focused = true
			m.Player = updated.(*PlayerModel)

			return m, cmd

		case "s":

			if m.focusIndex == 2 {
				if len(m.SongList.Songs) == m.SongList.CurrentSong+1 {
					m.SongList.CurrentSong = 0
				} else {
					m.SongList.CurrentSong = m.SongList.CurrentSong + 1
				}
				m.Player.content.Name = m.SongList.Songs[m.SongList.CurrentSong].Name
				return m, m.Player.startPlayBack()

			}
			return m, nil
		case "tab":
			if m.focusIndex == 0 {
				m.textinput.Blur()
				m.SongList.Focus()
				m.focusIndex = 1
			} else if m.focusIndex == 1 {
				m.textinput.Focus()
				m.SongList.Blur()
				m.focusIndex = 0
			} else if m.focusIndex == 2 {
				m.focusIndex = 1
			}
			return m, nil
		}
	case SongAddMsg:

		updated, cmd := m.SongList.Update(msg)
		m.SongList = updated.(SongListModel)
		return m, cmd
	}
	m.textinput, cmd = m.textinput.Update(msg)

	return m, cmd
}

func (m RootModel) View() string {
	SongView := lipgloss.NewStyle().
		Background(lipgloss.Color("#1c1c2b")).
		Width(30).
		Height(m.height-3).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		Render(lipgloss.NewStyle().Foreground(lipgloss.Color("80")).Render(m.SongList.View()))

	PlayerView := lipgloss.NewStyle().
		Background(lipgloss.Color("#1c1c2b")).
		Border(lipgloss.NormalBorder()).
		Align(lipgloss.Center).
		Width(m.width - 35).
		Height(m.height - 6).
		Render(m.Player.View())

	inputView := lipgloss.NewStyle().
		Background(lipgloss.Color("#1c1c2b")).
		Border(lipgloss.NormalBorder()).
		Width(m.width - 35).
		Foreground(lipgloss.Color("86")).
		Render(m.textinput.View())

	layout := lipgloss.JoinHorizontal(lipgloss.Bottom, SongView, lipgloss.JoinVertical(lipgloss.Top, PlayerView, inputView))

	return layout
}

func (m RootModel) commit(msg string) tea.Cmd {
	return func() tea.Msg {
		return SongAddMsg{
			msg: msg,
		}
	}
}

func filterAnyChar(songs []Song, chars string) []Song {
	var result []Song

	for _, song := range songs {
		if containsAnyChar(song.Name, chars) {
			result = append(result, song)
		}
	}

	return result
}

func containsAnyChar(word string, chars string) bool {
	for _, c := range chars {
		if strings.ContainsRune(word, c) {
			return true
		}
	}
	return false
}

func main() {
	ti := textinput.New()
	ti.Placeholder = "Search for your song"
	ti.Focus()
	ti.CharLimit = 400
	ti.Width = 300
	ti.SetValue("")

	Songs, err := getSongs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m := RootModel{
		textinput: ti,
		SongList: SongListModel{
			Songs:       Songs,
			CurrentSong: 0,
		},
		Player: &PlayerModel{
			content: Song{
				Default: `
  	___    __  _______  __
   /   |  /  |/  / __ \/ /
  / /| | / /|_/ / /_/ / /
 / ___ |/ /  / / ____/ /___
/_/  |_/_/  /_/_/   /_____/

A cool music player written in Go`,
			},
			isPaused: false,
		},
		focusIndex: 0,
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
