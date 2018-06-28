package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "prairiedog",
	Short: "prairiedog creates pangenome graphs",
	Long: `A pangenome graph generator with storage in Dgraph
				and Bagder. Implements a cross between a De Bruijn
				Graph and a Li-Stephen model. Source: github.com/superphy/prairiedog.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
