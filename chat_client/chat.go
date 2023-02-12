package main

import (
	"bufio"
	"fmt"
	"image/color"
	"net"
	"strings"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type ChatPage struct {
	messageEditor widget.Editor
	fab           *widget.Clickable
	fabIcon       *widget.Icon
	messagesList  *widget.List
	messages      []string
	connection    net.Conn
}

func (cp *ChatPage) HandleAction(gtx layout.Context) {
	if cp.fab.Clicked() {
		cp.SendMessage()
	}
	for _, ev := range cp.messageEditor.Events() {
		switch ev := ev.(type) {
		case widget.SubmitEvent:
			if len(ev.Text) != 0 {
				cp.SendMessage()
			}
		}

	}
}

func (cp *ChatPage) OnMessage() {
	for {
		reader := bufio.NewReader(connection)
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		cp.messages = append(cp.messages, msg)
		updateMessages <- msg
	}
	connection.Close()
}

func (cp *ChatPage) SendMessage() {
	msg := cp.messageEditor.Text()
	cp.messageEditor.SetText("")
	msg = userName + ": " + msg + "\n"
	fmt.Println(msg)
	_, err := connection.Write([]byte(msg))
	if err != nil {
		fmt.Println(err)
	}
	cp.messages = append(cp.messages, msg)
}

func newChatPage() *ChatPage {
	cp := ChatPage{}
	cp.messages = make([]string, 100)
	cp.fab = new(widget.Clickable)
	cp.fabIcon, _ = widget.NewIcon(icons.ContentSend)
	cp.messagesList = &widget.List{
		List: layout.List{
			Axis:        layout.Vertical,
			ScrollToEnd: true,
		},
	}
	return &cp
}

func (cp *ChatPage) PageLayout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		//Vertical alignment, from top to bottom
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				// Vertical alignment, from top to bottom
				Axis:    layout.Vertical,
				Spacing: layout.SpaceSides,
			}.Layout(gtx,
				layout.Flexed(1,
					func(gtx layout.Context) layout.Dimensions {
						margins := layout.Inset{
							Top:    unit.Dp(0),
							Right:  unit.Dp(5),
							Bottom: unit.Dp(5),
							Left:   unit.Dp(5),
						}
						return margins.Layout(gtx, func(gtx C) D {
							border := widget.Border{
								Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
								CornerRadius: unit.Dp(3),
								Width:        unit.Dp(2),
							}
							return border.Layout(gtx, func(gtx C) D {
								paddings := layout.Inset{
									Top:    unit.Dp(5),
									Right:  unit.Dp(5),
									Bottom: unit.Dp(5),
									Left:   unit.Dp(5),
								}
								return paddings.Layout(gtx, func(gtx C) D {
									return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Flexed(1, func(gtx C) D {

											gtx.Constraints.Min.X = gtx.Constraints.Max.X

											l := cp.messagesList

											return material.List(theme, l).Layout(gtx, len(cp.messages), func(gtx C, i int) D {
												paragraph := material.Label(theme, unit.Sp(float32(16)), strings.Trim(cp.messages[i], "\n"))
												// The text is centered
												paragraph.Alignment = text.Start
												// Return the laid out paragraph
												return paragraph.Layout(gtx)
											})
										}),
									)

								})

							})
						})
					},
				),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					margins := layout.Inset{
						Top:    unit.Dp(0),
						Right:  unit.Dp(5),
						Bottom: unit.Dp(5),
						Left:   unit.Dp(5),
					}
					return margins.Layout(gtx, func(gtx C) D {
						border := widget.Border{
							Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
							CornerRadius: unit.Dp(3),
							Width:        unit.Dp(2),
						}
						return border.Layout(gtx, func(gtx C) D {
							padding := layout.Inset{
								Top:    unit.Dp(5),
								Right:  unit.Dp(5),
								Bottom: unit.Dp(4),
								Left:   unit.Dp(5),
							}
							return padding.Layout(gtx, func(gtx C) D {
								cp.messageEditor.Alignment = text.Start
								cp.messageEditor.SingleLine = false
								cp.messageEditor.Submit = true
								// ... and material design ...
								ed := material.Editor(theme, &cp.messageEditor, "Write a message")

								// ... before laying it out, one inside the other
								return margins.Layout(gtx,
									func(gtx C) D {
										return ed.Layout(gtx)
									},
								)
							})
						})
					})
				}),
				layout.Rigid(func(gtx C) D {
					margins := layout.Inset{
						Top:    unit.Dp(5),
						Right:  unit.Dp(5),
						Bottom: unit.Dp(5),
						Left:   unit.Dp(0),
					}
					return margins.Layout(gtx, func(gtx C) D {
						bt := material.IconButton(theme, cp.fab, cp.fabIcon, "Sent Button")
						bt.Size = 16
						bt.Inset = layout.UniformInset(5)
						return bt.Layout(gtx)
					})
				}),
			)
		}),
	)
}
