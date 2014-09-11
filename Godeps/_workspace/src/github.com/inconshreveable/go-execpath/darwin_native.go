// +build darwin,cgo

package execpath

import (
	"fmt"
)

/*
#include <mach-o/dyld.h>
#include <string.h>
*/
import "C"

func GetNative() (string, error) {
	var buflen C.uint32_t = C.uint32_t(maxPathSize)
	buf := make([]C.char, buflen)

	ret := C._NSGetExecutablePath(&buf[0], &buflen)
	if ret == -1 {
		// buflen wasn't large enough, _NSGetExecutablePath set it to the necessary size
		// so recreate the buffer and try again
		buf = make([]C.char, buflen)
		ret = C._NSGetExecutablePath(&buf[0], &buflen)
		if ret == -1 {
			// this should never happen
			return "", fmt.Errorf("_NSGetExecutable failed to get the executable path")
		}
	}
	pathlen := C.strlen(&buf[0])
	return C.GoStringN(&buf[0], C.int(pathlen)), nil
}
