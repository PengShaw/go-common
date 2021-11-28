# go-cdc

```go
package main

import (
	"fmt"
	"github.com/PengShaw/go-common/go-cdc/mysql"
	"time"
)

type A struct {
	Id        uint       `binlog:"column:id"`
	Name      string     `binlog:"column:name"`
	IsAdmin   bool       `binlog:"column:isAdmin"`
	CreatedAt *time.Time `binlog:"column:createdAt"`
}

func main() {
	var oldCh = make(chan A)
	var newCh = make(chan A)

	var printer = func(prefix string, ch chan A) {
		for {
			a := <-ch
			fmt.Printf(prefix+" %v \n", a)
		}
	}
	go printer("old", oldCh)
	go printer("new", newCh)

	tables := []mysql.Table{
		mysql.Table{
			Name:  "origin",
			Model: A{},
			HandlerFunc: func(oldItem interface{}, newItem interface{}, table string) {
				if oldItem != nil {
					oldA := oldItem.(A)
					oldCh <- oldA
				}
				if newItem != nil {
					newA := newItem.(A)
					newCh <- newA
				}
			},
		},
	}
	cdc, err := mysql.NewCDC(&mysql.Options{
		Password: "ChangeIt",
		Database: "logs",
		Tables:   tables,
	})
	if err != nil {
		println("error!!! ", err.Error())
		return
	}
	cdc.Listen()
}
```

## MySQL 
```shell
docker run -d --name mysql \
    -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD=ChangeIt \
    mysql:8
``` 

```mysql
show variables like '%log_bin%';

select * from information_schema.processlist as p where p.command = 'Binlog Dump';

CREATE TABLE origin(id INT(11), name VARCHAR(25), isAdmin boolean, createdAt datetime );
```

