package loggable

const loggableTag = "gorm-loggable"
const LoggableUserTag = loggableTag + ":user"
const LoggablePrevVersion = loggableTag + ":prev_version"

const (
	actionCreate = "create"
	actionUpdate = "update"
	actionDelete = "delete"
)

const DefaultTableName = "change_logs"
