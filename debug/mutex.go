//go:build !debug

package debug

import (
	"sync"
)

type Mutex = sync.Mutex
