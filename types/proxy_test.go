package types

import (
	"context"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestProxy_Init(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		p    *Proxy
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Init(tt.args.ctx, false, 9040)
		})
	}
}

func TestProxy_publish(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		p       *Proxy
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.publish(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Proxy.publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
