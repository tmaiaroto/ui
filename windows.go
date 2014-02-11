// 8 february 2014
package main

import (
	"syscall"
	"unsafe"
)

// Window styles.
const (
	WS_BORDER = 0x00800000
	WS_CAPTION = 0x00C00000
	WS_CHILD = 0x40000000
	WS_CHILDWINDOW = 0x40000000
	WS_CLIPCHILDREN = 0x02000000
	WS_CLIPSIBLINGS = 0x04000000
	WS_DISABLED = 0x08000000
	WS_DLGFRAME = 0x00400000
	WS_GROUP = 0x00020000
	WS_HSCROLL = 0x00100000
	WS_ICONIC = 0x20000000
	WS_MAXIMIZE = 0x01000000
	WS_MAXIMIZEBOX = 0x00010000
	WS_MINIMIZE = 0x20000000
	WS_MINIMIZEBOX = 0x00020000
	WS_OVERLAPPED = 0x00000000
	WS_OVERLAPPEDWINDOW = (WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX)
	WS_POPUP = 0x80000000
	WS_POPUPWINDOW = (WS_POPUP | WS_BORDER | WS_SYSMENU)
	WS_SIZEBOX = 0x00040000
	WS_SYSMENU = 0x00080000
	WS_TABSTOP = 0x00010000
	WS_THICKFRAME = 0x00040000
	WS_TILED = 0x00000000
	WS_TILEDWINDOW = (WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX)
	WS_VISIBLE = 0x10000000
	WS_VSCROLL = 0x00200000
)

// Extended window styles.
const (
	WS_EX_ACCEPTFILES = 0x00000010
	WS_EX_APPWINDOW = 0x00040000
	WS_EX_CLIENTEDGE = 0x00000200
//	WS_EX_COMPOSITED = 0x02000000	// [Windows 2000:This style is not supported.]
	WS_EX_CONTEXTHELP = 0x00000400
	WS_EX_CONTROLPARENT = 0x00010000
	WS_EX_DLGMODALFRAME = 0x00000001
	WS_EX_LAYERED = 0x00080000
	WS_EX_LAYOUTRTL = 0x00400000
	WS_EX_LEFT = 0x00000000
	WS_EX_LEFTSCROLLBAR = 0x00004000
	WS_EX_LTRREADING = 0x00000000
	WS_EX_MDICHILD = 0x00000040
	WS_EX_NOACTIVATE = 0x08000000
	WS_EX_NOINHERITLAYOUT = 0x00100000
	WS_EX_NOPARENTNOTIFY = 0x00000004
	WS_EX_OVERLAPPEDWINDOW = (WS_EX_WINDOWEDGE | WS_EX_CLIENTEDGE)
	WS_EX_PALETTEWINDOW = (WS_EX_WINDOWEDGE | WS_EX_TOOLWINDOW | WS_EX_TOPMOST)
	WS_EX_RIGHT = 0x00001000
	WS_EX_RIGHTSCROLLBAR = 0x00000000
	WS_EX_RTLREADING = 0x00002000
	WS_EX_STATICEDGE = 0x00020000
	WS_EX_TOOLWINDOW = 0x00000080
	WS_EX_TOPMOST = 0x00000008
	WS_EX_TRANSPARENT = 0x00000020
	WS_EX_WINDOWEDGE = 0x00000100
)

// bizarrely, this value is given on the page for CreateMDIWindow, but not CreateWindow or CreateWindowEx
// I do it this way because Go won't let me shove the exact value into an int
var (
	_uCW_USEDEFAULT uint = 0x80000000
	CW_USEDEFAULT = int(_uCW_USEDEFAULT)
)

