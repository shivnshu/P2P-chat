package peer

import (
	"github.com/marcusolsson/tui-go"
	"github.com/shivnshu/P2P-chat/common/iface"
	"log"
	"strings"
	"time"
)

func (c *Peer) startChatBox() {
	history := tui.NewVBox()

	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	c.ChatHistoryBox = history

	root := tui.NewHBox(chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	c.UIPainter = ui

	input.OnSubmit(func(e *tui.Entry) {
		c.submitHandler(e.Text())
		input.SetText("")
	})

	c.printToChat("Welcome to P2P chat!")
	c.printToChat("Your alias is " + c.Self.Alias + ".")
	c.printToChat("Usage:")
	c.printToChat("    dummy_alias <= I am dummy message. // Send message to alias 'dummy_alias'")
	c.printToChat("    ALL <= I am dummy message. // Send message to all online aliases")

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}

}

func (c *Peer) printToChat(str string) {
	box := c.ChatHistoryBox
	// Concurrent access
	c.ChatHistoryBoxLock.Lock()
	box.Append(tui.NewHBox(
		tui.NewLabel(str),
		tui.NewSpacer(),
	))
	c.UIPainter.Repaint()
	c.ChatHistoryBoxLock.Unlock()
}

func (c *Peer) submitHandler(txt string) {
	msg := strings.Split(txt, "<=")
	if len(msg) < 2 {
		c.printToChat("Entered text does not follow the sending format.")
		return
	}
	toAlias := strings.TrimSpace(msg[0])
	message := strings.TrimSpace(msg[1])
	c.printToChat(txt)
	new_msg := iface.Message{}
	new_msg.ToAlias = toAlias
	new_msg.FromAlias = c.Self.Alias
	new_msg.Msg = message
	new_msg.Time = time.Now()
	new_msg.TTL = iface.DefaultTTL
	new_msg.MD5Hash = iface.CalculateMD5Hash(new_msg)
	c.sendMessage(new_msg)
}

func (c *Peer) recvMessage(msg iface.Message) {
	// Concurrent access
	c.ReadMsgsLock.Lock()
	// Do not display same message (identified by its hash) twice
	if _, ok := c.ReadMsgs[msg.MD5Hash]; ok {
		c.ReadMsgsLock.Unlock()
		return
	}
	c.ReadMsgs[msg.MD5Hash] = true
	c.ReadMsgsLock.Unlock()

	var message string
	message = msg.FromAlias
	message += " => "
	message += msg.Msg
	c.printToChat(message)
}
