//go:build linux
// +build linux

package watcher

import (
	"github.com/rjeczalik/notify"
)

const (
	WriteEvent = notify.Event(notify.InCloseWrite)
)
