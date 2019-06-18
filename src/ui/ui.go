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

func layout(g *gocui.Gui) error {

	var views = []string{LeftView, RightView, BottomView}
	maxX, maxY := g.Size()
	for _, view := range views {
		x0, y0, x1, y1 := viewPositions[view].getCoordinates(maxX, maxY)
		if v, err := g.SetView(view, x0, y0, x1, y1); err != nil {
			v.SelFgColor = gocui.ColorBlack
			v.SelBgColor = gocui.ColorGreen

			v.Title = " " + view + " "
			if err != gocui.ErrUnknownView {
				return err

			}
		}
	}
	return nil

}
