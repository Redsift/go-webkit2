// +build !webkit2gtk_2_4_0,!webkit2gtk_2_6_0
package webkit2

// #include <stdlib.h>
// #include <webkit2/webkit2.h>
//
// #cgo pkg-config: webkit2gtk-4.0
import "C"

import (
	"github.com/gotk3/gotk3/gdk"
	"unsafe"
)

const (
	SnapshotOptionTransparentBackground = C.WEBKIT_SNAPSHOT_OPTIONS_TRANSPARENT_BACKGROUND
)

func (v *WebView) SetBackgroundColor(color *gdk.RGBA) {
	C.webkit_web_view_set_background_color(v.webView, (*C.GdkRGBA)(unsafe.Pointer(color.Native())))
}
