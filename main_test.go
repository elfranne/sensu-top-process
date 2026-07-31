package main

import (
	"testing"

	"github.com/sensu/sensu-plugin-sdk/sensu"
)

func TestMain(t *testing.T) {
}

// checkArgs reads the package level plugin config, so these cases mutate it in
// place and cannot run in parallel.
func TestCheckArgs(t *testing.T) {
	cases := []struct {
		name    string
		cpu     float64
		memory  float32
		scheme  string
		sample  float64
		wantErr bool
	}{
		{"valid", 15.5, 20, "my_scheme", 1, false},
		{"cpu zero", 0, 20, "my_scheme", 1, true},
		{"cpu negative", -1, 20, "my_scheme", 1, true},
		{"cpu hundred", 100, 20, "my_scheme", 1, true},
		{"memory zero", 15.5, 0, "my_scheme", 1, true},
		{"memory negative", 15.5, -1, "my_scheme", 1, true},
		{"memory hundred", 15.5, 100, "my_scheme", 1, true},
		{"missing scheme", 15.5, 20, "", 1, true},
		{"sample zero", 15.5, 20, "my_scheme", 0, true},
		{"sample negative", 15.5, 20, "my_scheme", -1, true},
		{"sample sub second", 15.5, 20, "my_scheme", 0.5, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plugin.CPU = tc.cpu
			plugin.Memory = tc.memory
			plugin.Scheme = tc.scheme
			plugin.Sample = tc.sample

			status, err := checkArgs(nil)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got none")
				}
				if status != sensu.CheckStateWarning {
					t.Errorf("status = %d, want %d", status, sensu.CheckStateWarning)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if status != sensu.CheckStateOK {
				t.Errorf("status = %d, want %d", status, sensu.CheckStateOK)
			}
		})
	}
}

// Exercises the seed/sleep/read path end to end. Deliberately makes no
// assertion about the reported percentages: a loaded CI runner makes any
// specific value flaky.
func TestExecuteCheck(t *testing.T) {
	plugin.CPU = 15.5
	plugin.Memory = 20
	plugin.Scheme = "my_scheme"
	plugin.Expand = ""
	plugin.Sample = 0.2

	status, err := executeCheck(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != sensu.CheckStateOK {
		t.Errorf("status = %d, want %d", status, sensu.CheckStateOK)
	}
}
