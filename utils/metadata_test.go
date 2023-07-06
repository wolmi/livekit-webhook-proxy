package utils

import "testing"

func TestMetaData_Init(t *testing.T) {
	tests := []struct {
		name string
		m    *MetaData
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Init()
		})
	}
}
