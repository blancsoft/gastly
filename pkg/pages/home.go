package pages

import (
	"github.com/blancsoft/gastly/pkg/components"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Home struct {
	app.Compo

	updateAvailable  bool
	isAppInstallable bool
}

func (h *Home) OnAppUpdate(ctx app.Context) {
	h.updateAvailable = ctx.AppUpdateAvailable() // Reports that an app update is available.
}

func (h *Home) OnMount(ctx app.Context) {
	h.isAppInstallable = ctx.IsAppInstallable()
}

func (h *Home) Render() app.UI {
	root := app.Div().ID("root")
	nav := &components.NavBar{UpdateAvailable: h.updateAvailable}
	playground := &components.PlayGround{}
	return root.Body(
		nav,
		playground,
	)
}
