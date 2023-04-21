package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Editor struct {
	app.Compo
}

//func (e *Editor) OnMount(ctx app.Context) {
//	ctx.Defer(func(_ app.Context) {
//		editor := app.Window().Get()
//	})
//}

func (e *Editor) Render() app.UI {
	editor := app.Div().ID("editor")
	return app.Div().Class("six", "columns").Body(editor)
}
