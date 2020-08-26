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
			name: "stable Gradle version",
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
		{
			name: "RC Gradle version",
			out: `
Welcome to Gradle 4.10-rc-1!

Here are the highlights of this release:
 - Incremental Java compilation by default
 - Periodic Gradle caches cleanup
 - Gradle Kotlin DSL 1.0-RC1
 - Nested included builds
 - SNAPSHOT plugin versions in the ` + "`plugins {}`" + ` block

For more details see https://docs.gradle.org/4.10-rc-1/release-notes.html


------------------------------------------------------------
Gradle 4.10-rc-1
------------------------------------------------------------

Build time:   2018-08-09 06:19:37 UTC
Revision:     97951b7f541f1da43de291246cc7b17507e75a14

Kotlin DSL:   1.0-rc-1
Kotlin:       1.2.60
Groovy:       2.4.15
Ant:          Apache Ant(TM) version 1.9.11 compiled on March 23 2018
JVM:          1.8.0_241 (Oracle Corporation 25.241-b07)
OS:           Mac OS X 10.15.5 x86_64`,
			want:    "4.10-rc-1",
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
