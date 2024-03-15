package internal

import (
	"github.com/null93/mirdir/pkg/template"
	. "github.com/null93/mirdir/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	Version     = "0.0.0"
	SkipPrompt  = false
	DryRun      = false
	DeleteEmpty = false
	Preserve    = false
	Verbose     = false
)

var RootCmd = &cobra.Command{
	Use:     "mirdir TPL_DIR DST_DIR",
	Version: Version,
	Short:   "CLI tool that mirrors and templates a directory structure",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tmplDir := args[0]
		destDir := args[1]
		envVars := GetEnvironmentalVars()

		if !Exists(tmplDir) {
			ExitWithError(1, "TPL_DIR does not exist", nil, false)
		}

		if !IsDirectory(tmplDir) {
			ExitWithError(2, "TPL_DIR must be a directory", nil, false)
		}

		template, templateErr := template.NewTemplate(tmplDir)
		if templateErr != nil {
			ExitWithError(3, "Failed parsing template", templateErr, Verbose)
		}

		output, renderErr := template.Render(destDir, Preserve, envVars)
		if renderErr != nil {
			ExitWithError(4, "Failed rendering template", renderErr, Verbose)
		}

		for _, path := range output {
			isEmptyFile := path.IsEmptyFile()
			doesExist := Exists(path.Path)
			if DryRun {
				path.Print(!isEmptyFile || !DeleteEmpty, Verbose)
				continue
			}
			if isEmptyFile && doesExist && DeleteEmpty {
				if !SkipPrompt && !PromptDelete(path.Path) {
					continue
				}
				if deleteErr := path.Delete(); deleteErr != nil {
					ExitWithError(5, "Failed deleting file", deleteErr, Verbose)
				}
				continue
			} else if isEmptyFile && DeleteEmpty {
				continue
			} else if doesExist && !SkipPrompt && !PromptOverwrite(path.Path, path.IsDir()) {
				continue
			}
			if writeErr := path.Write(Preserve); writeErr != nil {
				ExitWithError(6, "Failed writing file", writeErr, Verbose)
			}
		}
	},
}

func init() {
	RootCmd.Flags().SortFlags = true
	RootCmd.Flags().BoolVarP(&SkipPrompt, "yes", "y", SkipPrompt, "don't ask, just do")
	RootCmd.Flags().BoolVarP(&DryRun, "dry-run", "d", DryRun, "print output without writing")
	RootCmd.Flags().BoolVarP(&DeleteEmpty, "delete-empty", "e", DeleteEmpty, "delete file if its only whitespace")
	RootCmd.Flags().BoolVarP(&Preserve, "preserve", "p", Preserve, "preserve permissions and ownership")
	RootCmd.Flags().BoolVarP(&Verbose, "verbose", "v", Verbose, "print verbose output")
}
