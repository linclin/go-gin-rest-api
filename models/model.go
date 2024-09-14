package models

import (
	"database/sql/driver"
	"fmt"
	"go-gin-rest-api/pkg/global"
	"time"
)

// 由于gorm提供的base model没有json tag, 使用自定义
type Model struct {
	Id        uint       `gorm:"column:Id;primary_key;comment:'自增编号'" json:"Id" rql:"filter,sort,column=Id" `
	CreatedAt time.Time  `gorm:"column:CreatedAt;comment:'创建时间'" json:"CreatedAt" rql:"filter,sort,column=CreatedAt,layout=2006-01-02 15:04:05"`
	UpdatedAt time.Time  `gorm:"column:UpdatedAt;comment:'更新时间'" json:"UpdatedAt" rql:"filter,sort,column=UpdatedAt,layout=2006-01-02 15:04:05"`
	DeletedAt *time.Time `gorm:"column:DeletedAt;comment:'删除时间(软删除)'" sql:"index" json:"DeletedAt" rql:"filter,sort,column=DeletedAt,layout=2006-01-02 15:04:05"`
}

// 表名设置
func (Model) TableName(name string) string {
	// 添加表前缀
	if global.Conf.Mysql != nil {
		return fmt.Sprintf("%s%s", global.Conf.Mysql.TablePrefix, name)
	}
	if global.Conf.Pgsql != nil {
		return fmt.Sprintf("%s%s", global.Conf.Pgsql.TablePrefix, name)
	}
	return name
}

// 自定义时间json转换
const TimeFormat = "2006-01-02 15:04:05"

type LocalTime struct {
	time.Time
}

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	// 空值不进行解析
	if len(data) == 2 {
		*t = LocalTime{Time: time.Time{}}
		return
	}
	// 指定解析的格式
	now, err := time.Parse(`"`+TimeFormat+`"`, string(data))
	*t = LocalTime{Time: now}
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	output := fmt.Sprintf("\"%s\"", t.Format(TimeFormat))
	return []byte(output), nil
}

// gorm 写入 mysql 时调用
func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// gorm 检出 mysql 时调用
func (t *LocalTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// 用于 fmt.Println 和后续验证场景
func (t LocalTime) String() string {
	return t.Format(TimeFormat)
}
