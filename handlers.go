package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/tj/go-dropbox"
	"github.com/tj/go-dropy"
	"github.com/yanzay/tbot/v2"
)

// Handle the /start command here
func (a *application) startHandler(m *tbot.Message) {
	msg := "This is a bot whose sole purpose is to play rock, paper, scissors with you.\nCommands:\n1. Use /play to play.\n2. Use /score to view current scores.\n3. Use /reset to reset scores."
	a.client.SendMessage(m.Chat.ID, msg)
}

// Handle the /play command here
func (a *application) playHandler(m *tbot.Message) {
	buttons := makeButtons()
	a.client.SendMessage(m.Chat.ID, "Pick an option:", tbot.OptInlineKeyboardMarkup(buttons))
}

// Handle the /score command here
func (a *application) scoreHandler(m *tbot.Message) {
	msg := fmt.Sprintf("Scores:\nWins: %v\nDraws: %v\nLosses: %v", a.wins, a.draws, a.losses)
	a.client.SendMessage(m.Chat.ID, msg)
}

// Handle the /reset command here
func (a *application) resetHandler(m *tbot.Message) {
	a.wins, a.draws, a.losses = 0, 0, 0
	msg := "Scores have been reset to 0."
	a.client.SendMessage(m.Chat.ID, msg)
}

// Handle buttton presses here
func (a *application) callbackHandler(cq *tbot.CallbackQuery) {
	humanMove := cq.Data
	msg := a.draw(humanMove)
	a.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	a.client.SendMessage(cq.Message.Chat.ID, msg)
}

// Default handler
func (a *application) defaultHandler(m *tbot.Message) {
	msg := fmt.Sprintf("We are going to log that in")
	a.client.SendMessage(m.Chat.ID, msg)

	if m.Photo != nil {
		photo, err := a.client.GetFile(m.Photo[0].FileID)
		if err != nil {
			log.Println(err)
			return
		}
		url := a.client.FileURL(photo)
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		out, err := os.Create(m.Document.FileName)
		if err != nil {
			log.Println(err)
			return
		}
		defer out.Close()
		io.Copy(out, resp.Body)
		a.dropboxUpload(m.Document.FileName)

	} else if m.Audio != nil {
		audio, err := a.client.GetFile(m.Audio.FileID)
		if err != nil {
			log.Println(err)
			return
		}
		url := a.client.FileURL(audio)
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		out, err := os.Create(m.Document.FileName)
		if err != nil {
			log.Println(err)
			return
		}
		defer out.Close()
		io.Copy(out, resp.Body)
	}

}

func (a *application) dropboxUpload(filename string) {
	token := os.Getenv("DROPBOX_ACCESS_TOKEN")
	client := dropy.New(dropbox.New(dropbox.NewConfig(token)))

	client.Upload(filename, strings.NewReader(filename))
}
