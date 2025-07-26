package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RootModel struct {
	textinput                 textinput.Model
	SongList                  SongListModel
	width, height, focusIndex int
	Player                    PlayerModel
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
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		case "enter":
			return m, m.commit(m.textinput.Value())
		case "j", "k":
			if m.focusIndex == 1 { // Song-list has focus
				var cmd tea.Cmd
				updated, cmd := m.SongList.Update(msg)
				m.SongList = updated.(SongListModel)
				return m, cmd
			}
		case "tab":
			if m.focusIndex == 0 {
				m.textinput.Blur()
				m.SongList.Focus()
				m.focusIndex = 1
			} else {
				m.textinput.Focus()
				m.SongList.Blur()
				m.focusIndex = 0
			}
			return m, nil
		}
	case SongAddMsg:
		m.SongList.Songs = append(m.SongList.Songs, Song{Name: msg.msg})
		m.textinput.SetValue("")
		return m, nil
	}
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m RootModel) View() string {
	SongView := lipgloss.NewStyle().
		Background(lipgloss.Color("#313244")).
		Border(lipgloss.NormalBorder()).
		Width(40).
		Height(m.height - 3).
		Render(lipgloss.NewStyle().Foreground(lipgloss.Color("80")).Render(m.SongList.View()))

	PlayerView := lipgloss.NewStyle().
		Background(lipgloss.Color("#313244")).
		Border(lipgloss.NormalBorder()).
		Align(lipgloss.Center).
		Width(m.width - 45).
		Height(m.height - 6).
		Render(m.Player.View())

	inputView := lipgloss.NewStyle().
		Background(lipgloss.Color("#313244")).
		Border(lipgloss.NormalBorder()).
		Width(m.width - 45).
		Foreground(lipgloss.Color("86")).
		Render(m.textinput.View())

	layout := lipgloss.JoinHorizontal(lipgloss.Bottom, SongView, lipgloss.JoinVertical(lipgloss.Top, PlayerView, inputView))

	return fmt.Sprintf("%s", layout)

}

func (m RootModel) commit(msg string) tea.Cmd {
	return func() tea.Msg {

		return SongAddMsg{
			msg: msg,
		}
	}
}

func main() {
	ti := textinput.New()
	ti.Placeholder = "Search for your song"
	ti.Focus()
	ti.CharLimit = 400
	ti.Width = 300

	Songs := []Song{}

	currentUser, err := user.Current()
	
	if err != nil {
		fmt.Println("Error ", err)
	}
	
	err = filepath.WalkDir(filepath.Join(currentUser.HomeDir, "Music"), func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if !d.IsDir() {
			Songs = append(Songs, Song{Name: d.Name()})
		}

		return nil
	})

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
		Player:     PlayerModel{},
		focusIndex: 0,
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
