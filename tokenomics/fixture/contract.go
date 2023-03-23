// SPDX-License-Identifier: ice License 1.0

package fixture

// Public API.

const (
	TestConnectorsOrder = 0
)

const (
	All StartLocalTestEnvironmentType = "all"
	DB  StartLocalTestEnvironmentType = "db"
	MB  StartLocalTestEnvironmentType = "mb"
)

type (
	StartLocalTestEnvironmentType string
)

// Private API.

const (
	applicationYAMLKey = "tokenomics"
)
