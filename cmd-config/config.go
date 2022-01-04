package cmdconfig

import (
	"flag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func SetConfigFlag() {
	flag.StringVar(&cfgFile, "config", "./config.yaml", "config file")
	flag.Parse()
}

func SetConfigFlagByCobra(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "./config.yaml", "config file")
}

// GetConfig to struct, configs should be the point of config struct
func GetConfig(configs interface{}) error {
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(configs); err != nil {
		return err
	}
	return nil
}
