package storage

import "errors"

// ObjectType is a type of an object that can be stored in the storage
type ObjectType string

const (
	// TypeBlob is a user file that was added to the repository
	TypeBlob ObjectType = "blob"
	// TypeTree is a user directory
	TypeTree ObjectType = "tree"
	// TypeCommit is a commit
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

// Encode object, a counterpart of Decode function
func (t ObjectType) Encode() []byte {
	return []byte(t.String())
}

// ErrUnknownType is returned when the type of object is not one of the supported types
var ErrUnknownType = errors.New("unknown object type")

// Decode object from data, Encode coutnerpart
func Decode(encoded []byte) (ObjectType, error) {
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
