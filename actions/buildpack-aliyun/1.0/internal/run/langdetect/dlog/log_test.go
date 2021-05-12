package dlog

import "testing"

func Test_handleFilePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "/.pipeline/container/context/git-checkout",
			args: args{
				path: "/.pipeline/container/context/git-checkout",
			},
			want: "${git-checkout}",
		},
		{
			name: "/.pipeline/container/context/git-checkout/pom.xml",
			args: args{
				path: "/.pipeline/container/context/git-checkout/pom.xml",
			},
			want: "${git-checkout}/pom.xml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleFilePath(tt.args.path); got != tt.want {
				t.Errorf("handleFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
