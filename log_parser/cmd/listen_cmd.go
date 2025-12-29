package cmd

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

type Listener struct {
}

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "listen current dir, parse new zip file",
	Long:  "If a new zip file appears in src dir, it will immediately perform log parsing and output the parsing result file to dest dir",
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
		srcDir, err := cmd.Flags().GetString("srcDir")
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

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		err = watcher.Add(srcDir)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Watching ", srcDir)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					if isZipFile(event.Name) {
						log.Println("Get new zip file:", event.Name)
						parser, err := NewParser(
							WithMerge(merge),
							WithUnzip(unzip),
							WithSrcZip(event.Name),
							WithDestDir(destDir),
							WithModules(modules...),
						)
						if err != nil {
							log.Fatalln(err)
						}
						err = parser.Parse()
						if err != nil {
							log.Fatalln(err)
						}
						log.Println("parsed succ")
						log.Println("Watching ", srcDir)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("watch error:", err)
			}
		}
	},
}

func isZipFile(path string) bool {
	return strings.EqualFold(filepath.Ext(path), ".zip")
}

func init() {
	// listenCmd.Flags().BoolP("unzip", "u", true, "unzip at first")
	listenCmd.Flags().Bool("merge", true, "merge same module log after parsed")
	listenCmd.Flags().StringP("srcDir", "s", "./", "listen zip file directory")
	listenCmd.Flags().StringP("destDir", "d", "./", "dest directory")
	listenCmd.Flags().StringSliceP(
		"modules",
		"m",
		[]string{"APP_PROXY", "EVENTTASK", "APP_PERSIST", "EVENTTASK_PERSIST"},
		"log modules to parse (e.g. APP_PROXY,EVENTTASK,APP_PERSIST,EVENTTASK_PERSIST)",
	)
	rootCmd.AddCommand(listenCmd)
}
