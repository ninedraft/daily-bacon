package timezone

import "github.com/lithammer/fuzzysearch/fuzzy"

func Find(name string) []string {
	return fuzzy.FindNormalizedFold(name, allTimezones())
}
