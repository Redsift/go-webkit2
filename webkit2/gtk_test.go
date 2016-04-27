package webkit2

import (
	"runtime"

	"github.com/auroralaboratories/gotk3/gtk"
)

func init() {
	runtime.LockOSThread()
	gtk.Init(nil)
}
