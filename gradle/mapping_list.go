package gradle

import "strings"

// mappingListSeparator joins the mapping-file list. The producing steps always
// emit this separator; DecodeMappingList also tolerates newline separators for
// hand-authored input.
const mappingListSeparator = "|"

// EncodeMappingList joins mapping-file paths into the positional list format
// that Android build steps export as BITRISE_MAPPING_PATH_LIST. Entry N is the
// mapping file for app artifact N in the matching app list; empty entries are
// kept (["a", "", "c"] -> "a||c") so a variant with no mapping does not shift
// the alignment of the entries that follow it.
func EncodeMappingList(paths []string) string {
	return strings.Join(paths, mappingListSeparator)
}

// DecodeMappingList parses the mapping-file list format back into a positional
// slice, PRESERVING empty entries so index alignment with the app list is kept.
// It tolerates the pipe separator plus newline and literal `\n` separators (the
// value may be set by hand in a step input), trims whitespace around the whole
// value and around each entry, and returns nil for an empty value.
func DecodeMappingList(list string) []string {
	list = strings.TrimSpace(list)
	if list == "" {
		return nil
	}

	fields := []string{list}
	for _, separator := range []string{"\n", `\n`, mappingListSeparator} {
		fields = splitEachOn(fields, separator)
	}

	paths := make([]string, len(fields))
	for i, field := range fields {
		paths[i] = strings.TrimSpace(field)
	}
	return paths
}

func splitEachOn(in []string, separator string) []string {
	var out []string
	for _, element := range in {
		out = append(out, strings.Split(element, separator)...)
	}
	return out
}
