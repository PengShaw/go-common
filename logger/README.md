# logger

```golang
package main

import (
	"github.com/PengShaw/go-common/logger"
)

func main(){
	logger.InitLogger()
	//logger.InitLoggerByOptions(&logger.Options{Filename: "./test.log"})
	//log :=  logger.GetLogger()
	//log :=  logger.GetLoggerByOptions(&logger.Options{Filename: "./test.log"})

	logger.Info("info ok")
	logger.Infof("info $s","ok")
	//log.Info("info ok")
	//log.Infof("info $s","ok")
}
```