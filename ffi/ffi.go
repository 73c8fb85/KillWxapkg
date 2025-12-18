package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"runtime"
	"runtime/debug"
	"sync"
	"unsafe"

	"github.com/Ackites/KillWxapkg/cmd"
	. "github.com/Ackites/KillWxapkg/internal/cmd"
	. "github.com/Ackites/KillWxapkg/internal/config"
	"github.com/Ackites/KillWxapkg/internal/decrypt"
	"github.com/Ackites/KillWxapkg/internal/restore"
	"github.com/Ackites/KillWxapkg/internal/unpack"
)

// Result FFI返回结构
type Result struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func resultJSON(success bool, message string, data string) *C.char {
	r := Result{Success: success, Message: message, Data: data}
	b, _ := json.Marshal(r)
	return C.CString(string(b))
}

//export DecryptWxapkg
func DecryptWxapkg(inputFile, appID *C.char) *C.char {
	data, err := decrypt.DecryptWxapkg(C.GoString(inputFile), C.GoString(appID))
	if err != nil {
		return resultJSON(false, err.Error(), "")
	}
	return resultJSON(true, "解密成功", string(data))
}

//export UnpackWxapkg
func UnpackWxapkg(inputFile, appID, outputDir *C.char) *C.char {
	inFile := C.GoString(inputFile)
	aid := C.GoString(appID)
	outDir := C.GoString(outputDir)

	data, err := decrypt.DecryptWxapkg(inFile, aid)
	if err != nil {
		return resultJSON(false, "解密失败: "+err.Error(), "")
	}

	files, err := unpack.UnpackWxapkg(data, outDir)
	if err != nil {
		return resultJSON(false, "解包失败: "+err.Error(), "")
	}

	filesJSON, _ := json.Marshal(files)
	return resultJSON(true, "解包成功", string(filesJSON))
}

//export ProcessWxapkg
func ProcessWxapkg(inputFile, appID, outputDir *C.char, restoreProject, pretty, save C.int) *C.char {
	inFile := C.GoString(inputFile)
	aid := C.GoString(appID)
	outDir := C.GoString(outputDir)

	configManager := NewSharedConfigManager()
	configManager.Set("appID", aid)
	configManager.Set("pretty", pretty == 1)
	configManager.Set("noClean", false)
	configManager.Set("sensitive", false)

	if outDir == "" {
		outDir = DetermineOutputDir(inFile, aid)
	}

	err := ProcessFile(inFile, outDir, aid, save == 1)
	if err != nil {
		return resultJSON(false, err.Error(), "")
	}

	if restoreProject == 1 {
		restore.ProjectStructure(outDir, true)
	}

	return resultJSON(true, "处理成功", outDir)
}

//export ProcessWxapkgBatch
func ProcessWxapkgBatch(input, appID, outputDir, fileExt *C.char, restoreProject, pretty, noClean, save C.int) *C.char {
	inPath := C.GoString(input)
	aid := C.GoString(appID)
	outDir := C.GoString(outputDir)
	ext := C.GoString(fileExt)
	if ext == "" {
		ext = ".wxapkg"
	}

	cmd.Execute(aid, inPath, outDir, ext, restoreProject == 1, pretty == 1, noClean == 1, save == 1, false)
	return resultJSON(true, "批量处理完成", outDir)
}

//export FreeString
func FreeString(s *C.char) {
	C.free(unsafe.Pointer(s))
}

var initOnce sync.Once

//export InitLibrary
func InitLibrary() {
	initOnce.Do(func() {
		// 初始化操作
	})
}

//export FreeMemory
func FreeMemory() {
	runtime.GC()
	debug.FreeOSMemory()
}

func main() {}
