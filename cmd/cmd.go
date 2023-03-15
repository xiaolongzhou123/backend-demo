package cmd

import (
	"fmt"

	"sso/cmd/http"
	"sso/pkg"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "sso-cli",
		Short: "Single Sign-On采用了OAuth和OpenID Connect协议",
		Long: `SSO是单点登录。它是一种身份验证机制，允许用户使用一个集中的身份验证系统来访问多个应用程序，而不需要为每个应用程序单独进行身份验证。在SSO中，用户只需登录一次，就可以在多个应用程序中使用相同的凭据来访问它们。这种机制可以提高用户体验、简化身份验证过程并提高安全性。本文中的SSO实现包括基于OAuth和OpenID Connect等协议的解决方案。
`,
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of sso",
		Long:  `All software has versions. This is sso's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("sso Static Site Generator v0.9 -- HEAD")
		},
	}

	runCmd = &cobra.Command{
		Use:   "start",
		Short: "this is how to start sso http server",
		Long:  `this is api server`,
		Run: func(cmd *cobra.Command, args []string) {
			http.Start()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $PWD/sso.yaml)")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(runCmd)
	// rootCmd.AddCommand(initCmd)
}

func initConfig() {
	c := pkg.Conf()
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		// home, err := os.UserHomeDir()
		// cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.SetConfigType("yaml")
		viper.SetConfigName("sso")
		viper.AddConfigPath(".")
		viper.ReadInConfig()

		if err := viper.Unmarshal(c); err != nil {
		}
	}

	fmt.Println("==initConfig:config:=", c.String())
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
