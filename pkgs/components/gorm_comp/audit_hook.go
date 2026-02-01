package gorm_comp

import (
	"context"
	"reflect"

	"github.com/dukk308/beetool.dev-go-starter/pkgs/constants"
	"gorm.io/gorm"
)

func registerAuditHook(db *gorm.DB) {
	db.Callback().Create().Before("gorm:before_create").Register("audit:before_create", auditBeforeCreate)
	db.Callback().Update().Before("gorm:before_update").Register("audit:before_update", auditBeforeUpdate)
	db.Callback().Delete().Before("gorm:before_delete").Register("audit:before_delete", auditBeforeDelete)
}

func getUserIDFromContext(ctx context.Context) *string {
	if ctx == nil {
		return nil
	}
	v := ctx.Value(constants.ContextKeyUserID)
	if v == nil {
		return nil
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return nil
	}
	return &s
}

func setAuditField(dest interface{}, fieldName string, value *string) {
	if dest == nil || value == nil {
		return
	}
	val := reflect.ValueOf(dest)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}
	f := val.FieldByName(fieldName)
	if !f.IsValid() || !f.CanSet() {
		return
	}
	if f.Kind() != reflect.Ptr {
		return
	}
	if f.Type().Elem() != reflect.TypeOf("") {
		return
	}
	f.Set(reflect.ValueOf(value))
}

func auditBeforeCreate(db *gorm.DB) {
	userID := getUserIDFromContext(db.Statement.Context)
	if userID == nil {
		return
	}
	dest := db.Statement.Dest
	setAuditField(dest, "CreatedBy", userID)
	setAuditField(dest, "UpdatedBy", userID)
}

func auditBeforeUpdate(db *gorm.DB) {
	userID := getUserIDFromContext(db.Statement.Context)
	if userID == nil {
		return
	}
	dest := db.Statement.Dest
	setAuditField(dest, "UpdatedBy", userID)
	if db.Statement.Schema != nil {
		if f, ok := db.Statement.Schema.FieldsByDBName["updated_by"]; ok && f != nil {
			db.Statement.SetColumn("updated_by", *userID)
		}
	}
}

func auditBeforeDelete(db *gorm.DB) {
	userID := getUserIDFromContext(db.Statement.Context)
	if userID == nil {
		return
	}
	dest := db.Statement.Dest
	setAuditField(dest, "DeletedBy", userID)
	if db.Statement.Schema != nil {
		if f, ok := db.Statement.Schema.FieldsByDBName["deleted_by"]; ok && f != nil {
			db.Statement.SetColumn("deleted_by", *userID)
		}
	}
}
