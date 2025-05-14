package cmdutils

import (
	"fmt"
	"sort"
	"strings"
	"time"
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

func ValidateDate(flagName, val string) error {
	if val == "" {
		return nil
	}
	if _, err := time.Parse(time.DateOnly, val); err != nil {
		return fmt.Errorf("invalid date for --%s: %q (must be YYYY-MM-DD)", flagName, val)
	}
	return nil
}
