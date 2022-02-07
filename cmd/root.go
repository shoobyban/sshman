package cmd

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sshman",
	Short: "SSH Key manager",
	Long:  `Deploy keys to remote hosts`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("Run sshman --help for more help")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initCorba)

	var cfgFile string
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "~/.sshman", "config (.env format) file")

	cfgPath := path.Dir(cfgFile)
	cfgName := path.Base(cfgFile)
	v := viper.New()
	v.AddConfigPath(cfgPath)
	v.SetConfigName(cfgName)
	v.SetEnvPrefix("SSHMAN")
	v.SetConfigType("env")
	if err := v.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", v.ConfigFileUsed())
	}

	var config backend.Config

	v.AutomaticEnv()
	v.BindEnv("storage")
	err := v.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	if config.StorageFilePath == "" {
		home, _ := os.UserHomeDir()
		config.StorageFilePath = path.Join(home, ".ssh", ".sshman")
	}
	backend.SetConfig(&config)
}

func initCorba() {
	flag.Parse()
}
