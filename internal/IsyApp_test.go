package internal

import (
	"testing"
	"time"

	"github.com/hspaay/iotc.golang/iotc"
	"github.com/hspaay/iotc.golang/messenger"
	"github.com/hspaay/iotc.golang/publisher"
	"github.com/stretchr/testify/assert"
)

const testConfigFolder = "../test"
const testCacheFolder = "../test/cache"
const testData = "isy99-testdata.xml"
const deckLightsID = "15 2D A 1"

var messengerConfig = &messenger.MessengerConfig{Domain: "test"}
var appConfig = &IsyAppConfig{}

func TestLoadConfig(t *testing.T) {
	_, err := publisher.NewAppPublisher(appID, testConfigFolder, testCacheFolder, appConfig, false)
	assert.NoError(t, err)
	assert.Equal(t, "isy99", appConfig.PublisherID)
	assert.Equal(t, "gateway", appConfig.GatewayID)
}

// Read ISY device and check if more than 1 node is returned. A minimum of 1 is expected if the device is online with
// an additional node for each connected node.
func TestReadIsy(t *testing.T) {
	_, err := publisher.NewAppPublisher(appID, testConfigFolder, testCacheFolder, appConfig, false)
	assert.NoError(t, err)

	isyAPI := NewIsyAPI(appConfig.GatewayAddress, appConfig.LoginName, appConfig.Password)
	// use a simulation file
	isyAPI.address = "file://../test/gateway-config.xml"
	isyDevice, err := isyAPI.ReadIsyGateway()
	assert.NoError(t, err)
	assert.NotEmptyf(t, isyDevice.configuration.AppVersion, "Expected an application version")

	// use a simulation file
	isyAPI.address = "file://../test/gateway-nodes.xml"
	isyNodes, err := isyAPI.ReadIsyNodes()
	if assert.NoError(t, err) {
		assert.True(t, len(isyNodes.Nodes) > 5, "Expected 5 ISY nodes. Got fewer.")
	}
}

func TestPollOnce(t *testing.T) {
	pub, err := publisher.NewAppPublisher(appID, testConfigFolder, testCacheFolder, appConfig, false)
	assert.NoError(t, err)

	app := NewIsyApp(appConfig, pub)
	pub.Start()
	assert.NoError(t, err)
	app.Poll(pub)
	time.Sleep(3 * time.Second)
	pub.Stop()
}

// This simulates the switch
func TestSwitch(t *testing.T) {
	pub, err := publisher.NewAppPublisher(appID, testConfigFolder, testCacheFolder, appConfig, false)
	assert.NoError(t, err)

	app := NewIsyApp(appConfig, pub)
	app.SetupGatewayNode(pub)

	// FIXME: load isy nodes from file

	pub.Start()
	assert.NoError(t, err)
	app.Poll(pub)
	// some time to publish stuff
	time.Sleep(3 * time.Second)

	// throw a switch
	deckSwitch := pub.GetNodeByID(deckLightsID)
	if assert.NotNil(t, deckSwitch) {
		switchInput := pub.Inputs.GetInput(deckSwitch.Address, iotc.InputTypeSwitch, iotc.DefaultInputInstance)
		// switchInput := deckSwitch.GetInput(iotc.InputTypeSwitch)

		app.logger.Infof("TestSwitch: --- Switching deck switch %s OFF", deckSwitch.Address)
		pubKey := pub.GetPublisherKey(switchInput.Address)
		pub.PublishSetInput(switchInput.Address, "false", pubKey)
		assert.NoError(t, err)
		time.Sleep(2 * time.Second)

		// fetch result
		switchOutput := pub.GetOutputByType(deckLightsID, iotc.OutputTypeSwitch, iotc.DefaultInputInstance)
		// switchOutput := deckSwitch.GetOutput(iotc.InputTypeSwitch)
		if assert.NotNil(t, switchOutput) {
			outputValue := pub.OutputValues.GetOutputValueByAddress(switchOutput.Address)
			assert.Equal(t, "false", outputValue.Value)

			app.logger.Infof("TestSwitch: --- Switching deck switch %s ON", deckSwitch.Address)
			if assert.NotNil(t, switchInput) {
				pub.PublishSetInput(switchInput.Address, "true", pubKey)
			}
			time.Sleep(2 * time.Second)
			outputValue = pub.OutputValues.GetOutputValueByAddress(switchOutput.Address)
			assert.Equal(t, "true", outputValue.Value)

			// be nice and turn the light back off
			pub.PublishSetInput(switchInput.Address, "false", pubKey)
		}
	}
	pub.Stop()
}

func TestStartStop(t *testing.T) {
	pub, err := publisher.NewAppPublisher(appID, testConfigFolder, testCacheFolder, appConfig, false)
	assert.NoError(t, err)

	// app := NewIsyApp(appConfig, pub)

	pub.Start()
	time.Sleep(time.Second * 100)
	pub.Stop()
}
