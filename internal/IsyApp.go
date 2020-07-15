// Package internal for basic ISY99x Insteon home automation hub access
// This implements common sensors and switches
package internal

import (
	"github.com/iotdomain/iotdomain-go/publisher"
)

// ConfigDefaultPollIntervalSec for polling the gateway
const ConfigDefaultPollIntervalSec = 15 * 60

// AppID application name used for configuration file and default publisherID
const appID = "isy99"

const defaultGatewayID = "gateway"

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
	isyAPI *IsyAPI // ISY gateway access
}

// NewIsyApp creates the app
// This creates a node for the gateway
func NewIsyApp(config *IsyAppConfig, pub *publisher.Publisher) *IsyApp {
	app := IsyApp{
		config: config,
		pub:    pub,
		// gatewayNodeAddr: nodes.MakeNodeDiscoveryAddress(pub.Zone, config.PublisherID, GatewayID),
		isyAPI: NewIsyAPI(config.GatewayAddress, config.LoginName, config.Password),
	}
	if app.config.GatewayID == "" {
		app.config.GatewayID = defaultGatewayID
	}
	if app.config.PublisherID == "" {
		app.config.PublisherID = appID
	}
	pub.SetPollInterval(60, app.Poll)
	pub.SetNodeInputHandler(app.HandleInputCommand)
	pub.SetNodeConfigHandler(app.HandleConfigCommand)
	// // Discover the node(s) and outputs. Use default for republishing discovery
	// isyPub.SetDiscoveryInterval(0, app.Discover)

	return &app
}

// Run the publisher until the SIGTERM  or SIGINT signal is received
func Run() {
	appConfig := &IsyAppConfig{PublisherID: appID, GatewayID: defaultGatewayID}
	isyPub, _ := publisher.NewAppPublisher(appID, "", "", appConfig, true)

	app := NewIsyApp(appConfig, isyPub)
	app.SetupGatewayNode(isyPub)

	isyPub.Start()
	isyPub.WaitForSignal()
	isyPub.Stop()
}
