package gui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/systray"
	"github.com/vjpristy/secre-tar/internal/assets"
	"github.com/vjpristy/secre-tar/internal/message"
	"github.com/vjpristy/secre-tar/internal/network"
)

type SecretaryGUI struct {
	app         fyne.App
	mainWindow  fyne.Window
	conn        *network.Connection
	messages    *widget.List
	input       *widget.Entry
	messageData []string
	updateChan  chan struct{}
}

func NewSecretaryGUI(conn *network.Connection) *SecretaryGUI {
	a := app.New()
	w := a.NewWindow("Secretary Application")

	sg := &SecretaryGUI{
		app:         a,
		mainWindow:  w,
		conn:        conn,
		messageData: []string{},
		updateChan:  make(chan struct{}, 100),
	}

	sg.createUI()
	go sg.handleIncomingMessages()
	go sg.handleUpdates()

	return sg
}

func (sg *SecretaryGUI) createUI() {
	sg.messages = widget.NewList(
		func() int { return len(sg.messageData) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(sg.messageData[id])
		},
	)

	sg.input = widget.NewMultiLineEntry()

	sendBtn := widget.NewButton("Send", sg.sendMessage)

	content := container.NewBorder(nil, container.NewBorder(nil, nil, nil, sendBtn, sg.input), nil, nil, sg.messages)

	sg.mainWindow.SetContent(content)
}

func (sg *SecretaryGUI) handleIncomingMessages() {
	for {
		messageData, err := sg.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}
		msg, err := message.Deserialize(messageData)
		if err != nil {
			log.Println("Error deserializing message:", err)
			continue
		}
		sg.messageData = append(sg.messageData, fmt.Sprintf("%s: %s", msg.From, msg.Content))
		sg.updateChan <- struct{}{}
	}
}

func (sg *SecretaryGUI) handleUpdates() {
	for range sg.updateChan {
		sg.messages.Refresh()
	}
}

func (sg *SecretaryGUI) Run() {
	go systray.Run(sg.onReady, sg.onExit)
	sg.mainWindow.ShowAndRun()
}

func (sg *SecretaryGUI) onReady() {
	iconBytes, err := assets.GetIcon("secretary_icon.png")
	if err != nil {
		log.Printf("Failed to load secretary icon: %v", err)
	} else {
		systray.SetIcon(iconBytes)
	}
	systray.SetTitle("Secretary")
	systray.SetTooltip("Secretary Application")

	mShow := systray.AddMenuItem("Show", "Show the application")
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				sg.mainWindow.Show()
			case <-mQuit.ClickedCh:
				systray.Quit()
				sg.app.Quit()
				return
			}
		}
	}()
}

func (sg *SecretaryGUI) onExit() {
	// Cleanup
}

func (sg *SecretaryGUI) sendMessage() {
	content := sg.input.Text
	if content == "" {
		return
	}

	msg := message.NewMessage("Secretary", "Boss", content)
	data, err := msg.Serialize()
	if err != nil {
		// Handle error
		return
	}

	err = sg.conn.WriteMessage(data)
	if err != nil {
		// Handle error
		return
	}

	sg.input.SetText("")
}
