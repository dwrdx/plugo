package plugo

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

// Config options.
type ConfigOptions struct {
	host string
	port string
}

// PluginManager .
type PluginManager struct {
	Options ConfigOptions
	Plugins []PluginProxy
	conns   []net.Conn
}

// NewPluginManager creates a new PluginManager.
func NewPluginManager(host string, port string) *PluginManager {
	return &PluginManager{
		Options: ConfigOptions{
			host: host,
			port: port,
		},
	}
}

// LoadPlugins loads plugins from a directory and adds them to the Plugins slice.
func (pm *PluginManager) LoadPlugins(pluginPath string) error {
	// check pluginPath is directory
	fileinfo, err := os.Stat(pluginPath)

	if err != nil {
		return err
	}
	if !fileinfo.IsDir() {
		return fmt.Errorf("pluginPath is not a directory")
	}

	// walk the directory and add plugins to the Plugins slice
	filepath.Walk(pluginPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if info.Mode().Perm()&0111 != 0 {
				pm.Plugins = append(pm.Plugins, PluginProxy{
					Name: info.Name(),
					Path: path,
				})
				fmt.Println("Loaded plugin: ", info.Name())
			}
		}
		return nil
	})

	return nil
}

// Serve starts the plugin manager.
func (pm *PluginManager) Serve() {
	listen, err := net.Listen("tcp", pm.Options.host+":"+pm.Options.port)
	if err != nil {
		os.Exit(1)
	}
	// close listener
	defer listen.Close()

	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				os.Exit(1)
			}
			pm.conns = append(pm.conns, conn)
			fmt.Println(pm.conns)
			defer conn.Close()
			go pm.handleRequest(conn)
			go pm.sendRequest(conn)
		}
	}()

	go pm.StartPlugins()

	for {
		time.Sleep(1 * time.Second)
	}
}

func (pm *PluginManager) handleRequest(conn net.Conn) {
	// incoming request
	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			return
		}
		// fmt.Println("Received message:", string(buffer))

		var requestMsg RequestMessage
		err = msgpack.Unmarshal(buffer, &requestMsg)
		// fmt.Println(requestMsg)
		// fmt.Println(reflect.ValueOf(requestMsg.Params))

		f := reflect.ValueOf(pm).MethodByName(requestMsg.Method)
		f.Call([]reflect.Value{reflect.ValueOf(requestMsg.Params)})

		if err != nil {
			panic(err)
		}

	}
}

func (pm *PluginManager) sendRequest(conn net.Conn) {
	// incoming request
	for {
		// buffer := make([]byte, 1024)
		// _, err := conn.Read(buffer)
		// if err != nil {
		// 	return
		// }
		// fmt.Println("Received message:", string(buffer))

		_, err := conn.Write([]byte("Hello World!"))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		time.Sleep(1 * time.Second)

		// var item Item
		// err = msgpack.Unmarshal(buffer, &item)
		// if err != nil {
		// 	panic(err)
		// }

	}
}

func (pm *PluginManager) StartPlugins() {
	for _, plugin := range pm.Plugins {
		plugin.Start()
	}
}

func (pm *PluginManager) Register(params map[string]interface{}) {
	fmt.Println(params)
	for _, plugin := range pm.Plugins {
		if plugin.Name == params["Name"].(string) {
			for _, conn := range pm.conns {
				if conn.RemoteAddr().(*net.TCPAddr).Port == int(params["Address"].(uint16)) {
					plugin.Registered = true
					fmt.Println("Registered plugin: ", plugin.Name)
					plugin.conn = conn
					break
				}
			}
		}
	}
}
