package webkit2

// #cgo pkg-config: webkit2gtk-4.0
// #include "webview.h"
//
// static WebKitWebView* to_WebKitWebView(GtkWidget* w) { return WEBKIT_WEB_VIEW(w); }
//
import "C"

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/sqs/gojs"
	"image"
	"unsafe"
)

type SnapshotOptions int

const (
	SnapshotOptionNone                         SnapshotOptions = C.WEBKIT_SNAPSHOT_OPTIONS_NONE
	SnapshotOptionIncludeSelectionHighlighting                 = C.WEBKIT_SNAPSHOT_OPTIONS_INCLUDE_SELECTION_HIGHLIGHTING
)

type SnapshotRegion int

const (
	RegionVisible      SnapshotRegion = C.WEBKIT_SNAPSHOT_REGION_VISIBLE
	RegionFullDocument                = C.WEBKIT_SNAPSHOT_REGION_FULL_DOCUMENT
)

// WebView represents a WebKit WebView.
//
// See also: WebView at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebView.html.
type WebView struct {
	*gtk.Widget
	webView *C.WebKitWebView
}

type RunJavaScriptResponse struct {
	CWebView   *C.WebKitWebView
	Reply      func(result *gojs.Value, err error)
	Autoremove bool
}

type GetSnapshotAsImageResponse struct {
	CWebView   *C.WebKitWebView
	Reply      func(result *image.RGBA, err error)
	Autoremove bool
}

type GetSnapshotAsCairoSurfaceResponse struct {
	CWebView   *C.WebKitWebView
	Reply      func(result *cairo.Surface, err error)
	Autoremove bool
}

// NewWebView creates a new WebView with the default WebContext and the default
// WebViewGroup.
//
// See also: webkit_web_view_new at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebView.html#webkit-web-view-new.
func NewWebView() *WebView {
	return newWebView(C.webkit_web_view_new())
}

// NewWebViewWithContext creates a new WebView with the given WebContext and the
// default WebViewGroup.
//
// See also: webkit_web_view_new_with_context at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebView.html#webkit-web-view-new-with-context.
func NewWebViewWithContext(ctx *WebContext) *WebView {
	return newWebView(C.webkit_web_view_new_with_context(ctx.webContext))
}

func newWebView(webViewWidget *C.GtkWidget) *WebView {
	obj := &glib.Object{glib.ToGObject(unsafe.Pointer(webViewWidget))}
	return &WebView{&gtk.Widget{glib.InitiallyUnowned{obj}}, C.to_WebKitWebView(webViewWidget)}
}

// Context returns the current WebContext of the WebView.
//
// See also: webkit_web_view_get_context at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebView.html#webkit-web-view-get-context.
func (v *WebView) Context() *WebContext {
	return &WebContext{C.webkit_web_view_get_context(v.webView)}
}

// LoadURI requests loading of the specified URI string.
//
// See also: webkit_web_view_load_uri at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebView.html#webkit-web-view-load-uri
func (v *WebView) LoadURI(uri string) {
	C.webkit_web_view_load_uri(v.webView, (*C.gchar)(C.CString(uri)))
}

// LoadHTML loads the given content string with the specified baseURI. The MIME
// type of the document will be "text/html".
//
// See also: webkit_web_view_load_html at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebView.html#webkit-web-view-load-html
func (v *WebView) LoadHTML(content, baseURI string) {
	C.webkit_web_view_load_html(v.webView, (*C.gchar)(C.CString(content)), (*C.gchar)(C.CString(baseURI)))
}

// Settings returns the current active settings of this WebView's WebViewGroup.
//
// See also: webkit_web_view_get_settings at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebView.html#webkit-web-view-get-settings.
func (v *WebView) Settings() *Settings {
	return newSettings(C.webkit_web_view_get_settings(v.webView))
}

// Title returns the current active title of the WebView.
//
// See also: webkit_web_view_get_title at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebView.html#webkit-web-view-get-title.
func (v *WebView) Title() string {
	return C.GoString((*C.char)(C.webkit_web_view_get_title(v.webView)))
}

// URI returns the current active URI of the WebView.
//
// See also: webkit_web_view_get_uri at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebView.html#webkit-web-view-get-uri.
func (v *WebView) URI() string {
	return C.GoString((*C.char)(C.webkit_web_view_get_uri(v.webView)))
}

