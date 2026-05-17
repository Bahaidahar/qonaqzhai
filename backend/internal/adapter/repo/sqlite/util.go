package sqlite

import (
	"fmt"
	"strings"
)

// pg rewrites `?` positional placeholders to Postgres `$N` form.
// It is unaware of string literals — do NOT use it with hard-coded `?` inside
// quoted SQL strings (we have no such cases).
func pg(q string) string {
	var b strings.Builder
	b.Grow(len(q) + 8)
	n := 0
	for i := 0; i < len(q); i++ {
		c := q[i]
		if c == '?' {
			n++
			b.WriteString(fmt.Sprintf("$%d", n))
			continue
		}
		b.WriteByte(c)
	}
	return b.String()
}
