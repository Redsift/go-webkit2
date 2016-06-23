package webkit2

// #include <webkit2/webkit2.h>
// #include "arrays.h"
import "C"
import "runtime"

// WebContext manages all aspects common to all WebViews.
//
// See also: WebKitWebContext at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebContext.html.
type WebContext struct {
	webContext *C.WebKitWebContext

	// Book keeping for the C allocations
	languageGCharArray **C.gchar
}

// DefaultWebContext returns the default WebContext.
//
// See also: webkit_web_context_get_default at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebContext.html#webkit-web-context-get-default.
func DefaultWebContext() *WebContext {
	wc := &WebContext{C.webkit_web_context_get_default(), nil}
	runtime.SetFinalizer(wc, (*WebContext).Free)
	return wc
}

// CacheModel describes the caching behavior.
//
// See also: WebKitCacheModel at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebContext.html#WebKitCacheModel.
type CacheModel int

// CacheModel enum values are described at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebContext.html#WebKitCacheModel.
const (
	DocumentViewerCacheModel CacheModel = iota
	WebBrowserCacheModel
	DocumentBrowserCacheModel
)

// CacheModel returns the current cache model.
//
// See also: webkit_web_context_get_cache_model at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebContext.html#webkit-web-context-get-cache-model.
func (wc *WebContext) CacheModel() CacheModel {
	return CacheModel(C.int(C.webkit_web_context_get_cache_model(wc.webContext)))
}

// SetCacheModel sets the current cache model.
//
// See also: webkit_web_context_set_cache_model at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebContext.html#webkit-web-context-set-cache-model.
func (wc *WebContext) SetCacheModel(model CacheModel) {
	C.webkit_web_context_set_cache_model(wc.webContext, C.WebKitCacheModel(model))
}

// ClearCache clears all resources currently cached.
//
// See also: webkit_web_context_clear_cache at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebContext.html#webkit-web-context-clear-cache.
func (wc *WebContext) ClearCache() {
	C.webkit_web_context_clear_cache(wc.webContext)
}

// SetPreferredLanguages set the list of preferred languages, sorted from most desirable to least desirable.
//
// See also: webkit_web_context_set_preferred_languages at
// http://webkitgtk.org/reference/webkit2gtk/stable/WebKitWebContext.html#webkit-web-context-set-preferred-languages.
func (wc *WebContext) SetPreferredLanguages(languages []string) {
	wc.freeLanguageGCharArray()
	wc.languageGCharArray = C.alloc_gchar_array((C.size_t)(len(languages) + 1))


	for i, s := range languages {
		cstr := C.CString(s)
		C.set_gchar_array(wc.languageGCharArray, C.int(i), (*C.gchar)(cstr))
	}

	C.set_gchar_array(wc.languageGCharArray, C.int(len(languages)), (*C.gchar)(nil))

	C.webkit_web_context_set_preferred_languages(wc.webContext, wc.languageGCharArray)
}

func (wc *WebContext) Free() {
	wc.freeLanguageGCharArray()
}

func (wc *WebContext) freeLanguageGCharArray() {
	if wc.languageGCharArray == nil {
		return
	}

	C.free_gchar_array(wc.languageGCharArray)
	wc.languageGCharArray = nil
}

