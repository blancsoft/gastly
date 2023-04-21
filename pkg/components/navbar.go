package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type NavBar struct {
	app.Compo

	LogoUrl         string
	UpdateAvailable bool
}

func (n *NavBar) onUpdateClick(ctx app.Context, e app.Event) {
	// Reloads the page to display the modifications.
	ctx.Reload()
}

func (n *NavBar) Render() app.UI {
	navList := app.Ul().Class("navbar-list", "align-center").
		Body(
			app.Li().Class("eight", "columns").
				Body(
					app.A().Href("/").Style("display", "inline-block").
						Body(app.Img().Src("/web/images/gastly.png").Style("width", "15rem")),
				),
			app.Li().Class("four", "columns", "u-pull-right").
				Body(
					app.Label().Hidden(true).Text("Go package name"),
					app.Input().ID("pkg").Class("u-full-width", "m-0").
						Type("text").Placeholder("github.com/blancsoft/pkg"),
					// Displays an Update button when an update is available.
					app.If(n.UpdateAvailable,
						app.Button().
							Text("Update app!").
							OnClick(n.onUpdateClick),
					),
				),
		)
	nav := app.Nav().Class("navbar", "padding").
		Body(navList)
	return nav
}
