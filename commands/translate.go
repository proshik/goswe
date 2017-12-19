package commands


import (
	"github.com/spf13/cobra"
	"fmt"
	"strings"
)


var translate = &cobra.Command{
	Use:   "translate",
	Short: "open translate screen",
	Run: translateRun,
}

func init() {
	//RootCmd.PersistentFlags().StringVar(&Cf, "config", "", "config file (default is $HOME/.gotrew/config.yaml)")
}

func translateRun(cmd *cobra.Command, args []string) {

	t, err := cmd.Flags().GetInt("times")
	if err != nil {
		panic(err)
	}

	fmt.Printf("times = %d\n", t)

	fmt.Println(strings.Join(args, " "))
}


