package webkit2

import (
	"github.com/auroralaboratories/gotk3/gtk"
	"runtime"
)

func init() {
	runtime.LockOSThread()
	gtk.Init(nil)
}
