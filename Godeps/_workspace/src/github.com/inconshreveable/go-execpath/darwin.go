// +build darwin,!cgo

package execpath

import (
	"fmt"
)

func GetNative() (string, error) {
	return "", fmt.Errorf("GetNative() not implemented when cross-compiling to darwin")
}
