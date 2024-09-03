//go:build !debug

package debug

import (
	"sync"
)

type RWMutex = sync.RWMutex
