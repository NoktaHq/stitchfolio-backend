package scopes

import (
	"context"
	"fmt"

	"github.com/loop-kar/pixie/constants"
	"github.com/loop-kar/pixie/log"
	"github.com/loop-kar/pixie/util"
	"gorm.io/gorm"
)

func TasksForCurrentUser() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		val, ok := db.Get(constants.USER_ID)
		if !ok || val == nil {
			return db
		}

		userID, ok := val.(*uint)
		if !ok || userID == nil || *userID == 0 {
			return db
		}
		tableName, err := getTableSQL(db)
		if err != nil {
			ctx := context.Background()
			log.Error(&ctx, err)
			return db
		}

		return db.Where(
			`(`+tableName+`.assigned_to_id = ? OR`+tableName+`.created_by_id = ?)`,
			*userID, *userID,
		)
	}
}

func GetTasks_Search(search string) func(db *gorm.DB) *gorm.DB {
	defaultReturn := func(db *gorm.DB) *gorm.DB { return db }

	if util.IsNilOrEmptyString(&search) {
		return defaultReturn
	}

	return func(db *gorm.DB) *gorm.DB {
		formattedSearch := util.EncloseWithPercentageOperator(search)
		whereClause := fmt.Sprintf(
			`("stich"."Tasks".title ILIKE %s OR "stich"."Tasks".description ILIKE %s)`,
			formattedSearch, formattedSearch,
		)
		return db.Where(whereClause)
	}
}

func GetTasks_Filter(filters string) func(db *gorm.DB) *gorm.DB {
	defaultReturn := func(db *gorm.DB) *gorm.DB { return db }

	if util.IsNilOrEmptyString(&filters) {
		return defaultReturn
	}

	fieldMap := map[string]string{
		"Status":     "status",
		"Priority":   "priority",
		"DueDate":    "due_date",
		"CompletedAt": "completed_at",
	}

	return func(db *gorm.DB) *gorm.DB {
		queryString := util.BuildQuery(filters, fieldMap)
		if !util.IsNilOrEmptyString(&queryString) {
			return db.Where(queryString)
		}
		return db
	}
}
