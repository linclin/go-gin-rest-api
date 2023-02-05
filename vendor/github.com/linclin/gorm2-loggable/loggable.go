package loggable

import (
	"encoding/json"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

// Interface is used to get metadata from your models.
type Interface interface {
	// Meta should return structure, that can be converted to json.
	Meta() interface{}
	// lock makes available only embedding structures.
	lock()
	// check if callback enabled
	isEnabled() bool
	// enable/disable loggable
	Enable(v bool)
}

// LoggableModel is a root structure, which implement Interface.
// Embed LoggableModel to your model so that Plugin starts tracking changes.
type LoggableModel struct {
	Disabled bool `gorm:"-"`
}

func (LoggableModel) Meta() interface{} { return nil }
func (LoggableModel) lock()             {}
func (l LoggableModel) isEnabled() bool { return !l.Disabled }
func (l LoggableModel) Enable(v bool)   { l.Disabled = !v }

// ChangeLog is a main entity, which used to log changes.
// Commonly, ChangeLog is stored in 'change_logs' table.
type ChangeLog struct {
	// Primary key of change logs.
	ID string `gorm:"primary_key;"`
	// Timestamp, when change log was created.
	CreatedAt int64 `gorm:"autoCreateTime"`
	// Action type.
	// On write, supports only 'create', 'update', 'delete',
	// but on read can be anything.
	Action string
	// ID of tracking object.
	// By this ID later you can find all object (database row) changes.
	ObjectID string `gorm:"index"`
	// Reflect name of tracking object.
	// It does not use package or module name, so
	// it may be not unique when use multiple types from different packages but with the same name.
	ObjectType string `gorm:"index"`
	// Raw representation of tracking object.
	// todo(@sas1024): Replace with []byte, to reduce allocations. Would be major version.
	RawObject string `gorm:"type:JSON"`
	// Raw representation of tracking object's meta.
	// todo(@sas1024): Replace with []byte, to reduce allocations. Would be major version.
	RawMeta string `gorm:"type:JSON"`
	// Raw representation of diff's.
	// todo(@sas1024): Replace with []byte, to reduce allocations. Would be major version.
	RawDiff string `gorm:"type:JSON"`
	// Free field to store something you want, e.g. who creates change log.
	// Not used field in gorm-loggable, but gorm tracks this field.
	CreatedBy string `gorm:"index"`
	// Field Object would contain prepared structure, parsed from RawObject as json.
	// Use RegObjectType to register object types.
	Object interface{} `gorm:"-"`
	// Field Meta would contain prepared structure, parsed from RawMeta as json.
	// Use RegMetaType to register object's meta types.
	Meta interface{} `gorm:"-"`
}

func (l *ChangeLog) prepareObject(objType reflect.Type) error {
	// Allocate new and try to decode change logs field RawObject to Object.
	obj := reflect.New(objType).Interface()
	err := json.Unmarshal([]byte(l.RawObject), obj)
	l.Object = obj
	return err
}

func (l *ChangeLog) prepareMeta(objType reflect.Type) error {
	// Allocate new and try to decode change logs field RawObject to Object.
	obj := reflect.New(objType).Interface()
	err := json.Unmarshal([]byte(l.RawMeta), obj)
	l.Meta = obj
	return err
}

// Diff returns parsed to map[string]interface{} diff representation from field RawDiff.
// To unmarshal diff to own structure, manually use field RawDiff.
func (l ChangeLog) Diff() (UpdateDiff, error) {
	var diff UpdateDiff
	err := json.Unmarshal([]byte(l.RawDiff), &diff)
	if err != nil {
		return nil, err
	}

	return diff, nil
}

func interfaceToString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	default:
		return fmt.Sprint(v)
	}
}

func fetchChangeLogMeta(db *gorm.DB) []byte {
	val, ok := db.Statement.Dest.(Interface)
	if !ok {
		return nil
	}
	data, err := json.Marshal(val.Meta())
	if err != nil {
		panic(err)
	}
	return data
}

func isLoggable(value interface{}) bool {
	_, ok := value.(Interface)
	return ok
}

func isEnabled(value interface{}) bool {
	v, ok := value.(Interface)
	return ok && v.isEnabled()
}

func getPrimaryKeyValue(db *gorm.DB, namePrimaryKey string) string {
	valuePrimaryKey := ""
	for _, field := range db.Statement.Schema.Fields {
		if namePrimaryKey == field.Name {
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
					// Get value from field
					if fieldValue, isZero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue.Index(i)); !isZero {
						valuePrimaryKey = interfaceToString(fieldValue)
					}
				}
			case reflect.Struct:
				// Get value from field
				fieldValue, isZero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue)
				if !isZero {
					valuePrimaryKey = interfaceToString(fieldValue)
				}
			}
		}
	}
	return valuePrimaryKey
}
