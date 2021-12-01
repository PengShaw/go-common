package mysql

import (
	"errors"
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"io"
	"os"
	"runtime/debug"
)

var DefaultWriter io.Writer = os.Stdout

type CDC struct {
	canal   *canal.Canal
	Options Options
}

type Options struct {
	Ip       string // default, "127.0.0.1"
	Port     int    // default, 3306
	User     string // default, root
	Password string
	Database string
	Tables   []Table
	Flavor   string // flavor is mysql or mariadb, default mysql
	ServerID uint32
}

type Table struct {
	Name        string
	HandlerFunc TableHandlerFunc
}

type TableHandlerFunc func(oldItem map[string]interface{}, newItem map[string]interface{}, table string)

func (o *Options) init() {
	if o.Ip == "" {
		o.Ip = "127.0.0.1"
	}
	if o.Port == 0 {
		o.Port = 3306
	}
	if o.User == "" {
		o.User = "root"
	}
	switch o.Flavor {
	case "mariadb":
		o.Flavor = "mariadb"
	default:
		o.Flavor = "mysql"
	}
	if o.ServerID == 0 {
		o.ServerID = 10001
	}
}

func NewCDC(options *Options) (*CDC, error) {
	options.init()
	cfg := canal.NewDefaultConfig()
	cfg.Addr = fmt.Sprintf("%s:%d", options.Ip, options.Port)
	cfg.User = options.User
	cfg.Password = options.Password
	cfg.Flavor = options.Flavor
	cfg.ServerID = options.ServerID
	cfg.Dump.TableDB = options.Database
	cfg.Dump.ExecutionPath = ""
	var tableNames []string
	for _, v := range options.Tables {
		tableNames = append(tableNames, v.Name)
	}
	cfg.Dump.Tables = tableNames

	c, err := canal.NewCanal(cfg)
	if err != nil {
		return nil, err
	}
	return &CDC{canal: c, Options: *options}, nil
}

func (cdc *CDC) Listen() error {
	coords, err := cdc.canal.GetMasterPos()
	if err != nil {
		return err
	}

	cdc.canal.SetEventHandler(&binlogHandler{cdc: cdc})
	return cdc.canal.RunFrom(coords)
}

type binlogHandler struct {
	cdc                     *CDC
	canal.DummyEventHandler // Dummy handler from external lib
}

func (h *binlogHandler) String() string {
	return "binlogHandler"
}

func (h *binlogHandler) OnRow(e *canal.RowsEvent) error {
	defer func() {
		if r := recover(); r != nil {
			_, _ = fmt.Fprintln(DefaultWriter, "Error: Recover from ", r, " ", string(debug.Stack()))
		}
	}()
	// check tables
	var currentTable Table
	tableName := e.Table.Schema + "." + e.Table.Name
	matched := false
	for _, table := range h.cdc.Options.Tables {
		key := h.cdc.Options.Database + "." + table.Name
		if key != tableName {
			continue
		}
		currentTable = table
		matched = true
		break
	}
	if !matched {
		_, _ = fmt.Fprintln(DefaultWriter, "Warn: Not Matched Table ", tableName)
		return nil
	}

	// handle each row item
	var current = 0
	var step = 1
	if e.Action == canal.UpdateAction {
		current = 1
		step = 2
	}
	for i := current; i < len(e.Rows); i += step {
		item, err := h.getRowItem(e, i)
		if err != nil {
			_, _ = fmt.Fprintln(DefaultWriter, "Error: get row item failed ", err.Error())
			return nil
		}
		switch e.Action {
		case canal.UpdateAction:
			oldItem, err := h.getRowItem(e, i-1)
			if err != nil {
				_, _ = fmt.Fprintln(DefaultWriter, "Error: when update action, get old row item failed ", err.Error())
				return nil
			}
			currentTable.HandlerFunc(oldItem, item, tableName)
		case canal.DeleteAction:
			currentTable.HandlerFunc(item, nil, tableName)
		case canal.InsertAction:
			currentTable.HandlerFunc(nil, item, tableName)
		}
	}
	return nil
}

func (*binlogHandler) getRowItem(e *canal.RowsEvent, rowIndex int) (res map[string]interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintln("Error: Recover from ", r, " ", string(debug.Stack()))
			err = errors.New(msg)
		}
	}()
	res = make(map[string]interface{})
	for id, column := range e.Table.Columns {
		res[column.Name] = e.Rows[rowIndex][id]
	}
	return
}
