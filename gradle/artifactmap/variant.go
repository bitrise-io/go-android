package artifactmap

import (
	"path/filepath"
	"strings"
)

// ArtifactVariant identifies the build variant an artifact belongs to: the
// module together with the merged Gradle variant name (for example
// "demoRelease"). It is only compared between artifacts, so its exact textual
// form is not a stable contract; an app artifact and the mapping.txt built for
// the same variant just need to produce equal values. The module is included so
// artifacts from different modules that share a variant name (for example both
// producing "release") do not collide.
type ArtifactVariant struct {
	Module  string
	Variant string
}

// VariantFromPath reports the build variant encoded in a Gradle build-output
// path. It reconciles the directory shapes the Android Gradle Plugin uses:
// APKs split the variant into flavor and build type
// (build/outputs/apk/demo/release/), AABs use a single merged directory
// (build/outputs/bundle/demoRelease/), and mapping files use either the merged
// directory (build/outputs/mapping/demoRelease/), the merged directory with a
// deeper "minify..." task subdirectory (intermediates/mapping/demoRelease/
// minifyDemoReleaseWithR8/), or the ProGuard-era split shape
// (outputs/mapping/demo/release/). Joining the split segments yields the same
// merged name, so an app artifact and its mapping resolve to equal
// ArtifactVariants.
//
// It anchors on the "outputs" (or "intermediates") directory and reads the
// artifact kind that follows it, scanning RIGHT-TO-LEFT so the marker closest
// to the file wins: a checkout directory that happens to be called "outputs"
// higher up the path cannot hijack parsing, and a flavor directory that
// happens to be named "apk"/"bundle"/"mapping" is not mistaken for the kind
// marker. ok is false when the path is not a recognised output or mapping
// path.
func VariantFromPath(path string) (variant ArtifactVariant, ok bool) {
	segments := strings.Split(filepath.ToSlash(path), "/")

	// segments[len-1] is the file name; the marker must sit at least two
	// levels above it (kind + at least one variant directory in between).
	for i := len(segments) - 2; i >= 0; i-- {
		if segments[i] != "outputs" && segments[i] != "intermediates" {
			continue
		}

		// The directories between the kind marker and the file name encode
		// the variant. An empty list means this "marker" has no room for a
		// variant (e.g. the file itself is named "apk"): not a real marker,
		// keep scanning left.
		variantSegments := segments[i+2 : max(i+2, len(segments)-1)]
		if len(variantSegments) == 0 {
			continue
		}

		module := moduleFromSegments(segments[:i])
		switch segments[i+1] {
		case "apk", "bundle":
			return ArtifactVariant{Module: module, Variant: mergeVariantSegments(variantSegments)}, true
		case "mapping":
			// Drop a trailing shrinker-task directory (e.g.
			// "minifyDemoReleaseWithR8") so both the merged and the
			// ProGuard-era split layout reduce to the variant directories.
			if len(variantSegments) > 1 && strings.HasPrefix(variantSegments[len(variantSegments)-1], "minify") {
				variantSegments = variantSegments[:len(variantSegments)-1]
			}
			return ArtifactVariant{Module: module, Variant: mergeVariantSegments(variantSegments)}, true
		}
		// Some other "outputs" child (e.g. logs): keep scanning for the real
		// marker.
	}

	return ArtifactVariant{}, false
}

// moduleFromSegments returns the module directory (the segment right before
// "build"), or "" when the path has no "build" segment.
func moduleFromSegments(segments []string) string {
	for i := len(segments) - 1; i >= 1; i-- {
		if segments[i] == "build" {
			return segments[i-1]
		}
	}
	return ""
}

// mergeVariantSegments joins variant directory segments into the single merged
// Gradle variant name: the first segment as-is, each following one capitalised
// on its first letter, so ["demo", "release"] and ["demoRelease"] both yield
// "demoRelease".
func mergeVariantSegments(segments []string) string {
	var builder strings.Builder
	for i, segment := range segments {
		if i == 0 {
			builder.WriteString(segment)
			continue
		}
		builder.WriteString(capitalizeFirst(segment))
	}
	return builder.String()
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