// GetSysColor values. These can be cast to HBRUSH (after adding 1) for WNDCLASS as well.
const (
	COLOR_3DDKSHADOW = 21
	COLOR_3DFACE = 15
	COLOR_3DHIGHLIGHT = 20
	COLOR_3DHILIGHT = 20
	COLOR_3DLIGHT = 22
	COLOR_3DSHADOW = 16
	COLOR_ACTIVEBORDER = 10
	COLOR_ACTIVECAPTION = 2
	COLOR_APPWORKSPACE = 12
	COLOR_BACKGROUND = 1
	COLOR_BTNFACE = 15
	COLOR_BTNHIGHLIGHT = 20
	COLOR_BTNHILIGHT = 20
	COLOR_BTNSHADOW = 16
	COLOR_BTNTEXT = 18
	COLOR_CAPTIONTEXT = 9
	COLOR_DESKTOP = 1
	COLOR_GRADIENTACTIVECAPTION = 27
	COLOR_GRADIENTINACTIVECAPTION = 28
	COLOR_GRAYTEXT = 17
	COLOR_HIGHLIGHT = 13
	COLOR_HIGHLIGHTTEXT = 14
	COLOR_HOTLIGHT = 26
	COLOR_INACTIVEBORDER = 11
	COLOR_INACTIVECAPTION = 3
	COLOR_INACTIVECAPTIONTEXT = 19
	COLOR_INFOBK = 24
	COLOR_INFOTEXT = 23
	COLOR_MENU = 4
//	COLOR_MENUHILIGHT = 29	// [Windows 2000:This value is not supported.]
//	COLOR_MENUBAR = 30		// [Windows 2000:This value is not supported.]
	COLOR_MENUTEXT = 7
	COLOR_SCROLLBAR = 0
	COLOR_WINDOW = 5
	COLOR_WINDOWFRAME = 6
	COLOR_WINDOWTEXT = 8
)

// ShowWindow settings.
const (
	SW_FORCEMINIMIZE = 11
	SW_HIDE = 0
	SW_MAXIMIZE = 3
	SW_MINIMIZE = 6
	SW_RESTORE = 9
	SW_SHOW = 5
	SW_SHOWDEFAULT = 10
	SW_SHOWMAXIMIZED = 3
	SW_SHOWMINIMIZED = 2
	SW_SHOWMINNOACTIVE = 7
	SW_SHOWNA = 8
	SW_SHOWNOACTIVATE = 4
	SW_SHOWNORMAL = 1
)

var (
	createWindowEx = user32.NewProc("CreateWindowExW")
	destroyWindow = user32.NewProc("DestroyWindow")
	enumChildWindows = user32.NewProc("EnumChildWindows")
	showWindow = user32.NewProc("ShowWindow")
)

// TODO use lpParam
func CreateWindowEx(dwExStyle uint32, lpClassName string, lpWindowName string, dwStyle uint32, x int, y int, nWidth int, nHeight int, hwndParent HWND, hMenu HMENU, hInstance HANDLE, lpParam interface{}) (hwnd HWND, err error) {
	r1, _, err := createWindowEx.Call(
		uintptr(dwExStyle),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpClassName))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpWindowName))),
		uintptr(dwStyle),
		uintptr(x),
		uintptr(y),
		uintptr(nWidth),
		uintptr(nHeight),
		uintptr(hwndParent),
		uintptr(hMenu),
		uintptr(hInstance),
		uintptr(0))
	if r1 == 0 {		// failure
		return NULL, err
	}
	return HWND(r1), nil
}

func DestroyWindow(hWnd HWND) (err error) {
	r1, _, err := destroyWindow.Call(uintptr(hWnd))
	if r1 == 0 {		// failure
		return err
	}
	return nil
}

type WNDENUMPROC func(hwnd HWND, lParam LPARAM) (cont bool)
type _WNDENUMPROC func(hwnd HWND, lParam LPARAM) int

func enumChildProc(p WNDENUMPROC) _WNDENUMPROC {
	return func(hwnd HWND, lParam LPARAM) int {
		if p(hwnd, lParam) {
			return TRUE
		}
		return FALSE
	}
}

// TODO figure out how to handle errors
func EnumChildWindows(hWndParent HWND, lpEnumFunc WNDENUMPROC, lParam LPARAM) (err error) {
	enumChildWindows.Call(
		uintptr(hWndParent),
		syscall.NewCallback(enumChildProc(lpEnumFunc)),
		uintptr(lParam))
	return nil
}

// TODO figure out how to handle errors
func ShowWindow(hWnd HWND, nCmdShow int) (previouslyVisible bool, err error) {
	r1, _, _ := showWindow.Call(
		uintptr(hWnd),
		uintptr(nCmdShow))
	return r1 != 0, nil
}

// WM_SETICON and WM_GETICON values.
const (
	ICON_BIG = 1
	ICON_SMALL = 0
	ICON_SMALL2 = 2		// WM_GETICON only?
)

