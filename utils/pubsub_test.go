package utils

import (
	"context"
	"testing"
)

func TestPubSub_Init(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		ps   *PubSub
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ps.Init(tt.args.ctx)
		})
	}
}
