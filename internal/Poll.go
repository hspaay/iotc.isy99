// Package internal to poll the ISY for device, node and parameter information
package internal

import (
	"fmt"
	"strings"

	"github.com/hspaay/iotc.golang/iotc"
	"github.com/hspaay/iotc.golang/publisher"
)

// IsyURL to contact ISY99x gateway
const IsyURL = "http://%s/rest/nodes"

// readIsyNodesValues reads the ISY Node values
// This will run http get on http://address/rest/nodes
// address is the ISY hostname or ip address.
// returns a node object and possible error
func (app *IsyApp) readIsyNodesValues(address string) (*IsyNodes, error) {
	isyURL := fmt.Sprintf(IsyURL, address)
	isyNodes := IsyNodes{}
	err := app.isyAPI.isyRequest(isyURL, &isyNodes)
	if err != nil {
		return nil, err
	}
	return &isyNodes, nil
}

// updateDevice updates the node discovery and output value from the provided isy node
func (app *IsyApp) updateDevice(isyNode *IsyNode) {
	nodeID := isyNode.Address
	pub := app.pub
	hasInput := false
	outputValue := isyNode.Property.Value

	// What node are we dealing with?
	deviceType := iotc.NodeTypeUnknown
	outputType := iotc.OutputTypeOnOffSwitch
	switch isyNode.Property.ID {
	case "ST":
		deviceType = iotc.NodeTypeOnOffSwitch
		outputType = iotc.OutputTypeOnOffSwitch
		hasInput = true
		if outputValue == "0" || strings.ToLower(outputValue) == "false" {
			outputValue = "false"
		} else {
			outputValue = "true"
		}
		break
	case "OL":
		deviceType = iotc.NodeTypeDimmer
		outputType = iotc.OutputTypeDimmer
		hasInput = true
		break
	case "RR":
		deviceType = iotc.NodeTypeUnknown
		break
	}
	// Add new discoveries
	node := pub.GetNodeByID(nodeID)
	if node == nil {
		pub.NewNode(nodeID, iotc.NodeType(deviceType))
		pub.NewNodeConfig(nodeID, iotc.NodeAttrName, iotc.DataTypeString, "Name of ISY node", isyNode.Name)
		pub.NewNodeConfig(nodeID, iotc.NodeAttrProduct, iotc.DataTypeString, "Device product name", isyNode.Type)
		pub.SetNodeStatus(nodeID, map[iotc.NodeStatus]string{
			iotc.NodeStatusRunState: iotc.NodeRunStateReady,
		})
	}

	output := pub.GetOutputByType(nodeID, outputType, iotc.DefaultOutputInstance)
	if output == nil {
		// Add an output and optionally an input for the node.
		// Most ISY nodes have only a single sensor. This is a very basic implementation.
		// Is it worth adding multi-sensor support?
		// https://wiki.universal-devices.com/index.php?title=ISY_Developers:API:REST_Interface#Properties
		pub.NewOutput(nodeID, outputType, iotc.DefaultOutputInstance)
		if hasInput {
			pub.NewInput(nodeID, outputType, iotc.DefaultInputInstance)
		}
	}

	//if output.Value() != isyNode.Property.Value {
	//	// this compares 0 with false, so lots of noise
	//	adapter.Logger().Debugf("Isy99Adapter.updateDevice. Update node %s, output %s[%s] from %s to %s",
	//		node.Id, output.IOType, output.Instance, output.Value(), isyNode.Property.Value)
	//}
	// let the adapter decide whether to repeat the same value based on config
	pub.UpdateOutputValue(nodeID, outputType, iotc.DefaultOutputInstance, outputValue)

}

// UpdateDevices discover ISY Nodes from config and ISY gateway
func (app *IsyApp) UpdateDevices() {
	// Discover the ISY nodes
	isyNodes, err := app.isyAPI.ReadIsyNodes()
	if err != nil {
		// Unexpected. What to do now?
		app.logger.Warningf("DiscoverNodes: Error reading nodes: %s", err)
		return
	}
	// Update new or changed ISY nodes
	for _, isyNode := range isyNodes.Nodes {
		app.updateDevice(isyNode)
	}
}

// ReadGateway reads the isy99 gateway device and its nodes
func (app *IsyApp) ReadGateway() error {
	pub := app.pub
	gwNodeID := app.config.GatewayID
	// gateway := app.GatewayNode()

	// app.isyAPI.address, _ = pub.GetNodeConfigValue(gwNodeID, iotc.NodeAttrAddress, app.config.GatewayAddress)
	// app.isyAPI.login, _ = pub.GetNodeConfigValue(gwNodeID, iotc.NodeAttrLoginName, app.config.LoginName)
	// app.isyAPI.password, _ = pub.GetNodeConfigValue(gwNodeID, iotc.NodeAttrPassword, app.config.Password)

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

// Poll polls the ISY gateway for updates to nodes and sensors
func (app *IsyApp) Poll(pub *publisher.Publisher) {
	err := app.ReadGateway()
	if err == nil {
		app.UpdateDevices()
	}
}
