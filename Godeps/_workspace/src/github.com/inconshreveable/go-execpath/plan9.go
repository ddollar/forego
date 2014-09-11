// +build plan9

package execpath

import (
	"fmt"
)

func GetNative() (string, error) {
	return "", fmt.Errorf("GetNative() not implemented on plan9")
}
