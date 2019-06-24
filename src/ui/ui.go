package ui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

type UI struct {
}

func New() *UI {
	return &UI{}
}

func (ui *UI) Init(g *gocui.Gui) {

	// g.Highlight = true
	// g.SelFgColor = gocui.ColorRed

	// help := NewHelpWidget("help", 1, 1, helpText)
	// status := NewStatusbarWidget("status", 1, 7, 50)
	// butdown := NewButtonWidget("butdown", 52, 7, "DOWN", statusDown(status))
	// butup := NewButtonWidget("butup", 58, 7, "UP", statusUp(status))
	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	// if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, toggleButton); err != nil {
	// 	log.Panicln(err)
	// }

	go g.Update(func(g *gocui.Gui) error {
		return updateView(g, "Hello")
	})

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func updateView(g *gocui.Gui, str string) error {

	v, err := g.View(RightView)
	if err != nil {
		return err
	}
	fmt.Fprintf(v, str)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
