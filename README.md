# Ampl - Go Terminal Music Player 🎵

A terminal-based, interactive music player built in Go using [Bubble Tea](https://github.com/charmbracelet/bubbletea) for UI and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for styling. Navigate your song library, search for songs, and play music right from your terminal!

---

## Features

* **Interactive terminal UI** with smooth navigation
* **Search songs** by typing keywords
* **Keyboard navigation**:

  * `Tab` – Switch focus between search input and song list
  * `Enter` – Select a song from the list
  * `j/k` – Move up/down the song list
  * `p` – Play/Pause the current song
  * `s` – Skip to the next song
  * `Esc` – Quit the application
* **Stylish UI** using Lip Gloss
* **Dynamic song list** with custom ASCII art for the player

---

## Installation

1. Clone this repository:

```bash
git clone https://github.com/winnerx0/ampl.git
cd ampl
```

2. Install dependencies:

```bash
go get github.com/charmbracelet/bubbletea
"go get github.com/charmbracelet/bubbles/textinput"
go get github.com/charmbracelet/lipgloss
```

3. Build and run the app:

```bash
go run main.go
```

---

## Usage

1. **Search**: Start typing in the search box to filter songs.
2. **Navigation**: Use `j` and `k` to move through the song list.
3. **Select Song**: Press `Enter` to play a selected song.
4. **Play/Pause**: Press `p` to toggle play and pause.
5. **Next Song**: Press `s` to skip to the next song.
6. **Focus**: Use `Tab` to switch between input box and song list.
7. **Exit**: Press `Esc` to quit the application.

## Contribution

Contributions are welcome! Feel free to submit issues or pull requests to improve the player, add real audio playback, or enhance the UI.

---

## License

This project is licensed under the MIT License.
