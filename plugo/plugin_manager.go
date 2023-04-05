package plugo

import (
	"errors"
	"net"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
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
	Plugins []*PluginProxy
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
func (pm *PluginManager) LoadPlugins(pluginPath string) {
	// check pluginPath is directory
	fileinfo, err := os.Stat(pluginPath)

	if err != nil {
		logrus.Fatal("pluginPath does not exist")
	}

	if !fileinfo.IsDir() {
		logrus.Fatal("pluginPath is not a directory")
	}

	// walk the directory and add plugins to the Plugins slice
	filepath.Walk(pluginPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if info.Mode().Perm()&0111 != 0 {
				pm.Plugins = append(pm.Plugins, &PluginProxy{
					Name:  info.Name(),
					Path:  path,
					State: PluginStateNotReady,
				})
				logrus.Info("Loaded plugin: ", info.Name())
			}
		}
		return nil
	})
}

// Serve starts the plugin manager and never returns.
func (pm *PluginManager) Serve() {
	listen, err := net.Listen("tcp", pm.Options.host+":"+pm.Options.port)
	if err != nil {
		os.Exit(1)
	}
	// close listener
	defer listen.Close()

	// start a goroutine to wait for connections and
	// handle incoming requests from plugins
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				os.Exit(1)
			}
			pm.conns = append(pm.conns, conn)
			defer conn.Close()

			go pm.handlePluginRegister(conn)
		}
	}()

	// start plugins
	pm.StartPlugins()

	// wait forever
	for {

	}
}

// handleRequest handles incoming requests from different plugins.
// one conn represents the connection to one of the plugins.
func (pm *PluginManager) handlePluginRegister(conn net.Conn) {
	// incoming request
	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
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
		case MSGPACK_RPC_TYPE_REQUEST:
			// unmarshal request message
			var request RequestMessage
			err = msgpack.Unmarshal(buffer, &request)

			if err != nil {
				logrus.Error("Failed to unmarshal request message: ", err)
				continue
			}

			if request.Method != "Register" {
				logrus.Error("Only Register method is allowed before plugin is registered")
				continue
			}
			CallMethodOfStruct(pm, request.Method, request.Params)

			response := ResponseMessage{
				Id:     request.Id,
				Type:   MSGPACK_RPC_TYPE_RESPONSE,
				Error:  "ok",
				Result: request.Params,
			}

			b, err := msgpack.Marshal(&response)
			if err != nil {
				logrus.Panic(err)
			}

			_, err = conn.Write(b)
			if err != nil {
				logrus.Fatal("Write data failed:", err.Error())
			}

			return
			break
		default:
			logrus.Error("Only request message is allowed before plugin is registered")
			break
		}

		if err != nil {
			panic(err)
		}

	}
}

// GetAllPluginNames returns all plugin names.
func (pm *PluginManager) GetAllPluginNames() []string {
	var names []string
	for _, plugin := range pm.Plugins {
		names = append(names, plugin.Name)
	}
	return names
}

// SyncCall calls a method of a plugin.
func (pm *PluginManager) SyncCall(name string, request *RequestMessage) error {
	plugin := pm.FindPlugin(name)
	if plugin == nil {
		logrus.Error("Plugin is not ready: ", name)
		return errors.New("Plugin is not ready: " + name)
	}

	if plugin.State != pluginStateBusy {
		plugin.State = pluginStateBusy
		plugin.Call(request)
		plugin.State = PluginStateReady
	} else {
		logrus.Error("Plugin is busy: ", name)
		return errors.New("Plugin is busy: " + name)
	}

	return nil
}

// StartPlugins starts all plugins.
func (pm *PluginManager) StartPlugins() {
	for _, plugin := range pm.Plugins {
		plugin.Start()
	}
}

// FindPlugin finds a registered plugin by name.
func (pm *PluginManager) FindPlugin(name string) *PluginProxy {
	for _, plugin := range pm.Plugins {
		if plugin.Name == name && plugin.State != PluginStateNotReady {
			return plugin
		}
	}
	return nil
}

// FindPluginByConn finds a registered plugin by connection.
func (pm *PluginManager) FindPluginByConn(conn net.Conn) *PluginProxy {
	for _, plugin := range pm.Plugins {
		if plugin.conn == conn {
			return plugin
		}
	}
	return nil
}

// Register registers a plugin and store the connection in the PluginProxy
// it is called by the plugin when it is ready to receive requests.
func (pm *PluginManager) Register(params map[string]interface{}) {
	for _, plugin := range pm.Plugins {
		// register plugin
		if plugin.State == PluginStateNotReady {
			if plugin.Name == params["Name"].(string) {
				for i, conn := range pm.conns {
					if conn.RemoteAddr().(*net.TCPAddr).Port == int(params["Address"].(uint16)) {
						plugin.conn = pm.conns[i]
						plugin.State = PluginStateReady
						logrus.Info("Registered plugin: ", plugin.Name)
						break
					}
				}
			}
		}
	}
}
