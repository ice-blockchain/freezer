// SPDX-License-Identifier: BUSL-1.1

package fixture

<<<<<<< HEAD
import "errors"

var errFixtureCleanupFailed = errors.New("economy fixture cleanup failed")
=======
import (
	"errors"
)

var errFixtureCleanup = errors.New("economy fixture cleanup failed")
>>>>>>> Start Mining endpoint implementation. Makefile coverage was chang…
