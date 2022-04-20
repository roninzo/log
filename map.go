package log

var (
	noArg = []interface{}{}
	noMap = Map{}
)

// Map is a shortcut for map[string]interface{} key, values, useful for JSON
// returns. Where the values can be any Go type, but convertable to a string.
type Map map[string]interface{}

func ParseArgs(args ...interface{}) ([]interface{}, Map) {
	n := len(args)
	switch n {
	case 0:
		return noArg, noMap
	default:
		last := n - 1
		v := args[last]
		switch m := v.(type) {
		case Map:
			if n > 1 {
				return args[:last], m
			}
			return noArg, m
		default:
			return args, noMap
		}
	}
}
