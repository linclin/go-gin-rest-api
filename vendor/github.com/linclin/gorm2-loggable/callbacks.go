package loggable

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var im = newIdentityManager()

type UpdateDiff map[string]interface{}

type DiffObject struct {
	Old interface{} `json:"old"`
	New interface{} `json:"new"`
}

// Hook for after_create.
func (p *Plugin) addCreated(db *gorm.DB) {
	if isLoggable(db.Statement.Dest) && isEnabled(db.Statement.Dest) {
		_ = p.addRecord(db, actionCreate)
	}
}

// Hook for after_update.
func (p *Plugin) addUpdated(db *gorm.DB) {
	if !isLoggable(db.Statement.Dest) || !isEnabled(db.Statement.Dest) {
		return
	}

	if p.opts.lazyUpdate {
		record, err := p.GetLastRecord(interfaceToString(db.Statement.Schema.PrimaryFields[0]), false)
		if err == nil {
			if isEqual(record.RawObject, db.Statement.Dest, p.opts.lazyUpdateFields...) {
				return
			}
		}
	}

	_ = p.addUpdateRecord(db, p.opts)
}

// Hook for after_delete.
func (p *Plugin) addDeleted(db *gorm.DB) {
	if isLoggable(db.Statement.Dest) && isEnabled(db.Statement.Dest) {
		_ = p.addRecord(db, actionDelete)
	}
}

func (p *Plugin) addUpdateRecord(db *gorm.DB, opts options) error {
	cl, err := newChangeLog(db, actionUpdate)
	if err != nil {
		return err
	}

	if opts.computeDiff {
		diff := computeUpdateDiff(db)

		if diff != nil {
			jd, err := json.Marshal(diff)
			if err != nil {
				return err
			}

			cl.RawDiff = string(jd)
		}
	}
	return db.Session(&gorm.Session{NewDB: true}).Table(p.tablename).Create(cl).Error
}

func newChangeLog(db *gorm.DB, action string) (*ChangeLog, error) {
	rawObject, err := json.Marshal(db.Statement.Dest)
	if err != nil {
		return nil, err
	}
	id := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	ui, ok := db.Get(LoggableUserTag)
	var u *User
	if !ok {
		u = &User{"null", "null", "null"}
	} else {
		u, ok = ui.(*User)
		if !ok {
			u = &User{"null", "default", "null"}
		}
	}
	us := `{"name":"null","id":"system","class":"null"}`
	ub, err := json.Marshal(u)
	if err == nil {
		us = string(ub)
	}

	scPrimaryFields := db.Statement.Schema.PrimaryFields
	objectID := ""
	for k, scPrimaryField := range scPrimaryFields {
		valuePrimaryField := interfaceToString(getPrimaryKeyValue(db, scPrimaryField.Name))
		if valuePrimaryField == "" {
			return nil, errors.New("value primary field is null")
		}
		if k == 0 {
			objectID += valuePrimaryField
		} else {
			objectID += "," + valuePrimaryField
		}
	}

	objectType := db.Statement.Table

	return &ChangeLog{
		ID:         id.String(),
		Action:     action,
		ObjectID:   objectID,
		ObjectType: objectType,
		RawObject:  string(rawObject),
		RawMeta:    string(fetchChangeLogMeta(db)),
		RawDiff:    "null",
		CreatedBy:  us,
	}, nil
}

// Writes new change log row to db.
func (p *Plugin) addRecord(db *gorm.DB, action string) error {
	cl, err := newChangeLog(db, action)
	if err != nil {
		return nil
	}

	return db.Session(&gorm.Session{NewDB: true}).Table(p.tablename).Create(cl).Error
}

func computeUpdateDiff(db *gorm.DB) UpdateDiff {
	fmt.Print("computeUpdateDiff \r\n")
	old, ok := db.Get(LoggablePrevVersion)
	fmt.Print("computeUpdateDiff LoggablePrevVersion \r\n", old, " \r\n")
	if !ok {
		return nil
	}

	ov := reflect.Indirect(reflect.ValueOf(old))
	nv := reflect.Indirect(reflect.ValueOf(db.Statement.Dest))
	names := getLoggableFieldNames(old)
	fmt.Print("computeUpdateDiff names \r\n", names, " \r\n")
	diff := make(UpdateDiff)

	for _, name := range names {

		ofv := ov.FieldByName(name).Interface()
		nfv := nv.FieldByName(name).Interface()
		fmt.Print("computeUpdateDiff names \r\n", name, ofv, nfv, reflect.DeepEqual(ofv, nfv), " \r\n")
		if !reflect.DeepEqual(ofv, nfv) {
			diff[ToSnakeCaseRegEx(name)] = DiffObject{
				Old: ofv,
				New: nfv,
			}
		}
	}

	return diff
}
