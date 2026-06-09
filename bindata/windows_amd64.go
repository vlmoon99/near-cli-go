//go:build windows && amd64
// +build windows,amd64

package bindata

import (
	_ "embed"
)

//go:embed tools/windows_amd64/near.exe
var NearCli []byte

//go:embed tools/windows_amd64/tinygo.zip
var TinyGoZip []byte
