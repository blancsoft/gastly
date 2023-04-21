package components

import (
	"github.com/blancsoft/gastly/pkg/lib"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type ViewerOpt struct {
	// Any value, `object`, `Array`, primitive type, even `Map` or `Set`.
	Value app.Value

	// Name of the root value. May be boolean or string. Defaults to "root"
	RootName string

	// Color theme. Defaults to 'light'
	Theme string

	// Indent width for nested objects. Defaults to 3
	IndentWidth int

	// Whether enable clipboard feature.
	EnableClipboard bool

	// Default inspect depth for nested objects.
	// If the number is set too large, it could result in performance issues._
	// Defaults to 5
	DefaultInspectDepth int

	// Hide items after reaching the count.
	// `Array` and `Object` will be affected.
	// _If the number is set too large, it could result in performance issues._
	// Defaults to 30
	MaxDisplayLength int

	// When an integer value is assigned, arrays will be displayed in groups by count of the value.
	// Groups are displayed with bracket notation and can be expanded and collapsed by clicking on the brackets.
	// Defaults to 100
	GroupArraysAfterLength int

	// Cut off the string after reaching the count.
	// Collapsed strings are followed by an ellipsis.
	// String content can be expanded and collapsed by clicking on the string value. Defaults to 50
	CollapseStringsAfterLength int

	// Whether add quotes on keys. Defaults to true
	QuotesOnKeys bool

	// Whether display data type labels. Defaults to true
	DisplayDataTypes bool

	// Whether display the size of array and object. Defaults to true
	DisplayObjectSize bool

	// Whether to highlight updates. Defaults to false
	HighlightUpdates bool
}

type Viewer struct {
	app.Compo

	// Selector to HTML element for mountpoint
	selector string

	// JSON viewer object
	viewer app.Value

	// Configuration options
	opt ViewerOpt
}

func (v *Viewer) OnMount(ctx app.Context) {
	ctx.Defer(func(_ app.Context) {
		jsonViewer := app.Window().Get("JsonViewer")
		app.Log("JsonViewer", jsonViewer)

		v.viewer = jsonViewer.New(lib.ConvertToJSValue(v.opt))
		v.viewer.Call("render", v.selector)
	})
}

func (v *Viewer) Render() app.UI {
	viewer := app.Div().ID("viewer")
	return app.Div().Class("six", "columns").Body(viewer)
}

//func onViewerChange(this js.Value, args []js.Value) any { return nil }
//func onViewerCopy(this js.Value, args []js.Value) any   { return nil }
//func onViewerSelect(this js.Value, args []js.Value) any { return nil }

func NewViewer(v any, selector string) Viewer {
	defaultOpts := ViewerOpt{
		Value:               app.ValueOf(v),
		IndentWidth:         2,
		EnableClipboard:     false,
		DefaultInspectDepth: 5,
		QuotesOnKeys:        false,
		DisplayDataTypes:    false,
		DisplayObjectSize:   false,
		HighlightUpdates:    true,
	}
	return Viewer{opt: defaultOpts, selector: selector}
}
