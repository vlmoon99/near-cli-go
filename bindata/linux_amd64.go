//go:build linux && amd64
// +build linux,amd64

package bindata

import (
	_ "embed"
)

//go:embed tools/linux_amd64/near
var NearCli []byte
