// Package internal for basic ISY99x Insteon home automation hub access
// This implements common sensors and switches
package internal

import (
	"github.com/hspaay/iotc.golang/publisher"
	"github.com/sirupsen/logrus"
)

// ConfigDefaultPollIntervalSec for polling the gateway
const ConfigDefaultPollIntervalSec = 15 * 60

// AppID application name used for configuration file and default publisherID
const appID = "isy99"

const gatewayID = "gateway"

// IsyAppConfig with application state, loaded from isy99.yaml
type IsyAppConfig struct {
	GatewayAddress string `yaml:"gatewayAddress"` // gateway IP address
	GatewayID      string `yaml:"gatewayId"`      // default is "gateway"
	LoginName      string `yaml:"login"`          // gateway login
	Password       string `yaml:"password"`       // gateway password
	PublisherID    string `yaml:"publisherId"`    // default is app ID
}

// IsyApp adapter main class
// to access multiple gatewways, run additional instances, or modify this code for multiple isyAPI instances
type IsyApp struct {
	config *IsyAppConfig
	pub    *publisher.Publisher
	logger *logrus.Logger
	isyAPI *IsyAPI // ISY gateway access
}

// // LoadConfiguration loads config for the adapter
// func (adapter *Isy99App) LoadConfiguration(configFolder string) error {
// 	err := adapter.MyZoneService.Load("isy99", configFolder, true)

// 	// ensure device is setup correctly after loading configuration
// 	if err == nil {
// 		adapter.PublisherNode().SetConfigDefault(nodes.AttrNamePollInterval, ConfigDefaultPollIntervalSec, nodes.DataTypeInt, "Interval in seconds to scan for new nodes")

// 		gateway := adapter.GatewayNode()
// 		// Set default data type and description of gateway parameters
// 		gateway.SetConfigDefault(nodes.AttrNameAddress, "", nodes.DataTypeString, "Hostname or IP address of the ISY gateway")
// 		config := gateway.SetConfigDefault(nodes.AttrNameLoginName, "", nodes.DataTypeString, "Secret login name of the ISY gateway")
// 		config.Secret = true
// 		config = gateway.SetConfigDefault(nodes.AttrNamePassword, "", nodes.DataTypeString, "Secret password of the ISY gateway")
// 		config.Secret = true
// 		adapter.isyAPI.log = adapter.Logger()
// 	}
// 	return err
// }

// // Start the module.
// // This loads the configuration and start polling and publishing sensor data in the background
// func (adapter *Isy99App) Start() error {
// 	interval, _ := adapter.PublisherNode().GetConfigInt(nodes.AttrNamePollInterval)
// 	err := adapter.MyZoneService.Start(adapter.commandHandler, nil, adapter.Poll, interval)
// 	adapter.Poll()
// 	return err
// }

// // Stop the adapter and its polling
// //func (adapter *Isy99Adapter) Stop() {
// //  adapter.MyZoneService.Stop()
// //}

// // NewIsy99Adapter returns a new instance of Isy99Adapter module
// func NewIsy99Adapter() *Isy99App {
// 	adapter := new(Isy99App)
// 	return adapter
// }

// NewIsyApp creates the app
// This creates a node for the gateway
func NewIsyApp(config *IsyAppConfig, pub *publisher.Publisher) *IsyApp {
	app := IsyApp{
		config: config,
		pub:    pub,
		logger: pub.Logger,
		// gatewayNodeAddr: nodes.MakeNodeDiscoveryAddress(pub.Zone, config.PublisherID, GatewayID),
		isyAPI: &IsyAPI{},
	}
	app.config.PublisherID = appID
	app.isyAPI.log = pub.Logger
	return &app
}

// Run the publisher until the SIGTERM  or SIGINT signal is received
func Run() {
	appConfig := &IsyAppConfig{PublisherID: appID, GatewayID: gatewayID}
	isyPub, _ := publisher.NewAppPublisher(appID, "", appConfig, true)

	app := NewIsyApp(appConfig, isyPub)
	app.SetupGatewayNode(isyPub)

	isyPub.SetPollInterval(60, app.Poll)
	isyPub.SetNodeInputHandler(app.InputHandler)
	// isyPub.SetNodeConfigHandler(app.ConfigHandler)

	// // Discover the node(s) and outputs. Use default for republishing discovery
	// isyPub.SetDiscoveryInterval(0, app.Discover)

	isyPub.Start()
	isyPub.WaitForSignal()
	isyPub.Stop()
}
