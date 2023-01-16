package server

import (
	"context"
	"errors"
	iserver "github.com/rzetelskik/iii/pkg/server"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
)

type serverOptions struct {
	addr         string
	redirectAddr string
}

func newServerOptions() *serverOptions {
	return &serverOptions{
		addr:         "localhost:8080",
		redirectAddr: "https://logowanie.uw.edu.pl/cas/login?locale=pl",
	}
}

func NewServerCommand() *cobra.Command {
	o := newServerOptions()

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run server",
		Long:  "Run server",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate()
			if err != nil {
				return err
			}

			err = o.Complete()
			if err != nil {
				return err
			}

			err = o.Run(cmd)
			if err != nil {
				return err
			}

			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().StringVar(&o.addr, "address", o.addr, "Address to listen on")
	cmd.Flags().StringVar(&o.redirectAddr, "redirect-address", o.redirectAddr, "Address to redirect to on login")

	return cmd
}

func (o *serverOptions) Validate() error {
	// TODO
	return nil
}

func (o *serverOptions) Complete() error {
	// TODO
	return nil
}

func (o *serverOptions) Run(cmd *cobra.Command) error {
	var err error

	cliflag.PrintFlags(cmd.Flags())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)
	go func() {
		s := <-c
		klog.Infof("Received first shutdown signal: %s. Shutting down gracefully.", s)
		cancel()
		<-c
		klog.Infof("Received second shutdown signal: %s. Exiting.", s)
		os.Exit(1)
	}()

	var wg sync.WaitGroup
	defer wg.Wait()

	server := iserver.NewServer(o.addr, o.redirectAddr)

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		klog.Info("Shutting down the server")
		defer klog.Info("Server shut down")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			klog.Fatalf("Couldn't terminate gracefully: %v", err)
		}
	}()

	klog.Infof("Starting web server on: %s", server.Addr)
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		klog.Fatalf("Couldn't listen on %s: %v", server.Addr, err)
	}

	return nil
}
