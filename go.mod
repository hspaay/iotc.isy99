module iotc.isy99

go 1.13

require (
	github.com/hspaay/iotc.golang v0.0.0-20200521044650-7be324a29524
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.6.0

)

// Temporary for testing iotc.golang
replace github.com/hspaay/iotc.golang => ../iotc.golang
