package query

import (
	"database/sql"
	"fmt"
)

type QDelete struct {
	DB *sql.DB
	param_index int
	param_value []any
	qtable      string
	qwhere      string
}

func (self *QDelete) Table(tbl_name string) *QDelete {
	self.qtable = tbl_name
	return self
}

func (self *QDelete) Where(condition string, operator string, value any) *QDelete {
	if self.qwhere == "" {
		self.qwhere = " WHERE"
	}
	self.param_index += 1
	self.param_value = append(self.param_value, value)
	self.qwhere += fmt.Sprintf(" %s %s $%d", condition, operator, self.param_index)
	return self
}

func (self *QDelete) OrWhere(condition string, operator string, value any) *QDelete {
	return self.Where(" OR "+condition, operator, value)
}

func (self *QDelete) AndWhere(condition string, operator string, value any) *QDelete {
	return self.Where(" AND "+condition, operator, value)
}

func (self *QDelete) GetQuery() string {
	query := fmt.Sprintf(" DELETE FROM %s %s",
		self.qtable,
		self.qwhere)
	return query
}

func (self *QDelete) Run() (sql.Result, error) {
	result, err := self.DB.Exec(self.GetQuery(), self.param_value...)
	return result, err
}
