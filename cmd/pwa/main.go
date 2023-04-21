package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"

	"github.com/blancsoft/gastly/pkg/pages"
)

func main() {
	// Client-side code
	{
		// Define the routes
		app.Route("/", &pages.Home{})

		// Start the app
		app.RunWhenOnBrowser()
	}

	// Server-/build-side code
	{
		// Parse the flags
		build := flag.Bool("build", false, "Create static build")
		out := flag.String("out", "dist", "Out directory for static build")
		path := flag.String("path", "", "Base path for static build")
		serve := flag.Bool("serve", false, "Build and serve the frontend")
		laddr := flag.String("laddr", "localhost:8081", "Address to serve the frontend on")

		flag.Parse()

		// Define the handler
		h := &app.Handler{
			Name:         "Gastly",
			Title:        "Gastly",
			LoadingLabel: "Go AST renderer in the browser powered by WASM",
			Description:  "An abstract syntax tree (AST) parser for rendering AST tree in the browser using WASM module",
			Author:       "Chumma Umenze",

			Image:           "/web/images/gastly.png",
			BackgroundColor: "#151515",
			Icon: app.Icon{
				Default: "/web/images/gastly.png",
			},
			Keywords: []string{},

			RawHeaders: []string{
				`<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/normalize/8.0.1/normalize.min.css" integrity="sha512-NhSC1YmyruXifcj/KFRWoC561YpHpc5Jtzgvbuzx5VozKpWvQ+4nXhPdFgmx8xqexRcpAglTj9sIBWINXa8x5w==" crossorigin="anonymous" referrerpolicy="no-referrer" />`,
				`<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/skeleton/2.0.4/skeleton.min.css" integrity="sha512-EZLkOqwILORob+p0BXZc+Vm3RgJBOe1Iq/0fiI7r/wJgzOFZMlsqTa29UEl6v6U6gsV4uIpsNZoV32YZqrCRCQ==" crossorigin="anonymous" referrerpolicy="no-referrer" />`,

				`<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/codemirror@6.65.7/lib/codemirror.css" integrity="sha256-60lOqXLSZh74b39qxlbdZ4bXIeScnBtG4euWfktvm/M=" crossorigin="anonymous" referrerpolicy="no-referrer">`,
				`<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/codemirror@6.65.7/theme/twilight.css" integrity="sha256-ltApMINjtnG8JfMP1DnRcEzkUFuMNil+PKQRoARd8js=" crossorigin="anonymous" referrerpolicy="no-referrer">`,

				`<script src="https://cdn.jsdelivr.net/npm/codemirror@6.65.7/lib/codemirror.js" integrity="sha256-pPHDA2XV+FlRoiFgsHp840BIy2QMKDVFyBkhMrRu+6g=" crossorigin="anonymous" referrerpolicy="no-referrer"></script>`,
				`<script src="https://cdn.jsdelivr.net/npm/codemirror@6.65.7/keymap/sublime.js" integrity="sha256-HVMqZ1lIRctxY7CQoG2pKHuJaUSrXyLfZ/FOF2vUNOU=" crossorigin="anonymous" referrerpolicy="no-referrer"></script>`,
				`<script src="https://cdn.jsdelivr.net/npm/codemirror@6.65.7/mode/go/go.js" integrity="sha256-TjsIvaL4XJ72jAPyTAxco7i6R5EntsvFDzPLlZPbrwA=" crossorigin="anonymous" referrerpolicy="no-referrer"></script>`,
				`<script src="https://cdn.jsdelivr.net/npm/codemirror@6.65.7/addon/selection/active-line.js" integrity="sha256-OvvPeINcm9w0LjmSxT2bdChnImE7saitydFA7chzfug=" crossorigin="anonymous" referrerpolicy="no-referrer"></script>`,
				`<script src="https://cdn.jsdelivr.net/npm/codemirror@6.65.7/addon/edit/matchbrackets.js" integrity="sha256-nQ5e5MGZ/L5IepAYYMLPfQByY4Umd91YWAmPyqhYDx4=" crossorigin="anonymous" referrerpolicy="no-referrer"></script>`,

				`<script src="https://cdn.jsdelivr.net/npm/@textea/json-viewer"></script>`,
			},
			Styles: []string{},
			Scripts: []string{
				"/web/scripts/main.js",
			},
		}

		// Create static build if specified
		if *build {
			// Deploy under a path
			if *path != "" {
				h.Resources = app.GitHubPages(*path)
			}

			if err := app.GenerateStaticWebsite(*out, h); err != nil {
				log.Fatalf("could not build: %v\n", err)
			}
			//for _, s := range append(h.Scripts, h.Styles...) {
			//	filename := filepath.Join(*out, s)
			//	if err := createStaticDir(filepath.Join(dir, path), ""); err != nil {
			//		return errors.New("creating web directory failed").Wrap(err)
			//	}
			//	os.Open()
			//}
		}

		// Serve if specified
		if *serve {
			log.Printf("Gastly listening on http://%v\n", *laddr)

			if err := http.ListenAndServe(*laddr, h); err != nil {
				log.Fatalf("could not open Gionic: %v\n", err)
			}
		}
	}
}
