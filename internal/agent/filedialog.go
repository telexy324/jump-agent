//go:build windows

package agent

import (
	"syscall"
	"unsafe"
)

var (
	ole32            = syscall.NewLazyDLL("ole32.dll")
	coInitializeEx   = ole32.NewProc("CoInitializeEx")
	coCreateInstance = ole32.NewProc("CoCreateInstance")
	coTaskMemFree    = ole32.NewProc("CoTaskMemFree")

	CLSID_FileOpenDialog = syscall.GUID{0xDC1C5A9C, 0xE88A, 0x4DDE, [8]byte{0xA5, 0xA1, 0x60, 0xF8, 0x2A, 0x20, 0xAE, 0xF7}}
	IID_IFileOpenDialog  = syscall.GUID{0xD57C7288, 0xD4AD, 0x4768, [8]byte{0xBE, 0x02, 0x9D, 0x96, 0x95, 0x32, 0xD9, 0x60}}
	IID_IShellItem       = syscall.GUID{0x43826D1E, 0xE718, 0x42EE, [8]byte{0xBC, 0x55, 0xA1, 0xE2, 0x61, 0xC3, 0x7B, 0xFE}}
)

const (
	COINIT_APARTMENTTHREADED = 0x2
	CLSCTX_INPROC_SERVER     = 1
	SIGDN_FILESYSPATH        = 0x80058000
)

func utf16PtrToString(p *uint16) string {
	if p == nil {
		return ""
	}
	// read until 0
	var s []uint16
	for i := 0; ; i++ {
		c := *(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + uintptr(i*2)))
		if c == 0 {
			break
		}
		s = append(s, c)
	}
	return syscall.UTF16ToString(s)
}

func SelectExecutable() (string, error) {
	coInitializeEx.Call(0, COINIT_APARTMENTTHREADED)

	var dialog uintptr
	hr, _, _ := coCreateInstance.Call(
		uintptr(unsafe.Pointer(&CLSID_FileOpenDialog)),
		0,
		CLSCTX_INPROC_SERVER,
		uintptr(unsafe.Pointer(&IID_IFileOpenDialog)),
		uintptr(unsafe.Pointer(&dialog)),
	)
	if hr != 0 {
		return "", syscall.Errno(hr)
	}

	vtable := **(**uintptr)(unsafe.Pointer(&dialog))

	// Show() -> vtable[3]
	show := *(*uintptr)(unsafe.Pointer(vtable + 3*unsafe.Sizeof(uintptr(0))))
	hr, _, _ = syscall.SyscallN(show, dialog, 0)
	if hr != 0 {
		return "", syscall.Errno(hr)
	}

	// GetResult -> vtable[20]
	getResult := *(*uintptr)(unsafe.Pointer(vtable + 20*unsafe.Sizeof(uintptr(0))))

	var item uintptr
	hr, _, _ = syscall.SyscallN(getResult, dialog, uintptr(unsafe.Pointer(&item)))
	if hr != 0 {
		return "", syscall.Errno(hr)
	}

	itemVTable := **(**uintptr)(unsafe.Pointer(&item))

	// GetDisplayName -> vtable[5]
	getDisplayName := *(*uintptr)(unsafe.Pointer(itemVTable + 5*unsafe.Sizeof(uintptr(0))))

	var pszPath uintptr
	hr, _, _ = syscall.SyscallN(getDisplayName, item, SIGDN_FILESYSPATH, uintptr(unsafe.Pointer(&pszPath)))
	if hr != 0 {
		return "", syscall.Errno(hr)
	}

	// ↓ 正确读取 UTF16 PWSTR
	path := utf16PtrToString((*uint16)(unsafe.Pointer(pszPath)))

	// ↓ 必须释放！
	coTaskMemFree.Call(pszPath)

	return path, nil
}
