package dimensions

import (
	"reflect"
	"sort"
	"testing"

	"github.com/gultekinmakif/llama-watch/internal/testutil"
)

func TestDetectKinds(t *testing.T) {
	root := t.TempDir()

	cases := []struct {
		name    string
		path    string
		dimType string
		want    []string
	}{
		{
			name:    "fees with both kinds",
			path:    testutil.WriteFile(t, root, "fees/aave-v2.ts", "export default { dailyFees: async () => {}, dailyRevenue: async () => {} };"),
			dimType: "fees",
			want:    []string{"dailyFees", "dailyRevenue"},
		},
		{
			name:    "fees with only dailyFees",
			path:    testutil.WriteFile(t, root, "fees/wbtc.ts", "export default { dailyFees: async () => {} };"),
			dimType: "fees",
			want:    []string{"dailyFees"},
		},
		{
			name:    "fees with neither key returns empty",
			path:    testutil.WriteFile(t, root, "fees/empty.ts", "export default { foo: 1, bar: 2 };"),
			dimType: "fees",
			want:    nil,
		},
		{
			name:    "dexs returns dailyVolume",
			path:    testutil.WriteFile(t, root, "dexs/uniswap-v3.ts", "module.exports = { dailyVolume: () => 0 };"),
			dimType: "dexs",
			want:    []string{"dailyVolume"},
		},
		{
			name:    "open-interest returns openInterestAtEnd",
			path:    testutil.WriteFile(t, root, "open-interest/gmx.ts", "export const openInterestAtEnd = () => 0;"),
			dimType: "open-interest",
			want:    []string{"openInterestAtEnd"},
		},
		{
			name:    "options returns both candidates when present",
			path:    testutil.WriteFile(t, root, "options/lyra.ts", "export default { dailyNotionalVolume: ()=>0, dailyPremiumVolume: ()=>0 };"),
			dimType: "options",
			want:    []string{"dailyNotionalVolume", "dailyPremiumVolume"},
		},
		{
			name:    "unknown dimType returns empty no error",
			path:    testutil.WriteFile(t, root, "mystery/x.ts", "export default { dailyFees: ()=>0 };"),
			dimType: "mystery",
			want:    nil,
		},
		{
			name:    "case-sensitive wrong case does not match",
			path:    testutil.WriteFile(t, root, "fees/wrongcase.ts", "export default { DailyFees: ()=>0 };"),
			dimType: "fees",
			want:    nil,
		},
		{
			name:    "word boundary suffix does not match",
			path:    testutil.WriteFile(t, root, "fees/suffix.ts", "export const dailyFeesAggregate = 0;"),
			dimType: "fees",
			want:    nil,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := DetectKinds(tc.path, tc.dimType)
			if err != nil {
				t.Fatalf("DetectKinds: %v", err)
			}
			sort.Strings(got)
			sort.Strings(tc.want)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDetectKinds_MissingFile(t *testing.T) {
	_, err := DetectKinds("/nonexistent/x.ts", "fees")
	if err == nil {
		t.Fatalf("want error on missing file, got nil")
	}
}
