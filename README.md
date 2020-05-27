# IoTConnect ISY99x publisher

ISY99 is a gateway for the Insteon protocol. It is end of life but they are still around. It is replaced by the ISY944 which also works with this publisher, albeit limitated to what the ISY99 can do.

## Installation


## Configuration

Two configuration files are expected:
1. ~/bin/iotc/config/messenger.yaml which is used by with all publishers
2. ~/bin/iotc/config/onewire.yaml with EDS gateway address and login name/password. The gateway can also be configured through the gateway node 'gatewayAddress' configuration.

See config files in ./test as examples
