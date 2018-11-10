package peer

import (
	"github.com/marcusolsson/tui-go"
	"log"
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

	c.printToChat(history, "Welcome to P2P chat!")
	c.printToChat(history, "Your alias is "+c.Self.Alias+".")
	c.printToChat(history, "Usage:")
	c.printToChat(history, "    dummy_alias <= I am dummy message. // Send message to alias 'dummy_alias'")
	c.printToChat(history, "    ALL <= I am dummy message. // Send message to all online aliases")

	input.OnSubmit(func(e *tui.Entry) {
		c.submitHandler(history, e.Text())
		input.SetText("")
	})

	root := tui.NewHBox(chat)

	ui, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}

func (c *Peer) printToChat(box *tui.Box, str string) {
	box.Append(tui.NewHBox(
		tui.NewLabel(str),
		tui.NewSpacer(),
	))
}

func (c *Peer) submitHandler(box *tui.Box, txt string) {
	c.printToChat(box, txt)
}
