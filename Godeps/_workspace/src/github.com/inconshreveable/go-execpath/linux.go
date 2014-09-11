// +build linux

package execpath

import (
	"os"
)

// GetNative returns the executable path by calling a platform-specific API. The API
// used for each platform is specified as follows:
//
// - Windows: GetModuleFileName
// - Linux: /proc/self/exe
// - OS X: _NSGetExecutablePath
// - FreeBSD: /proc/curproc/file
// - OpenBSD: /proc/curproc/file
// - NetBSD: /proc/curproc/exe
func GetNative() (string, error) {
	return os.Readlink("/proc/self/exe")
}
