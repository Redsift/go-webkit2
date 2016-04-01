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

// func (v *WebView) GetBackgroundColor(color gdk.RGBA) (gdk.RGBA, error) {
// 	rgba := &C.GdkRGBA{}

// 	C.webkit_web_view_get_background_color(v.webView, unsafe.Pointer(rgba))
// }
