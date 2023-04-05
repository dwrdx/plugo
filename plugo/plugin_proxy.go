package plugo

import (
	"net"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
)

type PluginProxyState int

const (
	PluginStateNotReady PluginProxyState = 0
	PluginStateReady    PluginProxyState = 1
	pluginStateBusy     PluginProxyState = 2
)

// PluginProxy .
type PluginProxy struct {
	Name  string
	Path  string
	State PluginProxyState
	Cmd   *exec.Cmd
	conn  net.Conn
}

// Start starts the plugin.
func (pp *PluginProxy) Start() error {
	logrus.Info("Starting plugin: ", pp.Name)
	pp.Cmd = exec.Command(pp.Path)
	pp.Cmd.Stdout = os.Stdout
	pp.Cmd.Stderr = os.Stderr
	pp.Cmd.Start()
	go pp.HandleConnection()
	logrus.Info("Started plugin: ", pp.Name, " with pid: ", pp.Cmd.Process.Pid)
	return nil
}

// Stop stops the plugin.
func (pp *PluginProxy) Stop() error {
	logrus.Info("Stopping plugin: ", pp.Name)
	pp.Cmd.Process.Kill()
	logrus.Info("Stopped plugin: ", pp.Name)
	return nil
}

// Restart restarts the plugin.
func (pp *PluginProxy) Restart() error {
	logrus.Info("Restarting plugin: ", pp.Name)
	pp.Stop()
	pp.Start()
	logrus.Info("Restarted plugin: ", pp.Name)
	return nil
}

func (pp *PluginProxy) HandleConnection() {
	for {
		if pp.State == PluginStateReady {
			buffer := make([]byte, 1024)
			_, err := pp.conn.Read(buffer)
			if err != nil {
				logrus.Error("Failed to read from connection: ", err)
				continue
			}

			var header MessageHeader
			err = msgpack.Unmarshal(buffer, &header)

			if err != nil {
				logrus.Error("Failed to unmarshal received message: ", err)
				continue
			}

			switch header.Type {
			case MSGPACK_RPC_TYPE_RESPONSE:
				// unmarshal response message
				var response ResponseMessage
				err = msgpack.Unmarshal(buffer, &response)

				if err != nil {
					logrus.Error("Failed to unmarshal response message: ", err)
					continue
				}

				logrus.Info("Received response from plugin: ", pp.Name)
				break

			case MSGPACK_RPC_TYPE_NOTIFY:
				logrus.Error("Notify message is not supported")
				break
			default:
				break
			}

		}
	}
}

// Call calls a method on the plugin.
func (pp *PluginProxy) Call(msg *RequestMessage) error {
	b, err := msgpack.Marshal(&msg)
	if err != nil {
		logrus.Panic(err)
	}

	_, err = pp.conn.Write(b)
	if err != nil {
		logrus.Fatal("Write data failed:", err.Error())
	}

	return nil
}

// Wait waits for the plugin to finish a call.
func (pp *PluginProxy) Wait() error {
	return nil
}
