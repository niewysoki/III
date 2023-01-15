package sender

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"

	senderassets "github.com/rzetelskik/iii/assets/sender"
	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
)

type senderOptions struct {
	hostAddr   string
	hostPort   string
	password   string
	username   string
	from       string
	to         string
	serverAddr string

	auth smtp.Auth
	msg  []byte
}

func newSenderOptions() *senderOptions {
	return &senderOptions{
		serverAddr: "localhost:8080",
	}
}

func NewSenderCommand() *cobra.Command {
	o := newSenderOptions()

	cmd := &cobra.Command{
		Use:   "sender",
		Short: "Run sender",
		Long:  "Run sender",
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

	cmd.Flags().StringVar(&o.hostAddr, "host", "", "SMTP host address")
	cmd.Flags().StringVar(&o.hostPort, "port", "", "SMTP host port")
	cmd.Flags().StringVar(&o.username, "username", "", "SMTP host username")
	cmd.Flags().StringVar(&o.password, "password", "", "SMTP host password")
	cmd.Flags().StringVar(&o.from, "from", "", "Sender address")
	cmd.Flags().StringVar(&o.to, "to", "", "Recipient address")
	cmd.Flags().StringVar(&o.serverAddr, "server", o.serverAddr, "Server address")

	cmd.MarkFlagRequired("host")
	cmd.MarkFlagRequired("port")
	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("from")
	cmd.MarkFlagRequired("to")

	return cmd
}

func (o *senderOptions) Validate() error {
	// TODO
	return nil
}

func (o *senderOptions) Complete() error {
	o.auth = smtp.PlainAuth("", o.username, o.password, o.hostAddr)

	t, err := template.New("msg").Parse(senderassets.Message)
	if err != nil {
		return fmt.Errorf("can't pasrse template: %w", err)
	}

	var buf bytes.Buffer
	inputs := map[string]any{
		"to":            o.to,
		"from":          o.from,
		"serverAddress": o.serverAddr,
		"user":          o.to,
	}

	err = t.Execute(&buf, inputs)
	if err != nil {
		return fmt.Errorf("can't execute template: %w", err)
	}

	o.msg = buf.Bytes()

	return nil
}

func (o *senderOptions) Run(cmd *cobra.Command) error {
	var err error

	cliflag.PrintFlags(cmd.Flags())

	err = smtp.SendMail(fmt.Sprintf("%s:%s", o.hostAddr, o.hostPort), o.auth, o.from, []string{o.to}, o.msg)
	if err != nil {
		return fmt.Errorf("can't send mail: %w", err)
	}

	return nil
}
