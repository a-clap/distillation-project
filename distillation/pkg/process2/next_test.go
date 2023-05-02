/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package process2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEndConditionTemperature(test *testing.T) {
	t := require.New(test)
	retTime := int64(0)
	retTemperature := float64(0)

	getTemp := func() float64 {
		return retTemperature
	}
	getTm := func() int64 {
		return retTime
	}

	e := newEndConditionTemperature(100, getTm, 10.0, getTemp)

	// Return temperature not over threshold
	end, left := e.end()
	t.False(end)
	t.EqualValues(100, left)

	// Temperature over threshold
	retTemperature = 10.1
	end, left = e.end()
	t.False(end)
	t.EqualValues(100, left)

	// Temperature still over threshold
	// But time is not satisfied
	// Unix() should be called once - byTime is not resetted
	retTime = 1
	end, left = e.end()
	t.False(end)
	t.EqualValues(99, left)

	// Temperature again below threshold
	retTemperature = 10.0
	end, left = e.end()
	t.False(end)
	t.EqualValues(100, left)

	// Happy path - temperature over threshold, and time get elapsed
	retTemperature = 10.1
	end, left = e.end()
	t.False(end)
	t.EqualValues(100, left)

	retTime = 101
	end, left = e.end()
	t.True(end)
	t.EqualValues(0, left)
}

func TestEndConditionTime(test *testing.T) {
	t := require.New(test)

	retTime := int64(0)
	getTm := func() int64 {
		return retTime
	}

	time := newEndConditionTime(100, getTm)

	end, left := time.end()
	t.False(end)
	t.EqualValues(100, left)

	retTime = 99
	end, left = time.end()
	t.False(end)
	t.EqualValues(1, left)

	// Now should be expired
	retTime = 100
	end, left = time.end()
	t.True(end)
	t.EqualValues(0, left)

	// Reset, should start from begin
	time.reset()
	end, left = time.end()
	t.False(end)
	t.EqualValues(100, left)

	// Now retTime was 100
	// To satisfy next, we must return 200
	retTime = 199
	end, left = time.end()
	t.False(end)
	t.EqualValues(1, left)

	// Much more time expired
	retTime = 300
	end, left = time.end()
	t.True(end)
	t.EqualValues(0, left)

}
