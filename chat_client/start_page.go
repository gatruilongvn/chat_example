package main

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type StartPage struct {
	nameInputTb widget.Editor
	errorText   string
}

func (st *StartPage) HandleAction(gtx layout.Context) {

}

func newStartPage() *StartPage {
	st := StartPage{}
	st.errorText = " "
	return &st
}

// Start page layout
func (st *StartPage) PageLayout(gtx layout.Context) layout.Dimensions {

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
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						title := material.Label(theme, unit.Sp(float32(14)), "Enter your name")
						title.Color = color.NRGBA{R: 127, G: 0, B: 0, A: 255}
						title.Alignment = text.Middle
						title.Layout(gtx)
						return title.Layout(gtx)
					},
				), layout.Rigid(
					func(gtx C) D {

						// Define insets ...
						margins := layout.Inset{
							Top:    unit.Dp(10),
							Right:  unit.Dp(25),
							Bottom: unit.Dp(0),
							Left:   unit.Dp(25),
						}
						paddings := layout.Inset{
							Top:    unit.Dp(5),
							Right:  unit.Dp(5),
							Bottom: unit.Dp(5),
							Left:   unit.Dp(5),
						}
						// ... and borders ...
						border := widget.Border{
							Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
							CornerRadius: unit.Dp(3),
							Width:        unit.Dp(2),
						}

						// ... and material design ...
						ed := material.Editor(theme, &st.nameInputTb, "")
						st.nameInputTb.Alignment = text.Middle
						st.nameInputTb.SingleLine = true
						// ... before laying it out, one inside the other
						return margins.Layout(gtx,
							func(gtx C) D {
								return border.Layout(gtx, func(gtx C) D {
									return paddings.Layout(gtx, func(gtx C) D {
										return ed.Layout(gtx)
									})
								})
							},
						)
					},
				),
				layout.Rigid(
					func(gtx C) D {

						// Define insets ...
						margins := layout.Inset{
							Top:    unit.Dp(5),
							Right:  unit.Dp(25),
							Bottom: unit.Dp(0),
							Left:   unit.Dp(25),
						}
						// ... and material design ...

						st.nameInputTb.Alignment = text.Middle
						// ... before laying it out, one inside the other
						return margins.Layout(gtx,
							func(gtx C) D {
								title := material.Label(theme, unit.Sp(float32(14)), st.errorText)
								title.Color = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
								title.Alignment = text.Middle
								title.Layout(gtx)
								return title.Layout(gtx)
							},
						)
					},
				),
			)
		}),
	)
}
