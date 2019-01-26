package cockroach

import (
	"strconv"
	"strings"
)

// setValuesList creates a list of column_name=$1 (the 1 replaced by the corresponding index) pairs and a flattened list
// of arguments; the arguments slice has capacity one more than the length of vals to allow for appending additional
// arguments without an allocation.
func setValuesList(vals map[string]interface{}) (string, []interface{}) {
	var b strings.Builder
	args := make([]interface{}, 0, len(vals)+1)
	first := true
	index := 1
	for k := range vals {
		if !first {
			b.WriteByte(',')
		}
		b.WriteString(strconv.Quote(k))
		b.WriteString("=$")
		b.WriteString(strconv.Itoa(index))
		first = false
		args = append(args, vals[k])
		index++
	}
	return b.String(), args
}
