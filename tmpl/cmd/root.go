package cmd

import (
	"github.com/KalleDK/go-tmpl/pkg/tmpl"

	"github.com/spf13/cobra"
)

var rootFlags struct {
	Input       string
	Output      string
	Context     string
	ContextType string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tmpl",
	Short: "Use go's template engine",
	Long:  `Use go's template engine to create files`,
	RunE: func(cmd *cobra.Command, args []string) error {
		t := tmpl.Template{
			Source:      rootFlags.Input,
			Destination: rootFlags.Output,
			Context:     rootFlags.Context,
			ContextType: tmpl.ContextType(rootFlags.ContextType),
		}
		err := t.Execute()
		if err != nil {
			return err
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringVarP(&rootFlags.Input, "input", "i", "-", "Input file '-' is stdin")
	rootCmd.Flags().StringVarP(&rootFlags.Output, "output", "o", "-", "Output file '-' is stdout")
	rootCmd.Flags().StringVarP(&rootFlags.Context, "context", "c", "-", "Context file '-' is env")
	rootCmd.Flags().StringVarP(&rootFlags.ContextType, "type", "t", "env", "Environment type")
}
