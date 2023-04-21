package components

import (
	"github.com/blancsoft/gastly/pkg/lib"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type PlayGround struct {
	app.Compo
}

func (p *PlayGround) Render() app.UI {
	playground := app.Div().Class("playground")
	editor := Editor{}
	viewer := NewViewer(lib.ConvertToJSValue(map[string]string{}), "#viewer")
	return playground.Body(&editor, &viewer)
}
