//go:build js && wasm

package platform

import (
	"io/fs"
	"syscall/js"

	"golang-rpg/assets"
	sceneconfig "golang-rpg/config"
)

var ready bool

func Browser() bool { return true }
func Controls() Input {
	value := js.Global().Get("appInput")
	if value.Type() != js.TypeObject {
		return Input{}
	}
	pressed := func(key string) bool { v := value.Get(key); return v.Type() == js.TypeBoolean && v.Bool() }
	input := Input{Left: pressed("left"), Right: pressed("right"), Up: pressed("up"), Down: pressed("down"), Attack: pressed("attack"), Use: pressed("use")}
	value.Set("use", false)
	return input
}
func Ready(width, height int) {
	if !ready {
		callback := js.Global().Get("appReady")
		if callback.Type() == js.TypeFunction {
			callback.Invoke(width, height)
			ready = true
		}
	}
}
func ReportError(message string) {
	callback := js.Global().Get("appError")
	if callback.Type() == js.TypeFunction {
		callback.Invoke(message)
	}
}
func Resources() (fs.FS, fs.FS) { return sceneconfig.Files(), assets.BrowserFiles() }
