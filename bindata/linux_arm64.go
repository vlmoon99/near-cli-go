//go:build linux && arm64
// +build linux,arm64

package bindata

import (
	_ "embed"
)

//go:embed tools/linux_arm64/near
var NearCli []byte
