package cmdutils

import (
	"fmt"
	"sort"
	"strings"
)

// ValidateStringEnum returns an error if val is not one of allowed.
func ValidateStringEnum(flagName, val string, allowed []string) error {
	for _, a := range allowed {
		if val == a {
			return nil
		}
	}

	copyAllowed := append([]string(nil), allowed...)
	sort.Strings(copyAllowed)
	return fmt.Errorf("invalid value %q for flag --%s; valid values are: %s", val, flagName, strings.Join(copyAllowed, ", "))
}
