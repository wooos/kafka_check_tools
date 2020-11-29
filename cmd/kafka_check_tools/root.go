package main

import (
	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "kafka_check_tools",
		Short: "kafka check tools",
		SilenceUsage: true,
	}

	cmd.AddCommand(newProducerCommand(), newConsumerCommand(), newCompletionCommand())

	return cmd
}
