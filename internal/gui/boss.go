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

type BossGUI struct {
	app         fyne.App
	mainWindow  fyne.Window
	conn        *network.Connection
	messages    *widget.List
	input       *widget.Entry
	messageData []string
	updateChan  chan struct{}
}

func NewBossGUI(conn *network.Connection) *BossGUI {
	a := app.New()
	w := a.NewWindow("Boss Application")

	bg := &BossGUI{
		app:         a,
		mainWindow:  w,
		conn:        conn,
		messageData: []string{},
		updateChan:  make(chan struct{}, 100),
	}

	bg.createUI()
	go bg.handleIncomingMessages()
	go bg.handleUpdates()

	return bg
}

func (bg *BossGUI) createUI() {
	bg.messages = widget.NewList(
		func() int { return len(bg.messageData) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(bg.messageData[id])
		},
	)

	bg.input = widget.NewMultiLineEntry()

	sendBtn := widget.NewButton("Send", bg.sendMessage)
	yesBtn := widget.NewButton("Yes", func() { bg.sendQuickReply("Yes") })
	noBtn := widget.NewButton("No", func() { bg.sendQuickReply("No") })
	waitBtn := widget.NewButton("Wait", func() { bg.sendQuickReply("Wait") })
	comeInBtn := widget.NewButton("Come in", func() { bg.sendQuickReply("Come in") })

	buttons := container.NewHBox(yesBtn, noBtn, waitBtn, comeInBtn)
	content := container.NewBorder(nil, container.NewBorder(nil, nil, nil, sendBtn, bg.input), nil, nil, bg.messages)

	bg.mainWindow.SetContent(container.NewBorder(nil, buttons, nil, nil, content))
}

func (bg *BossGUI) handleIncomingMessages() {
	for {
		messageData, err := bg.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}
		msg, err := message.Deserialize(messageData)
		if err != nil {
			log.Println("Error deserializing message:", err)
			continue
		}
		bg.messageData = append(bg.messageData, fmt.Sprintf("%s: %s", msg.From, msg.Content))
		bg.updateChan <- struct{}{}
	}
}

func (bg *BossGUI) handleUpdates() {
	for range bg.updateChan {
		bg.messages.Refresh()
	}
}

func (bg *BossGUI) Run() {
	go systray.Run(bg.onReady, bg.onExit)
	bg.mainWindow.ShowAndRun()
}

func (bg *BossGUI) onReady() {
	iconBytes, err := assets.GetIcon("boss_icon.png")
	if err != nil {
		log.Printf("Failed to load boss icon: %v", err)
	} else {
		systray.SetIcon(iconBytes)
	}
	systray.SetTitle("Boss")
	systray.SetTooltip("Boss Application")

	mShow := systray.AddMenuItem("Show", "Show the application")
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				bg.mainWindow.Show()
			case <-mQuit.ClickedCh:
				systray.Quit()
				bg.app.Quit()
				return
			}
		}
	}()
}

func (bg *BossGUI) onExit() {
	// Cleanup
}

func (bg *BossGUI) sendMessage() {
	content := bg.input.Text
	if content == "" {
		return
	}

	msg := message.NewMessage("Boss", "Secretary", content)
	data, err := msg.Serialize()
	if err != nil {
		// Handle error
		return
	}

	err = bg.conn.WriteMessage(data)
	if err != nil {
		// Handle error
		return
	}

	bg.input.SetText("")
}

func (bg *BossGUI) sendQuickReply(reply string) {
	msg := message.NewMessage("Boss", "Secretary", reply)
	data, err := msg.Serialize()
	if err != nil {
		// Handle error
		return
	}

	err = bg.conn.WriteMessage(data)
	if err != nil {
		// Handle error
		return
	}
}
