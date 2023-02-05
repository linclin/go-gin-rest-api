package loggable

import (
	"gorm.io/gorm"
)

// Plugin is a hook for gorm.
type Plugin struct {
	db        *gorm.DB
	opts      options
	tablename string
}

// Register initializes Plugin for provided gorm.DB.
// There is also available some options, that should be passed there.
// Options cannot be set after initialization.
func Register(db *gorm.DB, tablename string, opts ...Option) (Plugin, error) {
	// err := db.AutoMigrate(&ChangeLog{}).Error
	// if err != nil {
	// 	return Plugin{}, err
	// }
	o := options{}
	for _, option := range opts {
		option(&o)
	}
	tn := DefaultTableName
	if tablename != "" {
		tn = tablename
	}
	p := Plugin{db: db, opts: o, tablename: tn}
	callback := db.Callback()
	callback.Create().After("gorm:create").Register("loggable:after_create", p.addCreated)
	callback.Update().After("gorm:update").Register("loggable:after_update", p.addUpdated)
	callback.Delete().After("gorm:delete").Register("loggable:after_delete", p.addDeleted)
	return p, nil
}

// GetRecords returns all records by objectId.
// Flag prepare allows to decode content of Raw* fields to direct fields, e.g. RawObject to Object.
func (p *Plugin) GetRecords(objectId string, prepare bool) (changes []ChangeLog, err error) {
	defer func() {
		if prepare {
			for i := range changes {
				if t, ok := p.opts.metaTypes[changes[i].ObjectType]; ok {
					err = changes[i].prepareMeta(t)
					if err != nil {
						return
					}
				}
				if t, ok := p.opts.objectTypes[changes[i].ObjectType]; ok {
					err = changes[i].prepareObject(t)
					if err != nil {
						return
					}
				}
			}
		}
	}()
	return changes, p.db.Table(p.tablename).Where("object_id = ?", objectId).Find(&changes).Error
}

// GetLastRecord returns last by creation time (CreatedAt field) change log by provided object id.
// Flag prepare allows to decode content of Raw* fields to direct fields, e.g. RawObject to Object.
func (p *Plugin) GetLastRecord(objectId string, prepare bool) (change ChangeLog, err error) {
	defer func() {
		if prepare {
			if t, ok := p.opts.metaTypes[change.ObjectType]; ok {
				err := change.prepareMeta(t)
				if err != nil {
					return
				}
			}
			if t, ok := p.opts.objectTypes[change.ObjectType]; ok {
				err := change.prepareObject(t)
				if err != nil {
					return
				}
			}
		}
	}()
	return change, p.db.Table(p.tablename).Where("object_id = ?", objectId).Order("created_at DESC").Limit(1).Find(&change).Error
}
