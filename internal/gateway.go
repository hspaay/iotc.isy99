package internal

import (
	"github.com/hspaay/iotc.golang/iotc"
	"github.com/hspaay/iotc.golang/publisher"
)

// ReadGateway reads the isy99 gateway device and its nodes
func (app *IsyApp) ReadGateway() error {
	pub := app.pub
	gwNodeID := app.config.GatewayID

	isyDevice, err := app.isyAPI.ReadIsyGateway()

	prevStatus, _ := pub.GetNodeStatus(gwNodeID, iotc.NodeStatusRunState)
	if err != nil {
		// only report this once
		if prevStatus != iotc.NodeRunStateError {
			// gateway went down
			app.logger.Warningf("IsyApp.ReadGateway: ISY99x gateway is no longer reachable on address %s", app.isyAPI.address)
			pub.SetNodeStatus(gwNodeID, map[iotc.NodeStatus]string{
				iotc.NodeStatusRunState:  iotc.NodeRunStateError,
				iotc.NodeStatusLastError: "Gateway not reachable on address " + app.isyAPI.address,
			})
		}
		return err
	}

	if prevStatus != iotc.NodeRunStateReady {
		// gateway came back
		pub.SetNodeStatus(gwNodeID, map[iotc.NodeStatus]string{
			iotc.NodeStatusRunState:  iotc.NodeRunStateReady,
			iotc.NodeStatusLastError: "Connection restored to address " + app.isyAPI.address,
		})
		app.logger.Warningf("Isy99Adapter.ReadGateway: Connection restored to ISY99x gateway on address %s", app.isyAPI.address)
	}

	// Update the info we have on the gateway
	pub.Nodes.SetNodeAttr(gwNodeID, map[iotc.NodeAttr]string{
		iotc.NodeAttrName:            isyDevice.configuration.Platform,
		iotc.NodeAttrSoftwareVersion: isyDevice.configuration.App + " - " + isyDevice.configuration.AppVersion,
		iotc.NodeAttrModel:           isyDevice.configuration.Product.Description,
		iotc.NodeAttrManufacturer:    isyDevice.configuration.DeviceSpecs.Make,
		// iotc.NodeAttrLocalIP:         isyDevice.network.Interface.IP,
		iotc.NodeAttrLocalIP: app.isyAPI.address,
		iotc.NodeAttrMAC:     isyDevice.configuration.Root.ID,
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
		pub.NewNode(gwID, iotc.NodeTypeGateway)
		gatewayNode = pub.GetNodeByID(gwID)
	}
	pub.Nodes.UpdateNodeConfig(gatewayNode.Address, iotc.NodeAttrLocalIP, &iotc.ConfigAttr{
		Datatype:    iotc.DataTypeString,
		Description: "ISY gateway IP address",
		Secret:      true,
	})
	pub.Nodes.UpdateNodeConfig(gatewayNode.Address, iotc.NodeAttrLoginName, &iotc.ConfigAttr{
		Datatype:    iotc.DataTypeString,
		Description: "ISY gateway login name",
		Secret:      true,
	})
	pub.Nodes.UpdateNodeConfig(gatewayNode.Address, iotc.NodeAttrPassword, &iotc.ConfigAttr{
		Datatype:    iotc.DataTypeString,
		Description: "ISY gateway login password",
		Secret:      true,
	})

}
