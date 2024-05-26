package query

import (
	"database/sql"
	"fmt"
	"strings"
)

type QInsert struct {
	DB *sql.DB
	param_index int
	param_value []any
	qtable      string
	qcolumn     []string
	qvalues     []string
}

func (self *QInsert) Table(tbl_name string) *QInsert {
	self.qtable = tbl_name
	return self
}

func (self *QInsert) Column(col_name ...string) *QInsert {
	self.qcolumn = append(self.qcolumn, col_name...)
	return self
}

func (self *QInsert) Values(col_value ...string) *QInsert {
	for _, value := range col_value {
		self.param_index += 1
		self.param_value = append(self.param_value, value)
		self.qvalues = append(self.qvalues, fmt.Sprintf("$%d", self.param_index))
	}
	return self
}

func (self *QInsert) GetQuery() string {
	query := fmt.Sprintf(" INSERT INTO %s (%s) VALUES (%s)",
		self.qtable,
		strings.Join(self.qcolumn, ","),
		strings.Join(self.qvalues, ","))
	return query
}

func (self *QInsert) GetQueryReturn(col_return ...string) string {
	query := fmt.Sprintf(" INSERT INTO %s (%s) VALUES (%s) RETURNING %s",
		self.qtable,
		strings.Join(self.qcolumn, ","),
		strings.Join(self.qvalues, ","),
		strings.Join(col_return, ","))
	return query
}

func (self *QInsert) RunReturn(col_return ...string) (*sql.Row) {
	row := self.DB.QueryRow(self.GetQueryReturn(col_return...), self.param_value...)
	return row
}

func (self *QInsert) Run() (sql.Result, error) {
	result, err := self.DB.Exec(self.GetQuery(), self.param_value...)
	return result, err
}
