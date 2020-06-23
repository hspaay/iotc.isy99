module iotd.isy99

go 1.13

require (
	github.com/iotdomain/iotdomain-go v0.0.0-20200521044650-7be324a29524
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.6.0

)

// Temporary for testing iotdomain-go
replace github.com/iotdomain/iotdomain-go => ../iotdomain-go
