package scan

import (
	"context"
	"testing"
)

func TestScan(t *testing.T) {
	ctx := context.Background()
	hostinfo, err := NewScanner("")
	if err != nil {
		t.Error(err)
	}
	err = hostinfo.ScanSubnet(ctx)
	if err != nil {
		t.Error(err)
	}
}
