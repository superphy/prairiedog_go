package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/superphy/prairiedog/pangenome"
)

var (
	// Used for flags.
	cfgFile, projectBase, userLicense string

	rootCmd = &cobra.Command{
		Use:   "prairiedog",
		Short: "prairiedog creates pangenome graphs",
		Long: `A pangenome graph generator with storage in Dgraph
					and Bagder. Implements a cross between a De Bruijn
					Graph and a Li-Stephen model. Source: github.com/superphy/prairiedog.`,
		Run: func(cmd *cobra.Command, args []string) {
			pangenome.Run()
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "github.com/superphy/")
	rootCmd.PersistentFlags().StringP("author", "a", "NML", "National Microbiology Lab")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Apache 2.0")
	rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "NML chad.laing@canada.ca")
	viper.SetDefault("license", "Apache 2.0")

	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of prairiedog",
	Long:  `All software has versions. This is prairiedogs's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("prairiedog v0.0.1 -- HEAD")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
