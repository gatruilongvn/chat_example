package main

import (
	"fmt"
	"image/color"
	"net"
	"runtime"
	"strings"

	"gioui.org/font/gofont"
	"gioui.org/io/profile"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type UI struct {
	startPage      *StartPage
	chatPage       *ChatPage
	currentPage    Page
	connectBtn     widget.Clickable
	profiling      bool
	profile        profile.Event
	lastMallocs    uint64
	connectBtnText string
}

var (
	theme          *material.Theme
	userName       string
	connection     net.Conn
	ui             *UI
	updateMessages chan string
)

func newUI() *UI {
	u := &UI{}
	u.connectBtnText = "Connect"
	u.startPage = newStartPage()
	u.chatPage = newChatPage()
	ui = u
	return u
}

func rgb(c uint32) color.NRGBA {
	return argb((0xff << 24) | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func (u *UI) layoutTimings(gtx layout.Context) {
	if !u.profiling {
		return
	}
	for _, e := range gtx.Events(u) {
		if e, ok := e.(profile.Event); ok {
			u.profile = e
		}
	}
	profile.Op{Tag: u}.Add(gtx.Ops)
	var mstats runtime.MemStats
	runtime.ReadMemStats(&mstats)
	mallocs := mstats.Mallocs - u.lastMallocs
	u.lastMallocs = mstats.Mallocs
	layout.NE.Layout(gtx, func(gtx C) D {
		return layout.Inset{Top: unit.Dp(16)}.Layout(gtx, func(gtx C) D {
			txt := fmt.Sprintf("m: %d %s", mallocs, u.profile.Timings)
			lbl := material.Caption(theme, txt)
			lbl.Font.Variant = "Mono"
			return lbl.Layout(gtx)
		})
	})
}
func init() {
	theme = material.NewTheme(gofont.Collection())
	theme.Palette.Fg = rgb(0x333333)
}
func (u *UI) Layout(gtx layout.Context) {
	if ui.connectBtn.Clicked() {
		switch u.connectBtnText {
		case "Connect":
			u.connect()
		case "Dissconnect":
			u.dissconnect()
		}
	}
	if u.currentPage == nil {
		u.currentPage = u.startPage
	}

	u.currentPage.HandleAction(gtx)
	layout.Flex{
		//Vertical alignment, from top to bottom
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return u.connectBtnLayout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return u.currentPage.PageLayout(gtx)
		}),
	)

	u.layoutTimings(gtx)
}

func (u *UI) connect() {
	go func() {
		userName = u.startPage.nameInputTb.Text()
		userName = strings.TrimSpace(userName)

		if len(userName) == 0 {
			u.startPage.errorText = "Input your name, dude!"
		} else {
			conn, err := net.Dial("tcp", "localhost:3000")
			connection = conn
			if err != nil {
				u.startPage.errorText = "Connect to server error"
				fmt.Println(err)
			} else {
				u.connectBtnText = "Dissconnect"
				u.startPage.errorText = " "
				u.chatPage.connection = conn
				u.currentPage = u.chatPage
				go u.chatPage.OnMessage()
			}
		}
	}()

}

func (u *UI) dissconnect() {
	err := ui.chatPage.connection.Close()
	if err != nil {
		fmt.Println("error occur! Can not close connection")
	} else {
		u.connectBtnText = "Connect"
		u.currentPage = ui.startPage
	}
}

func (ui *UI) connectBtnLayout(gtx C) D {
	return layout.Flex{
		// Vertical alignment, from top to bottom
		Axis:    layout.Vertical,
		Spacing: layout.SpaceStart,
	}.Layout(gtx,
		layout.Rigid(
			func(gtx layout.Context) layout.Dimensions {
				// 		// The button

				// We start by defining a set of margins
				margins := layout.Inset{
					Top:    unit.Dp(25),
					Bottom: unit.Dp(25),
					Right:  unit.Dp(120),
					Left:   unit.Dp(120),
				}
				// Then we lay out within those margins
				return margins.Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						// The text on the button depends on program state
						var text = ui.connectBtnText
						btn := material.Button(theme, &ui.connectBtn, text)
						return btn.Layout(gtx)
					},
				)
			},
		),
	)
}
