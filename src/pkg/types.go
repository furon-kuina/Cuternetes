package c8s

import "github.com/docker/docker/api/types/strslice"

type Container struct {
	Image string
	Cmd   strslice.StrSlice `json:",omitempty"`
}
