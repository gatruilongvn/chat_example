package main

import (
	"context"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
)

type App struct {
	w *app.Window

	ui        *UI
	ctx       context.Context
	ctxCancel context.CancelFunc
}

// Page interface
type Page interface {
	PageLayout(gtx layout.Context) layout.Dimensions
	HandleAction(gtx layout.Context)
}

type (
	C = layout.Context
	D = layout.Dimensions
)

func main() {
	updateMessages = make(chan string)
	go func() {

		w := app.NewWindow(
			app.Size(unit.Dp(400), unit.Dp(600)),
			app.Title("Chat example"),
		)
		if err := newApp(w).run(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func (a *App) run() error {

	var ops op.Ops

	for {
		select {

		case <-updateMessages:
			a.w.Invalidate()
		case e := <-a.w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.StageEvent:
				if e.Stage >= system.StageRunning {
					if a.ctxCancel == nil {
						a.ctx, a.ctxCancel = context.WithCancel(context.Background())
					}
				} else {
					if a.ctxCancel != nil {
						a.ctxCancel()
						a.ctxCancel = nil
					}
				}
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				// register a global key listener for the escape key wrapping our entire UI.
				area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
				key.InputOp{
					Tag:  a.w,
					Keys: key.NameEscape + `|Short-P|` + key.NameBack,
				}.Add(gtx.Ops)

				// check for presses of global keyboard shortcuts and process them.
				for _, event := range gtx.Events(a.w) {
					switch event := event.(type) {
					case key.Event:
						switch event.Name {
						case key.NameEscape:
							return nil
						case "P":
							if event.Modifiers.Contain(key.ModShortcut) && event.State == key.Press {
								a.ui.profiling = !a.ui.profiling
								a.w.Invalidate()
							}
						}
					}
				}
				a.ui.Layout(gtx)
				area.Pop()
				e.Frame(gtx.Ops)
			}
		}
	}
}

func newApp(w *app.Window) *App {

	a := &App{
		w: w,
	}
	a.ui = newUI()
	return a
}
