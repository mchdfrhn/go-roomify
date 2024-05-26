package query

import (
	"database/sql"
	"fmt"
	"strings"
)

type QUpdate struct {
	DB *sql.DB
	param_index int
	param_value []any
	qtable      string
	qset        []string
	qwhere      string
}

func (self *QUpdate) Table(tbl_name string) *QUpdate {
	self.qtable = tbl_name
	return self
}

func (self *QUpdate) Set(col_name string, value string) *QUpdate {
	self.param_index += 1
	self.param_value = append(self.param_value, value)
	self.qset = append(self.qset, fmt.Sprintf("%s = $%d", col_name, self.param_index))
	return self
}

func (self *QUpdate) Where(condition string, operator string, value any) *QUpdate {
	if self.qwhere == "" {
		self.qwhere = " WHERE"
	}
	self.param_index += 1
	self.param_value = append(self.param_value, value)
	self.qwhere += fmt.Sprintf(" %s %s $%d", condition, operator, self.param_index)
	return self
}

func (self *QUpdate) OrWhere(condition string, operator string, value any) *QUpdate {
	return self.Where(" OR "+condition, operator, value)
}

func (self *QUpdate) AndWhere(condition string, operator string, value any) *QUpdate {
	return self.Where(" AND "+condition, operator, value)
}

func (self *QUpdate) GetQuery() string {
	query := fmt.Sprintf(" UPDATE %s SET %s %s",
		self.qtable,
		strings.Join(self.qset, ","),
		self.qwhere)
	return query
}

func (self *QUpdate) Run() (sql.Result, error) {
	result, err := self.DB.Exec(self.GetQuery(), self.param_value...)
	return result, err
}
