package microui

/*============================================================================
** layout
**============================================================================*/

func (ctx *Context) PushLayout(body Rect, scroll Vec2) {
	layout := Layout{}
	layout.Body = NewRect(body.X-scroll.X, body.Y-scroll.Y, body.W, body.H)
	layout.Max = NewVec2(-0x1000000, -0x1000000)

	// push()
	ctx.LayoutStack = append(ctx.LayoutStack, layout)

	ctx.LayoutRow(1, []int{0}, 0)
}

func (ctx *Context) LayoutBeginColumn() {
	ctx.PushLayout(ctx.LayoutNext(), NewVec2(0, 0))
}

func (ctx *Context) LayoutEndColumn() {
	b := ctx.GetLayout()
	// pop()
	expect(len(ctx.LayoutStack) > 0)
	ctx.LayoutStack = ctx.LayoutStack[:len(ctx.LayoutStack)-1]
	// inherit position/next_row/max from child layout if they are greater
	a := ctx.GetLayout()
	a.Position.X = mu_max(a.Position.X, b.Position.X+b.Body.X-a.Body.X)
	a.NextRow = mu_max(a.NextRow, b.NextRow+b.Body.Y-a.Body.Y)
	a.Max.X = mu_max(a.Max.X, b.Max.X)
	a.Max.Y = mu_max(a.Max.Y, b.Max.Y)
}

func (ctx *Context) LayoutRow(items int, widths []int, height int) {
	layout := ctx.GetLayout()

	expect(len(widths) <= MU_MAX_WIDTHS)
	copy(layout.Widths[:], widths)

	layout.Items = items
	layout.Position = NewVec2(layout.Indent, layout.NextRow)
	layout.Size.Y = height
	layout.ItemIndex = 0
}

// sets layout size.x
func (ctx *Context) LayoutWidth(width int) {
	ctx.GetLayout().Size.X = width
}

// sets layout size.y
func (ctx *Context) LayoutHeight(height int) {
	ctx.GetLayout().Size.Y = height
}

func (ctx *Context) LayoutSetNext(r Rect, relative bool) {
	layout := ctx.GetLayout()
	layout.Next = r
	if relative {
		layout.NextType = RELATIVE
	} else {
		layout.NextType = ABSOLUTE
	}
}

func (ctx *Context) LayoutNext() Rect {
	layout := ctx.GetLayout()
	style := ctx.Style
	var res Rect

	if layout.NextType != 0 {
		// handle rect set by `mu_layout_set_next`
		next_type := layout.NextType
		layout.NextType = 0
		res = layout.Next

		if next_type == ABSOLUTE {
			ctx.LastRect = res
			return ctx.LastRect
		}
	} else {
		// handle next row
		if layout.ItemIndex == layout.Items {
			ctx.LayoutRow(layout.Items, nil, layout.Size.Y)
		}

		// position
		res.X = layout.Position.X
		res.Y = layout.Position.Y

		// size
		if layout.Items > 0 {
			res.W = layout.Widths[layout.ItemIndex]
		} else {
			res.W = layout.Size.X
		}
		res.H = layout.Size.Y
		if res.W == 0 {
			res.W = style.Size.X + style.Padding*2
		}
		if res.H == 0 {
			res.H = style.Size.Y + style.Padding*2
		}
		if res.W < 0 {
			res.W += layout.Body.W - res.X + 1
		}
		if res.H < 0 {
			res.H += layout.Body.H - res.Y + 1
		}

		layout.ItemIndex++
	}

	// update position
	layout.Position.X += res.W + style.Spacing
	layout.NextRow = mu_max(layout.NextRow, res.Y+res.H+style.Spacing)

	// apply body offset
	res.X += layout.Body.X
	res.Y += layout.Body.Y

	// update max position
	layout.Max.X = mu_max(layout.Max.X, res.X+res.W)
	layout.Max.Y = mu_max(layout.Max.Y, res.Y+res.H)

	ctx.LastRect = res
	return ctx.LastRect
}
