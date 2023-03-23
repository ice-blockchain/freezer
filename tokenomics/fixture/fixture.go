// SPDX-License-Identifier: ice License 1.0

package fixture

import (
	"testing"

	"github.com/testcontainers/testcontainers-go"

	connectorsfixture "github.com/ice-blockchain/wintr/connectors/fixture"
	messagebrokerfixture "github.com/ice-blockchain/wintr/connectors/message_broker/fixture"
	storagefixture "github.com/ice-blockchain/wintr/connectors/storage/fixture"
)

func StartLocalTestEnvironment(tp StartLocalTestEnvironmentType) {
	var connectors []connectorsfixture.TestConnector
	switch tp {
	case DB:
		connectors = append(connectors, newDBConnector())
	case MB:
		connectors = append(connectors, newMBConnector())
	case All:
		connectors = WTestConnectors()
	default:
		connectors = WTestConnectors()
	}
	connectorsfixture.
		NewTestRunner(applicationYAMLKey, nil, connectors...).
		StartConnectorsIndefinitely()
}

//nolint:gocritic // Because that's exactly what we want.
func RunTests(
	m *testing.M,
	dbConnector *storagefixture.TestConnector,
	mbConnector *messagebrokerfixture.TestConnector,
	lifeCycleHooks ...*connectorsfixture.ConnectorLifecycleHooks,
) {
	*dbConnector = newDBConnector()
	*mbConnector = newMBConnector()

	var connectorLifecycleHooks *connectorsfixture.ConnectorLifecycleHooks
	if len(lifeCycleHooks) == 1 {
		connectorLifecycleHooks = lifeCycleHooks[0]
	}

	connectorsfixture.
		NewTestRunner(applicationYAMLKey, connectorLifecycleHooks, *dbConnector, *mbConnector).
		RunTests(m)
}

func WTestConnectors() []connectorsfixture.TestConnector {
	return []connectorsfixture.TestConnector{newDBConnector(), newMBConnector()}
}

func RTestConnectors() []connectorsfixture.TestConnector {
	return []connectorsfixture.TestConnector{newDBConnector()}
}

func newDBConnector() storagefixture.TestConnector {
	return storagefixture.NewTestConnector(applicationYAMLKey, TestConnectorsOrder)
}

func newMBConnector() messagebrokerfixture.TestConnector {
	return messagebrokerfixture.NewTestConnector(applicationYAMLKey, TestConnectorsOrder)
}

func RContainerMounts() []func(projectRoot string) testcontainers.ContainerMount {
	return nil
}

func WContainerMounts() []func(projectRoot string) testcontainers.ContainerMount {
	return nil
}
