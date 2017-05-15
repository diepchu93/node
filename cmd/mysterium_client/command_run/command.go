package command_run

import (
	"github.com/mysterium/node/bytescount"
	"github.com/mysterium/node/openvpn"
	"github.com/mysterium/node/server"
	"io"
	"time"
)

type commandRun struct {
	output      io.Writer
	outputError io.Writer

	mysteriumClient server.Client
	vpnClient       *openvpn.Client
}

func (cmd *commandRun) Run(options CommandOptions) error {
	vpnSession, err := cmd.mysteriumClient.SessionCreate(options.NodeKey)
	if err != nil {
		return err
	}

	vpnConfig, err := openvpn.NewClientConfigFromString(
		vpnSession.ConnectionConfig,
		options.DirectoryRuntime+"/client.ovpn",
	)
	if err != nil {
		return err
	}

	cmd.vpnClient = openvpn.NewClient(
		vpnConfig,
		options.DirectoryRuntime,
		bytescount.NewMiddleware(cmd.mysteriumClient, vpnSession.Id, 1*time.Minute),
	)
	if err := cmd.vpnClient.Start(); err != nil {
		return err
	}

	return nil
}

func (cmd *commandRun) Wait() error {
	return cmd.vpnClient.Wait()
}

func (cmd *commandRun) Kill() {
	cmd.vpnClient.Stop()
}
