package main

import (
	"fmt"
	"os"
	"os/user"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dhowden/tag"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type PlayerModel struct {
	content  Song
	Song     string
	metadata *tag.Metadata
}

type playerMsg struct {
	metadata tag.Metadata
}

func (m PlayerModel) startPlayBack() tea.Cmd {

	return func() tea.Msg {

		currentUser, err := user.Current()

		if err != nil {
			return ErrorMsg{
				msg: err.Error(),
			}
		}

		file, err := os.OpenFile(currentUser.HomeDir+"/Music/"+m.content.Name, os.O_RDONLY, 0755)

		if err != nil {
			return ErrorMsg{
				msg: err.Error(),
			}
		}
		metadata, err := tag.ReadFrom(file)

		streamer, format, err := mp3.Decode(file)

		if err != nil {
			return ErrorMsg{
				msg: err.Error(),
			}
		}

		go func() {
			defer streamer.Close()

			done := make(chan bool)
			speaker.Init(format.SampleRate, format.SampleRate.N(time.Millisecond*100))

			speaker.Play(beep.Seq(
				streamer,
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

func (m PlayerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:

			return m, m.startPlayBack()
		}
	case playerMsg:
		m.metadata = &msg.metadata
		return m, nil
	}

	return m, nil
}

func (m PlayerModel) View() string {

	if m.content.Name == "" {

		return fmt.Sprintf("%s", m.content.Default)
	}

	return fmt.Sprintf("%s", m.content.Name)
}
