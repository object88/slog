package cmd

import (
	"fmt"

	"github.com/object88/slog"
	"github.com/spf13/cobra"
	utilflag "k8s.io/apiserver/pkg/util/flag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

// InitializeCommands sets up the cobra commands
func InitializeCommands() *cobra.Command {
	rootCmd := createRootCommand()

	return rootCmd
}

type rootCommand struct {
	cobra.Command

	kubeConfigFlags *genericclioptions.ConfigFlags
	factory         cmdutil.Factory
}

func createRootCommand() *cobra.Command {
	var rc *rootCommand
	rc = &rootCommand{
		Command: cobra.Command{
			Use:   "slog [SERVICE]",
			Short: "slog reports status and logs for a Kubernetes service",
			PreRunE: func(cmd *cobra.Command, args []string) error {
				return rc.preexecute(cmd, args)
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				return rc.execute(cmd, args)
			},
		},
	}

	rc.generateKubFactory()

	return &rc.Command
}

func (rc *rootCommand) generateKubFactory() {
	flags := rc.PersistentFlags()
	flags.SetNormalizeFunc(utilflag.WarnWordSepNormalizeFunc)
	flags.SetNormalizeFunc(utilflag.WordSepNormalizeFunc)

	rc.kubeConfigFlags = genericclioptions.NewConfigFlags()
	rc.kubeConfigFlags.AddFlags(flags)
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(rc.kubeConfigFlags)
	matchVersionKubeConfigFlags.AddFlags(flags)

	rc.factory = cmdutil.NewFactory(matchVersionKubeConfigFlags)
}

func (rc *rootCommand) preexecute(cmd *cobra.Command, args []string) error {
	return nil
}

func (rc *rootCommand) execute(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		cmd.HelpFunc()(cmd, args)
		return nil
	}

	messageChannel := make(chan slog.Message)
	podStatusChannel := make(chan slog.PodStatus)

	s := slog.NewSlog(rc.factory, messageChannel, podStatusChannel)

	err := s.Connect()
	if err != nil {
		return err
	}

	go func() {
		err := s.Load(*rc.kubeConfigFlags.Namespace)
		if err != nil {
			fmt.Printf("Load failed:\n\t%s\n", err.Error())
		}
	}()
	// err = s.Load()
	// if err != nil {
	// 	return err
	// }

	// This is blocking
	// err = tui.Run(messageChannel, podStatusChannel)
	// if err != nil {
	// 	return err
	// }

	// Just, wait for ever.
	done := make(chan int)
	<-done

	return nil
}
