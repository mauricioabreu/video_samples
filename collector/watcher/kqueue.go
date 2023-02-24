//go:build darwin && !kqueue && cgo
// +build darwin,!kqueue,cgo

package watcher

import (
	"github.com/rjeczalik/notify"
)

const (
	WriteEvent = notify.Event(notify.FSEventsModified)
)
