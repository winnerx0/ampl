package main

import tea "github.com/charmbracelet/bubbletea"

type PlayerModel struct {
	
}

func (m PlayerModel) Init()(tea.Cmd){
	
	return nil
}

func (m PlayerModel) Update(msg tea.Msg)(tea.Model, tea.Cmd){
	return m, nil
}

func (m PlayerModel) View() string {
	return `
	___    __  _______  __ 
   /   |  /  |/  / __ \/ / 
  / /| | / /|_/ / /_/ / /  
 / ___ |/ /  / / ____/ /___
/_/  |_/_/  /_/_/   /_____/

A cool music player written in Go
                           
	`
}
