#include "webview.h"

void webkit2_gasync_callback(GObject *source_object, GAsyncResult *res, gpointer user_data) {
    go_genericGAsyncCallback(source_object, res, user_data);
}
