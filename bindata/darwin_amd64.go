//go:build darwin && amd64
// +build darwin,amd64

package bindata

import (
	_ "embed"
)

//go:embed tools/darwin_amd64/near
var NearCli []byte
