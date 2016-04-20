package webkit2

// #include "webview.h"
// #cgo pkg-config: webkit2gtk-4.0
import "C"


import (
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"
    "github.com/gotk3/gotk3/cairo"
    "github.com/sqs/gojs"
    "image"
    "unsafe"
)

//export go_genericGAsyncCallback
func go_genericGAsyncCallback(source *C.GObject, result *C.GAsyncResult, callbackId *C.char) {
    key := C.GoString(callbackId)

    if obj, ok := cgoget(key); ok {
        switch obj.(type) {
        case *RunJavaScriptResponse:
            var jserr *C.GError

            response := obj.(*RunJavaScriptResponse)

            if response.Autoremove {
                defer cgounregister(key)
            }

            if jsResult := C.webkit_web_view_run_javascript_finish(response.CWebView, result, &jserr); jsResult == nil {
                defer C.g_error_free(jserr)
                msg := C.GoString((*C.char)(jserr.message))
                response.Reply(nil, errors.New(msg))
            } else {
                ctxRaw := gojs.RawGlobalContext(unsafe.Pointer(C.webkit_javascript_result_get_global_context(jsResult)))
                jsValRaw := gojs.RawValue(unsafe.Pointer(C.webkit_javascript_result_get_value(jsResult)))
                ctx := (*gojs.Context)(gojs.NewGlobalContextFrom(ctxRaw))
                jsVal := ctx.NewValueFrom(jsValRaw)
                response.Reply(jsVal, nil)
            }

        case *GetSnapshotAsImageResponse:
            var snapErr *C.GError

            response := obj.(*GetSnapshotAsImageResponse)

            if response.Autoremove {
                defer cgounregister(key)
            }

            if snapResult := C.webkit_web_view_get_snapshot_finish(response.CWebView, result, &snapErr); snapResult == nil {
                defer C.g_error_free(snapErr)
                msg := C.GoString((*C.char)(snapErr.message))
                response.Reply(nil, errors.New(msg))
            } else {
                defer C.cairo_surface_destroy(snapResult)

                if C.cairo_surface_get_type(snapResult) != cairoSurfaceTypeImage ||
                    C.cairo_image_surface_get_format(snapResult) != cairoImageSurfaceFormatARGB32 {
                    response.Reply(nil, errors.New("Snapshot in unexpected format"))
                    return
                }

                w := int(C.cairo_image_surface_get_width(snapResult))
                h := int(C.cairo_image_surface_get_height(snapResult))
                stride := int(C.cairo_image_surface_get_stride(snapResult))
                data := unsafe.Pointer(C.cairo_image_surface_get_data(snapResult))
                surfaceBytes := C.GoBytes(data, C.int(stride*h))

                // convert from b,g,r,a or a,r,g,b(local endianness) to r,g,b,a
                testint, _ := binary.ReadUvarint(bytes.NewBuffer([]byte{0x1, 0}))

                if testint == 0x1 {
                    // Little: b,g,r,a -> r,g,b,a
                    for i := 0; i < w*h; i++ {
                        b := surfaceBytes[4*i+0]
                        r := surfaceBytes[4*i+2]
                        surfaceBytes[4*i+0] = r
                        surfaceBytes[4*i+2] = b
                    }
                } else {
                    // Big: a,r,g,b -> r,g,b,a
                    for i := 0; i < w*h; i++ {
                        a := surfaceBytes[4*i+0]
                        r := surfaceBytes[4*i+1]
                        g := surfaceBytes[4*i+2]
                        b := surfaceBytes[4*i+3]
                        surfaceBytes[4*i+0] = r
                        surfaceBytes[4*i+1] = g
                        surfaceBytes[4*i+2] = b
                        surfaceBytes[4*i+3] = a
                    }
                }

                rgba := &image.RGBA{
                    Pix:    surfaceBytes,
                    Stride: stride,
                    Rect:   image.Rect(0, 0, w, h),
                }

                response.Reply(rgba, nil)
            }

        case *GetSnapshotAsCairoSurfaceResponse:
            var snapErr *C.GError

            response := obj.(*GetSnapshotAsCairoSurfaceResponse)

            snapResult := C.webkit_web_view_get_snapshot_finish(response.CWebView, result, &snapErr)

            if snapResult == nil {
                defer C.g_error_free(snapErr)
                msg := C.GoString((*C.char)(snapErr.message))
                response.Reply(nil, errors.New(msg))
            } else {
                surface := cairo.NewSurface(uintptr(unsafe.Pointer(snapResult)), false)

                if status := surface.Status(); status == cairo.STATUS_SUCCESS {
                    response.Reply(surface, nil)
                } else {
                    response.Reply(nil, fmt.Errorf("Cairo surface error %d", status))
                }
            }
        }
    }
}
