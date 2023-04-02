package plugo

import (
	"fmt"
	"net"
	"os"
	"os/exec"
)

// PluginProxy .
type PluginProxy struct {
	Name       string
	Path       string
	Registered bool
	Cmd        *exec.Cmd
	conn       net.Conn
}

// Start starts the plugin.
func (pp *PluginProxy) Start() error {
	fmt.Println("Starting plugin: ", pp.Name)
	pp.Cmd = exec.Command(pp.Path)
	pp.Cmd.Stdout = os.Stdout
	pp.Cmd.Stderr = os.Stderr
	pp.Cmd.Start()
	fmt.Println("Started plugin: ", pp.Name, " with pid: ", pp.Cmd.Process.Pid)
	return nil
}

// Stop stops the plugin.
func (pp *PluginProxy) Stop() error {
	fmt.Println("Stopping plugin: ", pp.Name)
	pp.Cmd.Process.Kill()
	fmt.Println("Stopped plugin: ", pp.Name)
	return nil
}

// Restart restarts the plugin.
func (pp *PluginProxy) Restart() error {
	fmt.Println("Restarting plugin: ", pp.Name)
	pp.Stop()
	pp.Start()
	fmt.Println("Restarted plugin: ", pp.Name)
	return nil
}

// func (pp *PluginProxy) Call(msg Message) error {
// 	b, err := msgpack.Marshal(&msg)
// 	if err != nil {
// 		panic(err)
// 	}

// 	_, err = pp.conn.Write(b)
// 	if err != nil {
// 		println("Write data failed:", err.Error())
// 		os.Exit(1)
// 	}

// 	return nil
// }
