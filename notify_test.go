package sche

import (
	"context"
	"testing"
	"time"
)

func TestNotify(t *testing.T) {
	notify := Notify{ctx: context.Background()}
	notify.Wait(time.Now())
	notify.Notify()
}
