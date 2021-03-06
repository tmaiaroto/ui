// 8 december 2014

struct drawCellParams {
	intptr_t row;
	intptr_t column;
	LONG x;
	LONG y;
	LONG width;		// of column
	LONG height;		// rowHeight()
	LRESULT xoff;		// result of HDM_GETBITMAPMARGIN
};

static void drawCell(struct table *t, HDC dc, struct drawCellParams *p)
{
	RECT r;
	WCHAR msg[200];
	int n;
	HBRUSH background;
	int textColor;

	// TODO verify these two
	background = (HBRUSH) (COLOR_WINDOW + 1);
	textColor = COLOR_WINDOWTEXT;
	// TODO get rid of the selectedColumn bits
	if (t->selectedRow == p->row && t->selectedColumn == p->column) {
		// these are the colors wine uses (http://source.winehq.org/source/dlls/comctl32/listview.c)
		// the two for unfocused are also suggested by http://stackoverflow.com/questions/10428710/windows-forms-inactive-highlight-color
		background = (HBRUSH) (COLOR_HIGHLIGHT + 1);
		textColor = COLOR_HIGHLIGHTTEXT;
		if (GetFocus() != t->hwnd) {
			background = (HBRUSH) (COLOR_BTNFACE + 1);
			textColor = COLOR_BTNTEXT;
		}
		// TODO disabled
	}

	r.left = p->x;
	r.right = p->x + p->width;
	r.top = p->y;
	r.bottom = p->y + p->height;
	if (FillRect(dc, &r, background) == 0)
		panic("error filling Table cell background");

	// now offset the content to where inside the cell it should be
	r.left += p->xoff;
	// TODO vertical center content too

	if (SetTextColor(dc, GetSysColor(textColor)) == CLR_INVALID)
		panic("error setting Table cell text color");
	if (SetBkMode(dc, TRANSPARENT) == 0)
		panic("error setting transparent text drawing mode for Table cell");
	n = wsprintf(msg, L"(%d,%d)", p->row, p->column);
	if (DrawTextExW(dc, msg, n, &r, DT_END_ELLIPSIS | DT_LEFT | DT_NOPREFIX | DT_SINGLELINE, NULL) == 0)
		panic("error drawing Table cell text");
}

static void draw(struct table *t, HDC dc, RECT cliprect, RECT client)
{
	intptr_t i, j;
	int x = 0;
	HFONT prevfont, newfont;
	struct drawCellParams p;

	prevfont = selectFont(t, dc, &newfont);

	client.top += t->headerHeight;

	ZeroMemory(&p, sizeof (struct drawCellParams));
	p.height = rowHeight(t, dc, FALSE);
	p.xoff = SendMessageW(t->header, HDM_GETBITMAPMARGIN, 0, 0);

	p.y = client.top;
	for (i = t->vscrollpos; i < t->count; i++) {
		p.row = i;
		p.x = client.left - t->hscrollpos;
		for (j = 0; j < t->nColumns; j++) {
			p.column = j;
			p.width = columnWidth(t, p.column);
			drawCell(t, dc, &p);
			p.x += p.width;
		}
		p.y += p.height;
	}

	deselectFont(dc, prevfont, newfont);
}

HANDLER(drawHandlers)
{
	HDC dc;
	PAINTSTRUCT ps;
	RECT client;
	RECT r;

	if (uMsg != WM_PAINT && uMsg != WM_PRINTCLIENT)
		return FALSE;
	if (GetClientRect(t->hwnd, &client) == 0)
		panic("error getting client rect for Table painting");
	if (uMsg == WM_PAINT) {
		dc = BeginPaint(t->hwnd, &ps);
		if (dc == NULL)
			panic("error beginning Table painting");
		r = ps.rcPaint;
	} else {
		dc = (HDC) wParam;
		r = client;
	}
	draw(t, dc, r, client);
	if (uMsg == WM_PAINT)
		EndPaint(t->hwnd, &ps);
	// this is correct for WM_PRINTCLIENT; see http://stackoverflow.com/a/27362258/3408572
	*lResult = 0;
	return TRUE;
}
