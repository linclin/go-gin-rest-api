
CREATE TABLE change_logs (
	id varchar(36) NOT NULL,
	created_at BIGINT NULL,
    action varchar(20) NOT NULL,
    object_id varchar(50) NOT NULL,
    object_type varchar(50) NOT NULL,
    raw_object TEXT NULL,
    raw_meta TEXT NULL,
    raw_diff TEXT NULL,
    created_by TEXT NOT NULL,
	CONSTRAINT change_logs_PK PRIMARY KEY (id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_general_ci;
CREATE INDEX change_logs_object_id_IDX USING BTREE ON db_template (object_id);
CREATE INDEX change_logs_object_type_IDX USING BTREE ON db_template (object_type);



CREATE TRIGGER tr_b_ins_change_logs_timestamp
BEFORE INSERT
ON change_logs FOR EACH ROW
BEGIN
    IF new.id = "00000000-0000-0000-0000-000000000000" THEN SET new.id = UUID(); END IF;
	SET new.created_at = UNIX_TIMESTAMP();
END;