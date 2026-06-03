package store

import (
	"testing"
)

// TestMemoryStore defines table-driven tests for the in-memory store.
// Table-driven tests let you define multiple test cases in a compact, declarative form.
func TestMemoryStore(t *testing.T) {
	// tests is a slice of anonymous structs — this is the test table.
	// Each struct describes a single test case with inputs and expected outputs.
	tests := []struct {
		name       string // test case name shown by t.Run
		code       string // the short code we will Save
		original   string // the original URL to save
		lookupCode string // the code we will use for Get
		wantErr    bool   // whether we expect Get to return an error
	}{
		{
			name:       "happy path - save and retrieve",
			code:       "abc123",
			original:   "https://example.com",
			lookupCode: "abc123",
			wantErr:    false,
		},
		{
			name:       "miss - code never saved",
			code:       "", // no save for this test (we won't call Save with this code)
			original:   "",
			lookupCode: "missing",
			wantErr:    true,
		},
		{
			name:       "overwrite - save same code twice with different URLs",
			code:       "dup",
			original:   "https://second.example.com", // expected after overwrite
			lookupCode: "dup",
			wantErr:    false,
		},
	}

	// Loop over test cases and run each as a subtest with t.Run.
	for _, tc := range tests {
		// capture range variable for the closure to avoid the common gotcha
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			// subtests may run in parallel if desired; not using t.Parallel() here.
			// Create a fresh store for each subtest to isolate state.
			store := NewMemoryStore()

			// For the overwrite test we need to save twice with different URLs.
			if tc.name == "overwrite - save same code twice with different URLs" {
				// First save some initial URL.
				if err := store.Save(tc.code, "https://first.example.com"); err != nil {
					t.Fatalf("Save (first) returned unexpected error: %v", err)
				}
				// Second save should overwrite the first.
				if err := store.Save(tc.code, tc.original); err != nil {
					t.Fatalf("Save (second) returned unexpected error: %v", err)
				}
			} else if tc.name == "miss - code never saved" {
				// Intentionally do not call Save to simulate a missing code.
			} else {
				// For normal cases call Save once.
				if err := store.Save(tc.code, tc.original); err != nil {
					t.Fatalf("Save returned unexpected error: %v", err)
				}
			}

			// Call Get with the lookupCode.
			got, err := store.Get(tc.lookupCode)

			// Error checks based on tc.wantErr:
			if tc.wantErr {
				// We expected an error. If err is nil, that's a test failure.
				if err == nil {
					t.Errorf("Get(%q) expected error, got nil, url=%q", tc.lookupCode, got)
				}
				// If an error was returned, test passes for this expectation; nothing more to check.
				return
			}

			// We did not expect an error.
			if err != nil {
				t.Errorf("Get(%q) unexpected error: %v", tc.lookupCode, err)
				return
			}

			// If no error, check that returned URL matches expected URL.
			if got != tc.original {
				t.Errorf("Get(%q) = %q, want %q", tc.lookupCode, got, tc.original)
			}
		})
	}
}
