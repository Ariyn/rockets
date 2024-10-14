package rockets

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestMassObject_GetGravity(t *testing.T) {
	earth := MassObject{
		Position: Vector3D{0, 0, 0},
		Mass:     5.97219e+24,
		R:        6.375427e+06,
	}

	obj1 := MassObject{
		Position: Vector3D{earth.R, 0, 0},
		Mass:     1 * 1e3,
		R:        1,
	}

	gravity := obj1.GetGravity(&earth)
	assert.Equal(t, 9.80665, math.Round(gravity/obj1.Mass*1e5)/1e5)
}
