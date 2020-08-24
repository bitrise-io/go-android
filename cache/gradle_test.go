package cache

import "testing"

func Test_parseGradlewVersion(t *testing.T) {
	tests := []struct {
		name    string
		out     string
		want    string
		wantErr bool
	}{
		{
			name: "",
			out: `
------------------------------------------------------------
Gradle 6.1.1
------------------------------------------------------------

Build time:   2020-01-24 22:30:24 UTC
Revision:     a8c3750babb99d1894378073499d6716a1a1fa5d

Kotlin:       1.3.61
Groovy:       2.5.8
Ant:          Apache Ant(TM) version 1.10.7 compiled on September 1 2019
JVM:          1.8.0_241 (Oracle Corporation 25.241-b07)
OS:           Mac OS X 10.15.5 x86_64`,
			want:    "6.1.1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseGradleVersion(tt.out)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGradlewVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseGradlewVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
