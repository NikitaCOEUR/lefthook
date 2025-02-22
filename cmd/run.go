package cmd

import (
	"github.com/spf13/cobra"

	"github.com/evilmartians/lefthook/internal/config"
	"github.com/evilmartians/lefthook/internal/lefthook"
)

func newRunCmd(opts *lefthook.Options) *cobra.Command {
	runArgs := lefthook.RunArgs{}

	runHookCompletions := func(cmd *cobra.Command, args []string, toComplete string) (ret []string, compDir cobra.ShellCompDirective) {
		compDir = cobra.ShellCompDirectiveNoFileComp
		if len(args) != 0 {
			return
		}
		ret = lefthook.ConfigHookCompletions(opts)
		ret = append(ret, config.AvailableHooks[:]...)
		return
	}

	runHookCommandCompletions := func(cmd *cobra.Command, args []string, toComplete string) (ret []string, compDir cobra.ShellCompDirective) {
		compDir = cobra.ShellCompDirectiveNoFileComp
		if len(args) == 0 {
			return
		}
		ret = lefthook.ConfigHookCommandCompletions(opts, args[0])
		return
	}

	runCmd := cobra.Command{
		Use:               "run hook-name [git args...]",
		Short:             "Execute group of hooks",
		Example:           "lefthook run pre-commit",
		ValidArgsFunction: runHookCompletions,
		Args:              cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// args[0] - hook name
			// args[1:] - git hook arguments, number and value depends on the hook
			return lefthook.Run(opts, runArgs, args[0], args[1:])
		},
	}

	runCmd.Flags().BoolVarP(
		&runArgs.Force, "force", "f", false,
		"force execution of commands that can be skipped",
	)

	runCmd.Flags().BoolVarP(
		&runArgs.NoTTY, "no-tty", "n", false,
		"run hook non-interactively, disable spinner",
	)

	runCmd.Flags().BoolVar(
		&runArgs.AllFiles, "all-files", false,
		"run hooks on all files",
	)

	runCmd.Flags().StringSliceVar(&runArgs.Files, "files", nil, "run on specified files. takes precedence over --all-files")

	runCmd.Flags().StringSliceVar(&runArgs.RunOnlyCommands, "commands", nil, "run only specified commands")
	_ = runCmd.RegisterFlagCompletionFunc("commands", runHookCommandCompletions)

	return &runCmd
}
