package internal

import (
	"fmt"
	"time"

	"github.com/iotdomain/iotdomain-go/publisher"
	"github.com/iotdomain/iotdomain-go/types"
)

// ReadGateway reads the isy99 gateway device and its nodes
func (app *IsyApp) ReadGateway() error {
	pub := app.pub
	gwNodeID := app.config.GatewayID
	startTime := time.Now()
	isyDevice, err := app.isyAPI.ReadIsyGateway()
	endTime := time.Now()
	latency := endTime.Sub(startTime)

	prevStatus, _ := pub.GetNodeStatus(gwNodeID, types.NodeStatusRunState)
	if err != nil {
		// only report this once
		if prevStatus != types.NodeRunStateError {
			// gateway went down
			app.logger.Warningf("IsyApp.ReadGateway: ISY99x gateway is no longer reachable on address %s", app.isyAPI.address)
			pub.SetNodeStatus(gwNodeID, map[types.NodeStatus]string{
				types.NodeStatusRunState:  types.NodeRunStateError,
				types.NodeStatusLastError: "Gateway not reachable on address " + app.isyAPI.address,
			})
		}
		return err
	}

	pub.SetNodeStatus(gwNodeID, map[types.NodeStatus]string{
		types.NodeStatusRunState:    types.NodeRunStateReady,
		types.NodeStatusLastError:   "Connection restored to address " + app.isyAPI.address,
		types.NodeStatusLatencyMSec: fmt.Sprintf("%d", latency.Milliseconds()),
	})
	app.logger.Warningf("Isy99Adapter.ReadGateway: Connection restored to ISY99x gateway on address %s", app.isyAPI.address)

	// Update the info we have on the gateway
	pub.Nodes.SetNodeAttr(gwNodeID, map[types.NodeAttr]string{
		types.NodeAttrName:            isyDevice.configuration.Platform,
		types.NodeAttrSoftwareVersion: isyDevice.configuration.App + " - " + isyDevice.configuration.AppVersion,
		types.NodeAttrModel:           isyDevice.configuration.Product.Description,
		types.NodeAttrManufacturer:    isyDevice.configuration.DeviceSpecs.Make,
		// types.NodeAttrLocalIP:         isyDevice.network.Interface.IP,
		types.NodeAttrLocalIP: app.isyAPI.address,
		types.NodeAttrMAC:     isyDevice.configuration.Root.ID,
	})
	return nil
}

// SetupGatewayNode creates the gateway node if it doesn't exist
// This set the default gateway address in its configuration
func (app *IsyApp) SetupGatewayNode(pub *publisher.Publisher) {
	gwID := app.config.GatewayID
	app.logger.Infof("SetupGatewayNode. ID=%s", gwID)

	gatewayNode := pub.GetNodeByID(gwID)
	if gatewayNode == nil {
		pub.NewNode(gwID, types.NodeTypeGateway)
		gatewayNode = pub.GetNodeByID(gwID)
	}
	pub.Nodes.UpdateNodeConfig(gatewayNode.Address, types.NodeAttrLocalIP, &types.ConfigAttr{
		DataType:    types.DataTypeString,
		Description: "ISY gateway IP address",
		Secret:      true,
	})
	pub.Nodes.UpdateNodeConfig(gatewayNode.Address, types.NodeAttrLoginName, &types.ConfigAttr{
		DataType:    types.DataTypeString,
		Description: "ISY gateway login name",
		Secret:      true,
	})
	pub.Nodes.UpdateNodeConfig(gatewayNode.Address, types.NodeAttrPassword, &types.ConfigAttr{
		DataType:    types.DataTypeString,
		Description: "ISY gateway login password",
		Secret:      true,
	})

}
