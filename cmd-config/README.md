# cmdconfig

```go
package main

import (
	cmdconfig "github.com/PengShaw/go-common/cmd-config"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "root cmd",
}

var subCmd = &cobra.Command{
	Use:   "sub",
	Short: "sub cmd",
	Run: func(cmd *cobra.Command, args []string) {
		println("run sub cmd")
	},
}

type config struct {
	host string
}

func main() {
	cmdconfig.SetConfigFlag(rootCmd)
	rootCmd.AddCommand(subCmd)

	var c config
	if err := cmdconfig.GetConfig(&c); err != nil {
		os.Exit(1)
	}
}
```