// JavaScriptGlobalContext returns the global JavaScript context used by
// WebView.
//
// See also: webkit_web_view_get_javascript_global_context at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebView.html#webkit-web-view-get-javascript-global-context
func (v *WebView) JavaScriptGlobalContext() *gojs.Context {
	return (*gojs.Context)(gojs.NewGlobalContextFrom((gojs.RawGlobalContext)(unsafe.Pointer(C.webkit_web_view_get_javascript_global_context(v.webView)))))
}

// RunJavaScript runs script asynchronously in the context of the current page
// in the WebView. Upon completion, resultCallback will be called with the
// result of evaluating the script, or with an error encountered during
// execution. To get the stack trace and other error logs, use the
// ::console-message signal.
//
// See also: webkit_web_view_run_javascript at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebView.html#webkit-web-view-run-javascript
func (v *WebView) RunJavaScript(script string, resultCallback func(result *gojs.Value, err error)) {
	var callbackIdPtr unsafe.Pointer

	if resultCallback != nil {
		callbackId := cgoregister(``, &RunJavaScriptResponse{
			CWebView:   v.webView,
			Reply:      resultCallback,
			Autoremove: true,
		})

		callbackIdPtr = unsafe.Pointer(C.CString(callbackId))
	}

	C.webkit_web_view_run_javascript(v.webView,
		(*C.gchar)(C.CString(script)),
		nil,
		(C.GAsyncReadyCallback)(C.webkit2_gasync_callback),
		callbackIdPtr)
}

// Destroy destroys the WebView's corresponding GtkWidget and marks its internal
// WebKitWebView as nil so that it can't be accidentally reused.
func (v *WebView) Destroy() {
	v.Widget.Destroy()
	v.webView = nil
}

// LoadEvent denotes the different events that happen during a WebView load
// operation.
//
// See also: WebKitLoadEvent at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebView.html#WebKitLoadEvent.
type LoadEvent int

// LoadEvent enum values are described at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebView.html#WebKitLoadEvent.
const (
	LoadStarted LoadEvent = iota
	LoadRedirected
	LoadCommitted
	LoadFinished
)

// http://cairographics.org/manual/cairo-cairo-surface-t.html#cairo-surface-type-t
const cairoSurfaceTypeImage = 0

// http://cairographics.org/manual/cairo-Image-Surfaces.html#cairo-format-t
const cairoImageSurfaceFormatARGB32 = 0

// GetSnapshot runs asynchronously, taking a snapshot of the WebView.
// Upon completion, resultCallback will be called with a copy of the underlying
// bitmap backing store for the frame, or with an error encountered during
// execution.
//
// See also: webkit_web_view_get_snapshot at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebView.html#webkit-web-view-get-snapshot
func (v *WebView) GetSnapshotWithOptions(region SnapshotRegion, options SnapshotOptions, resultCallback func(result *image.RGBA, err error)) {
	var callbackIdPtr unsafe.Pointer

	if resultCallback != nil {
		callbackId := cgoregister(``, &GetSnapshotAsImageResponse{
			CWebView:   v.webView,
			Reply:      resultCallback,
			Autoremove: true,
		})

		callbackIdPtr = unsafe.Pointer(C.CString(callbackId))
	}

	C.webkit_web_view_get_snapshot(v.webView,
		(C.WebKitSnapshotRegion)(region), // FullDocument is the only working region at this point
		(C.WebKitSnapshotOptions)(options),
		nil,
		(C.GAsyncReadyCallback)(C.webkit2_gasync_callback),
		callbackIdPtr)
}

func (v *WebView) GetSnapshot(resultCallback func(result *image.RGBA, err error)) {
	v.GetSnapshotWithOptions(RegionFullDocument, SnapshotOptionNone, resultCallback)
}

func (v *WebView) GetSnapshotSurfaceWithOptions(region SnapshotRegion, options SnapshotOptions, resultCallback func(result *cairo.Surface, err error)) {
	var callbackIdPtr unsafe.Pointer

	if resultCallback != nil {
		callbackId := cgoregister(``, &GetSnapshotAsCairoSurfaceResponse{
			CWebView:   v.webView,
			Reply:      resultCallback,
			Autoremove: true,
		})

		callbackIdPtr = unsafe.Pointer(C.CString(callbackId))
	}

	C.webkit_web_view_get_snapshot(v.webView,
		(C.WebKitSnapshotRegion)(region), // FullDocument is the only working region at this point
		(C.WebKitSnapshotOptions)(options),
		nil,
		(C.GAsyncReadyCallback)(C.webkit2_gasync_callback),
		callbackIdPtr)
}
