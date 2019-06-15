package ui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func Init(g *gocui.Gui) {

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

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func layout(g *gocui.Gui) error {

	maxX, maxY := g.Size()
	v, err := g.SetView("size", 1, 1, maxX/2-1, maxY-1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Clear()
	fmt.Fprintf(v, "%d, %d", maxX, maxY)
	return nil
}
