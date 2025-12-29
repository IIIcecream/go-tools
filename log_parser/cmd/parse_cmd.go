package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "parse log files once",
	Run: func(cmd *cobra.Command, args []string) {
		unzip := true
		// unzip, err := cmd.Flags().GetBool("unzip")
		// if err != nil {
		// 	log.Fatalln(err)
		// 	return
		// }
		merge, err := cmd.Flags().GetBool("merge")
		if err != nil {
			log.Fatalln(err)
			return
		}
		srcZip, err := cmd.Flags().GetString("srcZip")
		if err != nil {
			log.Fatalln(err)
			return
		}
		destDir, err := cmd.Flags().GetString("destDir")
		if err != nil {
			log.Fatalln(err)
			return
		}
		modules, err := cmd.Flags().GetStringSlice("modules")
		if err != nil {
			log.Fatalln(err)
			return
		}
		parser, err := NewParser(
			WithMerge(merge),
			WithUnzip(unzip),
			WithSrcZip(srcZip),
			WithDestDir(destDir),
			WithModules(modules...),
		)
		if err != nil {
			log.Fatalln(err)
			return
		}
		err = parser.Parse()
		if err != nil {
			log.Fatalln(err)
			return
		}
	},
}

func init() {
	// parseCmd.Flags().BoolP("unzip", "u", true, "unzip at first")
	parseCmd.Flags().Bool("merge", true, "merge same module log after parsed")
	parseCmd.Flags().StringP("srcZip", "s", "./", "src zip file path")
	parseCmd.Flags().StringP("destDir", "d", "./", "dest path")
	parseCmd.Flags().StringSliceP(
		"modules",
		"m",
		[]string{"APP_PROXY", "EVENTTASK", "APP_PERSIST", "EVENTTASK_PERSIST"},
		"log modules to parse (e.g. APP_PROXY,EVENTTASK,APP_PERSIST,EVENTTASK_PERSIST)",
	)
	rootCmd.AddCommand(parseCmd)
}
