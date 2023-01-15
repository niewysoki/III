package main

import (
	"flag"
	"fmt"
	"github.com/rzetelskik/iii/pkg/cmd/server"
	"github.com/rzetelskik/iii/pkg/sender"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"os"
)

func main() {
	var err error

	klog.InitFlags(flag.CommandLine)
	if err = flag.Set("logtostderr", "false"); err != nil {
		panic(err)
	}
	if err = flag.Set("alsologtostderr", "false"); err != nil {
		panic(err)
	}
	flag.Parse()

	klog.SetOutput(os.Stdout)

	//cmdutil.InstallKlog(command)

	cmd := &cobra.Command{}

	cmd.AddCommand(server.NewServerCommand())
	cmd.AddCommand(sender.NewSenderCommand())

	InstallKlog(cmd)

	err = cmd.Execute()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func InstallKlog(cmd *cobra.Command) {
	level := flag.CommandLine.Lookup("v").Value.(*klog.Level)
	levelPtr := (*int32)(level)
	cmd.PersistentFlags().Int32Var(levelPtr, "loglevel", *levelPtr, "Set the level of log output (0-10).")
	if cmd.PersistentFlags().Lookup("v") == nil {
		cmd.PersistentFlags().Int32Var(levelPtr, "v", *levelPtr, "Set the level of log output (0-10).")
	}
	cmd.PersistentFlags().Lookup("v").Hidden = true

	// Enable directory prefix.
	err := flag.CommandLine.Lookup("add_dir_header").Value.Set("true")
	if err != nil {
		panic(err)
	}
}
