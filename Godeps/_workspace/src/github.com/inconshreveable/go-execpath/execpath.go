/*
Package execpath provides APIs for returning the absolute path to the
executable file of the running program.

	exePath, err := execpath.Get()
	if err != nil {
		fmt.Println(err.Error())
	}

	// Prints /path/to/your/executable
	fmt.Printf(exePath)

Determining the path to the executable of the currently running program is
difficult. Many programs read os.Args[0], sometimes combining it with the
working directory to determine the path, but this approach is fragile and
yields incorrect results because:

- The program may be in the PATH, and invoked without specifying any path to the executable.

- The program may have changed directories after it was invoked.

- The value of os.Args[0] is not required to be set to the invoking path.
Most shells follow a convention of setting os.Args[0] to the path used to
invoke the executable, but your program may be invoked by an unconventional
shell, forked from a non-shell process, started as a Windows service, etc.

Each operating system has a platform-specific API that is guaranteed to return
a proper path to the executable file. Package execpath provides GetNative(),
which uses the following platform-specific APIs which will always return the
proper executable path, if successful:

- Windows: GetModuleFileName

- Linux: /proc/self/exe

- OS X: _NSGetExecutablePath

- FreeBSD: /proc/curproc/file

- OpenBSD: /proc/curproc/file

- NetBSD: /proc/curproc/exe

Package execpath also implements the naive heuristic methods which involve
examining os.Args[0] or checking directories in the PATH.

Lastly, execpath provides Get(), the API you should prefer, which calls
GetNative() and falls back to using heuristic methods if GetNative() fails.
*/
package execpath

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

var maxPathSize = 1024

// Get returns the path to the current executable or an error on failure.
// This is the preferred API because it attempts to use the native
// method and falls back to heuristic methods on failure. In order, it calls:
//
// GetNative()
//
// GetArg0() if GetNative() fails
//
// GetPath() if GetArg0() fails
func Get() (path string, err error) {
	if path, err = GetNative(); err == nil {
		return
	}

	if path, err = GetArg0(); err == nil {
		return
	}

	if path, err = GetPath(); err == nil {
		return
	}

	return
}

// returns true if the path exists, false otherwise
// an error is returned on failure (e.g. insufficient
// permissions)
func pathExists(p string) (bool, error) {
	_, err := os.Stat(p)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func makeAbsolute(p string) (fullPath string, err error) {
	if path.IsAbs(p) {
		//it's an absolute path, we're done
		fullPath = p
		return
	}

	// the path is relative, combine p with the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	fullPath = path.Join(wd, p)
	return
}

// GetArg0 returns the executable path by examining the value of os.Args[0],
// combining it with the working directory if necessary. Because os.Args[0] is
// set by convention, the returned path is not guaranteed to be correct, even
// if GetArg0 does not return an error.
func GetArg0() (p string, err error) {
	p, err = makeAbsolute(os.Args[0])
	if err != nil {
		return
	}

	var exists bool
	if exists, err = pathExists(p); !exists {
		err = fmt.Errorf("Can't determine executable path from os.Args[0]")
	}

	return
}

// GetPath returns the executable path by searching for the the executable name
// (os.Args[0]) in each directory of the PATH environment variable. Because
// os.Args[0] is set by convention, the returned path is not guaranteed to be
// correct, even if GetPath does not return an error.
func GetPath() (p string, err error) {
	return exec.LookPath(os.Args[0])
}
