package query

import (
	"database/sql"
	"fmt"
	"strings"
)

type QSelect struct {
	DB          *sql.DB
	param_index int
	param_value []any
	qtable      string
	qcolumn     []string
	qjoin       string
	qwhere      string
	qlimit      string
	qoffset     string
	qorder      string
	qgroup      string
}

func (self *QSelect) Table(tbl_name ...string) *QSelect {
	self.qtable = fmt.Sprintf(" FROM %s", strings.Join(tbl_name, ","))
	return self
}

func (self *QSelect) Column(col_name ...string) *QSelect {
	self.qcolumn = append(self.qcolumn, col_name...)
	return self
}

func (self *QSelect) Where(condition string, operator string, value any) *QSelect {
	if self.qwhere == "" {
		self.qwhere = " WHERE"
	}
	self.param_index += 1
	self.param_value = append(self.param_value, value)
	self.qwhere += fmt.Sprintf(" %s %s $%d", condition, operator, self.param_index)
	return self
}

func (self *QSelect) WhereColumn(condition string, operator string, value any) *QSelect {
	if self.qwhere == "" {
		self.qwhere = " WHERE"
	}
	self.qwhere += fmt.Sprintf(" %s %s %s", condition, operator, value)
	return self
}

func (self *QSelect) OrWhere(condition string, operator string, value any) *QSelect {
	return self.Where(" OR "+condition, operator, value)
}

func (self *QSelect) AndWhere(condition string, operator string, value any) *QSelect {
	return self.Where(" AND "+condition, operator, value)
}

func (self *QSelect) Join(tbl_name string, condition string) *QSelect {
	self.qjoin += fmt.Sprintf(" JOIN %s ON %s ", tbl_name, condition)
	return self
}

func (self *QSelect) Limit(value int) *QSelect {
	self.param_index += 1
	self.param_value = append(self.param_value, value)
	self.qlimit += fmt.Sprintf(" LIMIT $%d ", self.param_index)
	return self
}

func (self *QSelect) Offset(value int) *QSelect {
	self.param_index += 1
	self.param_value = append(self.param_value, value)
	self.qlimit += fmt.Sprintf(" OFFSET $%d ", self.param_index)
	return self
}

func (self *QSelect) GroupBy(col_name string) *QSelect {
	self.param_index += 1
	self.param_value = append(self.param_value, col_name)
	self.qlimit += fmt.Sprintf(" GROUP BY $%d", self.param_index)
	return self
}

func (self *QSelect) OrderBy(value string, order string) *QSelect {
	self.param_index += 1
	self.param_value = append(self.param_value, value)
	self.qlimit += fmt.Sprintf(" ORDER BY $%d %s", self.param_index, order)
	return self
}

func (self *QSelect) GetQuery() string {
	query := fmt.Sprintf(" SELECT %s %s %s %s %s %s %s %s",
		strings.Join(self.qcolumn, ","),
		self.qtable,
		self.qjoin,
		self.qwhere,
		self.qgroup,
		self.qorder,
		self.qlimit,
		self.qoffset)
	return query
}

func (self *QSelect) Run() (*sql.Rows, error) {
	rows, err := self.DB.Query(self.GetQuery(), self.param_value...)
	return rows, err
}

func (self *QSelect) RunRow() (*sql.Row) {
	row := self.DB.QueryRow(self.GetQuery(), self.param_value...)
	return row
}
