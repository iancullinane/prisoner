package ui

import "github.com/jroimartin/gocui"

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
