// 4 december 2014

// TODO find a better place for this
static LONG rowHeight(struct table *t, HDC dc, BOOL select)
{
	BOOL release;
	HFONT prevfont, newfont;
	TEXTMETRICW tm;

	release = FALSE;
	if (dc == NULL) {
		dc = GetDC(t->hwnd);
		if (dc == NULL)
			panic("error getting Table DC for rowHeight()");
		release = TRUE;
	}
	if (select)
		prevfont = selectFont(t, dc, &newfont);
	if (GetTextMetricsW(dc, &tm) == 0)
		panic("error getting text metrics for rowHeight()");
	if (select)
		deselectFont(dc, prevfont, newfont);
	if (release)
		if (ReleaseDC(t->hwnd, dc) == 0)
			panic("error releasing Table DC for rowHeight()");
	return tm.tmHeight;
}

#define rowht(t) rowHeight(t, NULL, TRUE)

struct rowcol {
	intptr_t row;
	intptr_t column;
};

static struct rowcol clientCoordToRowColumn(struct table *t, POINT pt)
{
	RECT r;
	struct rowcol rc;
	intptr_t i;

	// initial values for the PtInRect() check
	rc.row = -1;
	rc.column = -1;

	if (GetClientRect(t->hwnd, &r) == 0)
		panic("error getting Table client rect in clientCoordToRowColumn()");
	r.top += t->headerHeight;
	if (PtInRect(&r, pt) == 0)
		return rc;

	// the row is easy
	pt.y -= t->headerHeight;
	rc.row = (pt.y / rowht(t)) + t->vscrollpos;

	// the column... not so much
	// we scroll p.x, then subtract column widths until we cross the left edge of the control
	pt.x += t->hscrollpos;
	rc.column = 0;
	for (i = 0; i < t->nColumns; i++) {
		pt.x -= columnWidth(t, i);
		// use <, not <=, here:
		// assume r.left and t->hscrollpos == 0;
		// given the first column is 100 wide,
		// pt.x == 0 (first pixel of col 0) -> p.x - 100 == -100 < 0 -> break
		// pt.x == 99 (last pixel of col 0) -> p.x - 100 == -1 < 0 -> break
		// pt.x == 100 (first pixel of col 1) -> p.x - 100 == 0 >= 0 -> next column
		if (pt.x < r.left)
			break;
		rc.column++;
	}
	// TODO what happens if the break was never taken?

	return rc;
}

// same as client coordinates, but stored in a lParam (like the various mouse messages provide)
static struct rowcol lParamToRowColumn(struct table *t, LPARAM lParam)
{
	POINT pt;

	pt.x = GET_X_LPARAM(lParam);
	pt.y = GET_Y_LPARAM(lParam);
	return clientCoordToRowColumn(t, pt);
}

// returns TRUE if the row is visible (even partially visible) and thus has a rectangle in the client area; FALSE otherwise
static BOOL rowColumnToClientRect(struct table *t, struct rowcol rc, RECT *r)
{
	RECT client;
	RECT out;			// don't change r if we return FALSE
	LONG height;
	intptr_t xpos;
	intptr_t i;

	if (rc.row < t->vscrollpos)
		return FALSE;
	rc.row -= t->vscrollpos;		// align with client.top

	if (GetClientRect(t->hwnd, &client) == 0)
		panic("error getting Table client rect in rowColumnToClientRect()");
	client.top += t->headerHeight;

	height = rowht(t);
	out.top = client.top + (rc.row * height);
	if (out.top >= client.bottom)		// >= because RECT.bottom is the first pixel outside the rectangle
		return FALSE;
	out.bottom = out.top + height;

	// and again the columns are the hard part
	// so we start with client.left - t->hscrollpos, then keep adding widths until we get to the column we want
	xpos = client.left - t->hscrollpos;
	for (i = 0; i < rc.column; i++)
		xpos += columnWidth(t, i);
	// did we stray too far to the right? if so it's not visible
	if (xpos >= client.right)		// >= because RECT.right is the first pixel outside the rectangle
		return FALSE;
	out.left = xpos;
	out.right = xpos + columnWidth(t, rc.column);
	// and is this too far to the left?
	if (out.right < client.left)		// < because RECT.left is the first pixel inside the rect
		return FALSE;

	*r = out;
	return TRUE;
}

// TODO idealCoordToRowColumn/rowColumnToIdealCoord?
