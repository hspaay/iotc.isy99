// Package internal handles node configuration commands
package internal

import "github.com/hspaay/iotc.golang/iotc"

// HandleConfigCommand for handling node configuration changes
// Not supported
func (app *IsyApp) HandleConfigCommand(node *iotc.NodeDiscoveryMessage, config iotc.NodeAttrMap) iotc.NodeAttrMap {
	app.logger.Infof("IsyApp.HandleConfigCommand for %s. Ignored as this isn't supported", node.Address)
	return nil
}
