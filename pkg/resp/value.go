package resp

type Type byte

const (
	TypeNull Type = iota
	TypeSimpleString
	TypeSimpleError
	TypeInteger
	TypeBulkString
	TypeArray
)

type Value struct {
	typ Type

	integer int
	string  []byte
	array   []Value
}

// func UnmarshalValue(r io.Reader) (Value, int, error) {
// 	typ := make([]byte, 1)
// 	if _, err := r.Read(typ); err != nil {
// 		return Value{}, 0, err
// 	}
//
// 	switch typ[0] {
// 	case '_':
// 		val := Value{typ: TypeNull}
// 		if _, err := r.Read(make([]byte, 2)); err != nil {
// 			return Value{}, 1, err
// 		}
//
// 		return val, 3, nil
// 	case '+':
// 		val := Value{typ: TypeSimpleString}
// 	case '-':
// 		val := Value{typ: TypeSimpleError}
// 	case ':':
// 		val := Value{typ: TypeInteger}
// 	case '$':
// 		val := Value{typ: TypeBulkString}
// 	case '*':
// 		val := Value{typ: TypeArray}
// 	default:
// 		return
// 	}
// }
