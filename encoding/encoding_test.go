// Copyright (C) 2018 Colin Walker
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnderscore(t *testing.T) {
	require.Equal(t, "created_at", underscore("CreatedAt"))
	require.Equal(t, "created", underscore("Created"))
	require.Equal(t, "api", underscore("API"))
	require.Equal(t, "test_api", underscore("Test_API"))
	require.Equal(t, "test_upper", underscore("TestUPPER"))
}
