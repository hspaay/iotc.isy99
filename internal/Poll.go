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
		pub.UpdateNodeConfig(nodeID, iotc.NodeAttrName, &iotc.ConfigAttr{
			DataType:    iotc.DataTypeString,
			Description: "Name of ISY node",
			Default:     isyNode.Name,
		})
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
			pub.NewInput(nodeID, iotc.InputType(outputType), iotc.DefaultInputInstance)
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

// Poll polls the ISY gateway for updates to nodes and sensors
func (app *IsyApp) Poll(pub *publisher.Publisher) {
	err := app.ReadGateway()
	if err == nil {
		app.UpdateDevices()
	}
}
