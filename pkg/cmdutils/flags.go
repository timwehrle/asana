package cmdutils

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"
)

// ValidateStringEnum returns an error if val is not one of allowed.
func ValidateStringEnum(flagName, val string, allowed []string) error {
	if slices.Contains(allowed, val) {
		return nil
	}
	copyAllowed := slices.Clone(allowed)
	sort.Strings(copyAllowed)
	return fmt.Errorf("invalid value %q for flag --%s; valid values are: %s", val, flagName, strings.Join(copyAllowed, ", "))
}

// ValidateDate returns an error if val is not a valid date in YYYY-MM-DD format.
func ValidateDate(flagName, val string) error {
	if val == "" {
		return nil
	}
	if _, err := time.Parse(time.DateOnly, val); err != nil {
		return fmt.Errorf("invalid date for --%s: %q (must be YYYY-MM-DD)", flagName, val)
	}
	return nil
}
