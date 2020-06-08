package internal

import (
	"github.com/hspaay/iotc.golang/iotc"
	"github.com/hspaay/iotc.golang/nodes"
	"github.com/hspaay/iotc.golang/publisher"
)

// 		gateway := adapter.GatewayNode()
// 		// Set default data type and description of gateway parameters
// 		gateway.SetConfigDefault(nodes.AttrNameAddress, "", nodes.DataTypeString, "Hostname or IP address of the ISY gateway")
// 		config := gateway.SetConfigDefault(nodes.AttrNameLoginName, "", nodes.DataTypeString, "Secret login name of the ISY gateway")
// 		config.Secret = true
// 		config = gateway.SetConfigDefault(nodes.AttrNamePassword, "", nodes.DataTypeString, "Secret password of the ISY gateway")
// 		config.Secret = true
// 		adapter.isyAPI.log = adapter.Logger()

// SetupGatewayNode creates the gateway node if it doesn't exist
// This set the default gateway address in its configuration
func (app *IsyApp) SetupGatewayNode(pub *publisher.Publisher) {
	gwID := app.config.GatewayID
	app.logger.Infof("SetupGatewayNode. ID=%s", gwID)

	gatewayNode := pub.GetNodeByID(gwID)
	if gatewayNode == nil {
		pub.NewNode(gwID, iotc.NodeTypeGateway)
		gatewayNode = pub.GetNodeByID(gwID)
	}
	config := pub.NewNodeConfig(gwID, iotc.NodeAttrAddress, iotc.DataTypeString, "ISY gateway IP address", "")
	config = nodes.NewNodeConfig(iotc.NodeAttrLoginName, iotc.DataTypeString, "ISY gateway login name", "")
	config.Secret = true
	pub.Nodes.UpdateNodeConfig(gatewayNode.Address, config)

	config = nodes.NewNodeConfig(iotc.NodeAttrPassword, iotc.DataTypeString, "ISY gateway login password", "")
	config.Secret = true
	pub.Nodes.UpdateNodeConfig(gatewayNode.Address, config)
}
