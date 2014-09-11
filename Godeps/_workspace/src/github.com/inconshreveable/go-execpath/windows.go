// +build windows

package execpath

import (
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

func GetNative() (path string, err error) {
	h, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		return
	}
	defer syscall.FreeLibrary(h)

	addr, err := syscall.GetProcAddress(h, "GetModuleFileNameW")
	if err != nil {
		return
	}

	bufSize := uint32(maxPathSize)
	for {
		buf := make([]uint16, bufSize)
		retptr, _, errno := syscall.Syscall(addr, 3,
			/* NULL */ uintptr(0),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&bufSize)))

		if errno != 0 {
			err = errno
			return
		}

		retval := uint32(retptr)
		if retval == 0 {
			err = fmt.Errorf("GetModuleFileName failed to get executable path")
			return
		} else if retval == bufSize {
			// buffer is too small, double it and try again
			bufSize = bufSize * 2
		} else {
			// success
			path = string(utf16.Decode(buf[0:retval]))
			return
		}
	}
}
