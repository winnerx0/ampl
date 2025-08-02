package main

import (
	"fmt"
	"os"
	"os/user"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dhowden/tag"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type PlayerModel struct {
	content  Song
	Song     string
	metadata *tag.Metadata
	height   int
	focused  bool
	isPaused bool
	ctrl     *beep.Ctrl
}

type playerMsg struct {
	metadata tag.Metadata
}

func (m *PlayerModel) startPlayBack() tea.Cmd {
	return func() tea.Msg {
		currentUser, err := user.Current()
		if err != nil {
			return ErrorMsg{
				msg: err.Error(),
			}
		}

		file, err := os.OpenFile(currentUser.HomeDir+"/Music/"+m.content.Name, os.O_RDONLY, 0o755)
		if err != nil {
			return ErrorMsg{
				msg: err.Error(),
			}
		}

		metadata, err := tag.ReadFrom(file)
		if err != nil {
			return ErrorMsg{
				msg: err.Error(),
			}
		}

		streamer, format, err := mp3.Decode(file)
		if err != nil {
			return ErrorMsg{
				msg: err.Error(),
			}
		}

		m.ctrl = &beep.Ctrl{
			Streamer: streamer,
			Paused:   false,
		}
		go func() {
			defer streamer.Close()

			done := make(chan bool)
			speaker.Init(format.SampleRate, format.SampleRate.N(time.Millisecond*100))

			speaker.Play(beep.Seq(
				m.ctrl,
				beep.Callback(func() {
					done <- true
				}),
			))

			<-done
		}()

		return playerMsg{
			metadata: metadata,
		}
	}
}

func (m PlayerModel) Init() tea.Cmd {
	return nil
}

func (m *PlayerModel) Focus() tea.Cmd {
	m.focused = true
	return nil
}

func (m *PlayerModel) Blur() {
	m.focused = false
}

func (m *PlayerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:

			return m, m.startPlayBack()
		}

		switch msg.String() {
		case "p":
			if m.ctrl != nil {
				speaker.Lock()
				m.isPaused = !m.isPaused
				m.ctrl.Paused = m.isPaused
				speaker.Unlock()
			}
			return m, nil
		}
	case playerMsg:
		m.metadata = &msg.metadata
		return m, nil

	}

	return m, nil
}

func (m PlayerModel) View() string {
	if m.content.Name == "" {
		return m.content.Default
	}

	pause := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).Render("p to pause/play")
	skip := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).Render("s to skip")
	reset := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).Render("r to reset")

	buttonLayout := lipgloss.JoinHorizontal(lipgloss.Bottom, pause, skip, reset)

	Metadata := fmt.Sprintf("%s\n\n%s", m.content.Name, m.metadata)
	SongData := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).Render(Metadata)

	SongLayout := lipgloss.JoinVertical(lipgloss.Center, SongData, buttonLayout)

	return lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center).Height(m.height - 6).Render(SongLayout)
}
