//go:build darwin && arm64
// +build darwin,arm64

package bindata

import (
	_ "embed"
)

//go:embed tools/darwin_arm64/near
var NearCli []byte
