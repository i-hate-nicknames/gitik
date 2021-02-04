package constants

import "errors"

type ObjectType string

const (
	TypeBlob   ObjectType = "blob"
	TypeTree   ObjectType = "tree"
	TypeCommit ObjectType = "commit"
)

func (t ObjectType) String() string {
	switch t {
	case TypeBlob:
		return "blob"
	case TypeTree:
		return "tree"
	case TypeCommit:
		return "commit"
	default:
		return "_unknown"
	}
}

func (t ObjectType) Encode() []byte {
	return []byte(t.String())
}

var ErrUnknownType = errors.New("unknown object type")

func DecodeType(encoded []byte) (ObjectType, error) {
	var otype ObjectType
	switch string(encoded) {
	case "blob":
		otype = TypeBlob
	case "tree":
		otype = TypeTree
	case "commit":
		otype = TypeCommit
	default:
		return otype, ErrUnknownType
	}
	return otype, nil
}
