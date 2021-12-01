# go-cdc

```go
package main

import (
	"fmt"
	"github.com/PengShaw/go-common/go-cdc/mysql"
)

func main() {
	var oldCh = make(chan map[string]interface{})
	var newCh = make(chan map[string]interface{})

	var printer = func(prefix string, ch chan map[string]interface{}) {
		for {
			a := <-ch
			fmt.Printf(prefix+" %v \n", a)
		}
	}
	go printer("old", oldCh)
	go printer("new", newCh)

	tables := []mysql.Table{
		mysql.Table{
			Name: "origin",
			HandlerFunc: func(oldItem map[string]interface{}, newItem map[string]interface{}, table string) {
				oldCh <- oldItem
				newCh <- newItem
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

select *
from information_schema.processlist as p
where p.command = 'Binlog Dump';

CREATE TABLE origin
(
    id        INT(11),
    name      VARCHAR(25),
    isAdmin   boolean,
    createdAt datetime
);
```

