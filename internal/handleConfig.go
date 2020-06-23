// Package internal handles node configuration commands
package internal

import "github.com/iotdomain/iotdomain-go/types"

// HandleConfigCommand for handling node configuration changes
// Not supported
func (app *IsyApp) HandleConfigCommand(node *types.NodeDiscoveryMessage, config types.NodeAttrMap) types.NodeAttrMap {
	app.logger.Infof("IsyApp.HandleConfigCommand for %s. Ignored as this isn't supported", node.Address)
	return nil
}
