package validation

import (
	"fmt"
	"strings"
)

func EscapeString(value any) string {
	switch value.(type) {
	case string:
		return "'" + strings.ReplaceAll(value.(string), "'", "''") + "'"
	default:
		return fmt.Sprintf("%v", value)
	}
}
