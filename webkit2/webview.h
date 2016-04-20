#ifndef GO_WEBVIEW_H
#define GO_WEBVIEW_H

#include <stdlib.h>
#include <webkit2/webkit2.h>
#include <cairo/cairo.h>

static WebKitWebView* to_WebKitWebView(GtkWidget* w) { return WEBKIT_WEB_VIEW(w); }

// callback declarations
void webkit2_gasync_callback(GObject *source_object, GAsyncResult *res, gpointer user_data);

#endif
