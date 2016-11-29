package appleopensource

import "testing"

func TestPlatform_String(t *testing.T) {
	tests := []struct {
		name string
		p    Platform
		want string
	}{
		{
			name: "macos",
			p:    MacOS,
			want: "macos",
		},
		{
			name: "xcode",
			p:    Xcode,
			want: "xcode",
		},
		{
			name: "ios",
			p:    IOS,
			want: "ios",
		},
		{
			name: "server",
			p:    Server,
			want: "server",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("Platform.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
