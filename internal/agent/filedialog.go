//go:build windows
// +build windows

package agent

import (
	"syscall"
	"unsafe"
)

var (
	modOle32             = syscall.NewLazyDLL("ole32.dll")
	procCoInitializeEx   = modOle32.NewProc("CoInitializeEx")
	procCoCreateInstance = modOle32.NewProc("CoCreateInstance")
)

const (
	COINIT_APARTMENTTHREADED = 0x2
	CLSCTX_INPROC_SERVER     = 1
)

var (
	CLSID_FileOpenDialog = syscall.GUID{0xDC1C5A9C, 0xE88A, 0x4DDE, [8]byte{0xA5, 0xA1, 0x60, 0xF8, 0x2A, 0x20, 0xAE, 0xF7}}
	IID_IFileOpenDialog  = syscall.GUID{0xD57C7288, 0xD4AD, 0x4768, [8]byte{0xBE, 0x02, 0x9D, 0x96, 0x95, 0x32, 0xD9, 0x60}}
	IID_IShellItem       = syscall.GUID{0x43826D1E, 0xE718, 0x42EE, [8]byte{0xBC, 0x55, 0xA1, 0xE2, 0x61, 0xC3, 0x7B, 0xFE}}
)

func SelectExecutable() (string, error) {
	// 初始化 COM
	hr, _, _ := procCoInitializeEx.Call(0, uintptr(COINIT_APARTMENTTHREADED))
	if hr != 0 {
		return "", syscall.Errno(hr)
	}

	var pfd uintptr
	hr, _, _ = procCoCreateInstance.Call(
		uintptr(unsafe.Pointer(&CLSID_FileOpenDialog)),
		0,
		uintptr(CLSCTX_INPROC_SERVER),
		uintptr(unsafe.Pointer(&IID_IFileOpenDialog)),
		uintptr(unsafe.Pointer(&pfd)),
	)

	if hr != 0 {
		return "", syscall.Errno(hr)
	}

	// 调用 Show()，弹出选择窗口
	vtable := *(*uintptr)(unsafe.Pointer(pfd))
	show := *(*uintptr)(unsafe.Pointer(vtable + 5*unsafe.Sizeof(uintptr(0))))

	hr, _, _ = syscall.SyscallN(show, pfd, 0)
	if hr != 0 {
		return "", syscall.Errno(hr)
	}

	// 获取结果
	getResult := *(*uintptr)(unsafe.Pointer(vtable + 7*unsafe.Sizeof(uintptr(0))))
	var psi uintptr
	_, _, _ = syscall.SyscallN(getResult, pfd, uintptr(unsafe.Pointer(&psi)))

	// 取文件路径
	siTable := *(*uintptr)(unsafe.Pointer(psi))
	getDisplayName := *(*uintptr)(unsafe.Pointer(siTable + 5*unsafe.Sizeof(uintptr(0))))

	var pszPath uintptr
	_, _, _ = syscall.SyscallN(getDisplayName, psi, 0, uintptr(unsafe.Pointer(&pszPath)))

	path := syscall.UTF16ToString((*[1 << 16]uint16)(unsafe.Pointer(pszPath))[:])

	return path, nil
}
