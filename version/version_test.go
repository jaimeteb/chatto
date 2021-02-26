package version

import "testing"

func TestBuildStr(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "build string should be formatted correctly",
			want: "version: v0.0.0\ncommit: 0000000000000000000000000000000000000000\nbuilt at: 0001-01-01 00:00:00 +0000 UTC\nbuilt by: dev\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildStr(); got != tt.want {
				t.Errorf("BuildStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
