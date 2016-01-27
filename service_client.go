package converger

import (
	"time"

	"github.com/cloudfoundry-incubator/consuladapter"
	"github.com/cloudfoundry-incubator/locket"
	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/ifrit"
)

const ConvergerLockSchemaKey = "converge_lock"

func ConvergerLockSchemaPath() string {
	return locket.LockSchemaPath(ConvergerLockSchemaKey)
}

type ServiceClient interface {
	NewConvergerLockRunner(logger lager.Logger, convergerID string, retryInterval, lockTTL time.Duration) ifrit.Runner
}

type serviceClient struct {
	consulClient consuladapter.Client
	clock        clock.Clock
}

func NewServiceClient(consulClient consuladapter.Client, clock clock.Clock) ServiceClient {
	return serviceClient{
		consulClient: consulClient,
		clock:        clock,
	}
}

func (c serviceClient) NewConvergerLockRunner(logger lager.Logger, convergerID string, retryInterval, lockTTL time.Duration) ifrit.Runner {
	return locket.NewLock(logger, c.consulClient, ConvergerLockSchemaPath(), []byte(convergerID), c.clock, retryInterval, lockTTL)
}
