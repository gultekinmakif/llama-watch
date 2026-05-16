// readUpstreamSHAs is intentionally not tested here: it shells out to git and
// needs real repos. Cover it via an operator-run integration check.
package main

import "testing"

func TestShasUnchanged(t *testing.T) {
	cases := []struct {
		name    string
		current map[string]string
		last    map[string]string
		want    bool
	}{
		{
			name:    "equal maps",
			current: map[string]string{"a": "1", "b": "2"},
			last:    map[string]string{"b": "2", "a": "1"},
			want:    true,
		},
		{
			name:    "different sha for same repo",
			current: map[string]string{"a": "1", "b": "2"},
			last:    map[string]string{"a": "1", "b": "X"},
			want:    false,
		},
		{
			name:    "different repo set same size",
			current: map[string]string{"a": "1", "b": "2"},
			last:    map[string]string{"a": "1", "c": "2"},
			want:    false,
		},
		{
			name:    "last empty",
			current: map[string]string{"a": "1"},
			last:    map[string]string{},
			want:    false,
		},
		{
			name:    "last nil",
			current: map[string]string{"a": "1"},
			last:    nil,
			want:    false,
		},
		{
			name:    "current empty",
			current: map[string]string{},
			last:    map[string]string{"a": "1"},
			want:    false,
		},
		{
			name:    "both nil",
			current: nil,
			last:    nil,
			want:    false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := shasUnchanged(tc.current, tc.last); got != tc.want {
				t.Fatalf("shasUnchanged(%v, %v) = %v, want %v", tc.current, tc.last, got, tc.want)
			}
		})
	}
}
