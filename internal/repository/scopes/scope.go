package scopes

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/loop-kar/pixie/constants"
	"github.com/loop-kar/pixie/log"
	"github.com/loop-kar/pixie/util"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

const CHANNEL_ID = "channel_id"
const And = " AND "
const OR = " OR "

func IsActive(params ...string) func(db *gorm.DB) *gorm.DB {

	if len(params) == 0 {
		return func(db *gorm.DB) *gorm.DB {

			if db.Statement.Model == nil {
				return db.Where("is_active")
			}

			tableName, err := getTableName(db)
			if err != nil {
				ctx := context.Background()
				log.Error(&ctx, err)
				return db
			}

			return db.Where(fmt.Sprintf("%s.is_active", tableName), true)
		}
	}

	//chain the where conditions and return

	var IsActive string = "is_active"
	var And string = " AND "

	return func(db *gorm.DB) *gorm.DB {

		whereClause := ""
		funk.ForEach(params, func(param string) {
			if !strings.HasPrefix(param, "E") {
				param = util.EncloseWithSymbol(param, "\"")
			}

			whereClause = whereClause + fmt.Sprintf("%s.%s = true", param, IsActive) + And
		})

		//removing the last 'and' word
		whereClause = strings.TrimSuffix(whereClause, And)

		return db.Where(whereClause)
	}

}

func Channel(params ...string) func(db *gorm.DB) *gorm.DB {

	if len(params) == 0 {
		return func(db *gorm.DB) *gorm.DB {

			var channelId uint
			if id, ok := db.Get(constants.CHANNEL_ID); ok {
				channelId = id.(uint)
			}

			//System Admin needs access to all Data
			if channelId == 0 {
				return db
			}

			if db.Statement.Model == nil {
				return db.Where("channel_id = ?", channelId)
			}

			stmt := &gorm.Statement{DB: db}
			if err := stmt.Parse(db.Statement.Model); err != nil {
				return db
			}
			tableName := stmt.Schema.Table

			return db.Where(fmt.Sprintf("%s.channel_id", tableName), channelId)

		}
	}

	//chain the where conditions and return

	return func(db *gorm.DB) *gorm.DB {

		var channelId uint
		if id, ok := db.Get(constants.CHANNEL_ID); ok {
			channelId = id.(uint)
		}

		//System Admin needs access to all Data
		if channelId == 0 {
			return db
		}

		whereClause := ""
		funk.ForEach(params, func(param string) {
			if !strings.HasPrefix(param, "E") {
				param = util.EncloseWithSymbol(param, "\"")
			}

			whereClause = whereClause + fmt.Sprintf("%s.%s = %d", param, CHANNEL_ID, channelId) + And
		})

		//removing the last 'and' word
		whereClause = strings.Trim(whereClause, And)

		return db.Where(whereClause)
	}

}

func AccessibleChannels(accesibleChannelIds []uint) func(db *gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {

		var channelId uint
		if id, ok := db.Get(constants.CHANNEL_ID); ok {
			channelId = id.(uint)
		}

		//System Admin needs access to all Data
		if channelId == 0 {
			return db
		}

		return db.Where("channel_id", accesibleChannelIds)

	}

}

func ILike(query string, params ...string) func(db *gorm.DB) *gorm.DB {

	if len(params) == 0 || util.IsNilOrEmptyString(&query) {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	return func(db *gorm.DB) *gorm.DB {

		whereClause := ""
		query = util.EncloseWithPercentageOperator(query)
		funk.ForEach(params, func(param string) {

			whereClause = whereClause + fmt.Sprintf(`%s ILIKE %s`, param, query) + OR
		})

		//removing the last 'and' word
		whereClause = strings.Trim(whereClause, OR)

		return db.Where(whereClause)
	}

}

func SelectFields(params ...string) func(db *gorm.DB) *gorm.DB {

	if len(params) == 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Select("id", params)
	}

}

// WithAuditInfo adds the created_by and updated_by fields to the select clause by performing left joins with the Users table.
// It takes the table name as a parameter to construct the join conditions and select clause.
// We can also pass the alias of the table if the query is using any alias, for example "E" in case of orders table. In that case the created_by and updated_by fields will be added as "o.created_by" and "o.updated_by"
func WithAuditInfo() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		// tableName, err := getTableSQL(db)
		// if err != nil {
		// 	ctx := context.Background()
		// 	log.Error(&ctx, err)
		// 	return db
		// }

		// return db.
		// 	Joins(`LEFT JOIN "stich"."Users" cu ON cu.id = ` + tableName + `.created_by_id`).
		// 	Joins(`LEFT JOIN "stich"."Users" uu ON uu.id = ` + tableName + `.updated_by_id`).
		// 	Select(`
		// 		` + tableName + `.*,
		// 		COALESCE(cu.first_name || ' ' || cu.last_name, '') AS created_by,
		// 		COALESCE(uu.first_name || ' ' || uu.last_name, '') AS updated_by
		// 	`)
		return db
	}
}

func getTableName(db *gorm.DB) (string, error) {
	stmt := &gorm.Statement{DB: db}
	if db.Statement.Model == nil {
		return "", errors.New("Missing Model in Statement")
	}
	if err := stmt.Parse(db.Statement.Model); err != nil {
		return "", err
	}
	return stmt.Schema.Table, nil
}

func getTableSQL(db *gorm.DB) (string, error) {
	stmt := &gorm.Statement{DB: db}
	if db.Statement.Model == nil {
		return "", errors.New("Missing Model in Statement")
	}
	if err := stmt.Parse(db.Statement.Model); err != nil {
		return "", err
	}
	return stmt.TableExpr.SQL, nil
}