// Window messages.
const (
	MN_GETHMENU = 0x01E1
	WM_ERASEBKGND = 0x0014
	WM_GETFONT = 0x0031
	WM_GETTEXT = 0x000D
	WM_GETTEXTLENGTH = 0x000E
	WM_SETFONT = 0x0030
	WM_SETICON = 0x0080
	WM_SETTEXT = 0x000C
)

// WM_INPUTLANGCHANGEREQUEST values.
const (
	INPUTLANGCHANGE_BACKWARD = 0x0004
	INPUTLANGCHANGE_FORWARD = 0x0002
	INPUTLANGCHANGE_SYSCHARSET = 0x0001
)

// WM_NCCALCSIZE return values.
const (
	WVR_ALIGNTOP = 0x0010
	WVR_ALIGNRIGHT = 0x0080
	WVR_ALIGNLEFT = 0x0020
	WVR_ALIGNBOTTOM = 0x0040
	WVR_HREDRAW = 0x0100
	WVR_VREDRAW = 0x0200
	WVR_REDRAW = 0x0300
	WVR_VALIDRECTS = 0x0400
)

// WM_SHOWWINDOW reasons (lParam).
const (
	SW_OTHERUNZOOM = 4
	SW_OTHERZOOM = 2
	SW_PARENTCLOSING = 1
	SW_PARENTOPENING = 3
)

// WM_SIZE values.
const (
	SIZE_MAXHIDE = 4
	SIZE_MAXIMIZED = 2
	SIZE_MAXSHOW = 3
	SIZE_MINIMIZED = 1
	SIZE_RESTORED = 0
)

// WM_SIZING edge values (wParam).
const (
	WMSZ_BOTTOM = 6
	WMSZ_BOTTOMLEFT = 7
	WMSZ_BOTTOMRIGHT = 8
	WMSZ_LEFT = 1
	WMSZ_RIGHT = 2
	WMSZ_TOP = 3
	WMSZ_TOPLEFT = 4
	WMSZ_TOPRIGHT = 5
)

// WM_STYLECHANGED and WM_STYLECHANGING values (wParam).
const (
	GWL_EXSTYLE = -20
	GWL_STYLE = -16
)

// Window notifications.
const (
	WM_ACTIVATEAPP = 0x001C
	WM_CANCELMODE = 0x001F
	WM_CHILDACTIVATE = 0x0022
	WM_CLOSE = 0x0010
	WM_COMPACTING = 0x0041
	WM_CREATE = 0x0001
	WM_DESTROY = 0x0002
//	WM_DPICHANGED = 0x02E0		// Windows 8.1 and newer only
	WM_ENABLE = 0x000A
	WM_ENTERSIZEMOVE = 0x0231
	WM_EXITSIZEMOVE = 0x0232
	WM_GETICON = 0x007F
	WM_GETMINMAXINFO = 0x0024
	WM_INPUTLANGCHANGE = 0x0051
	WM_INPUTLANGCHANGEREQUEST = 0x0050
	WM_MOVE = 0x0003
	WM_MOVING = 0x0216
	WM_NCACTIVATE = 0x0086
	WM_NCCALCSIZE = 0x0083
	WM_NCCREATE = 0x0081
	WM_NCDESTROY = 0x0082
	WM_NULL = 0x0000
	WM_QUERYDRAGICON = 0x0037
	WM_QUERYOPEN = 0x0013
	WM_QUIT = 0x0012
	WM_SHOWWINDOW = 0x0018
	WM_SIZE = 0x0005
	WM_SIZING = 0x0214
	WM_STYLECHANGED = 0x007D
	WM_STYLECHANGING = 0x007C
//	WM_THEMECHANGED = 0x031A		// Windows XP and newer only
//	WM_USERCHANGED = 0x0054			// Windows XP only: [Note  This message is not supported as of Windows Vista.; also listed as not supported by server Windows]
	WM_WINDOWPOSCHANGED = 0x0047
	WM_WINDOWPOSCHANGING = 0x0046
)

type MINMAXINFO struct {
	PtReserved		POINT
	PtMaxSize			POINT
	PtMaxPosition		POINT
	PtMinTrackSize		POINT
	PtMaxTrackSize	POINT
}

func (l LPARAM) MINMAXINFO() *MINMAXINFO {
	return (*MINMAXINFO)(unsafe.Pointer(l))
}
