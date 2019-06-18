package ui

// position
type position struct {
	prc    float32
	margin int
}

func (p position) getCoordinate(max int) int {
	// value = prc * MAX + abs
	return int(p.prc*float32(max)) - p.margin
}

//viewPosition
type viewPosition struct {
	x0, y0, x1, y1 position
}

// Views
const (
	LeftView   = "left"
	RightView  = "right"
	BottomView = "bottom"
	// helpView = "help"
)

var viewPositions = map[string]viewPosition{
	LeftView: {
		position{0.0, 0},
		position{0.0, 0},
		position{0.3, 2},
		position{0.9, 2},
	},
	RightView: {
		position{0.3, 0},
		position{0.0, 0},
		position{1.0, 2},
		position{0.9, 2},
	},
	BottomView: {
		position{0.0, 0},
		position{0.89, 0},
		position{1.0, 2},
		position{1.0, 2},
	},
}

func (vp viewPosition) getCoordinates(maxX, maxY int) (int, int, int, int) {
	var x0 = vp.x0.getCoordinate(maxX)
	var y0 = vp.y0.getCoordinate(maxY)
	var x1 = vp.x1.getCoordinate(maxX)
	var y1 = vp.y1.getCoordinate(maxY)
	return x0, y0, x1, y1
}
