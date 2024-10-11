package rockets

import (
	"math"
	"reflect"
	"testing"
)

func TestVector3D_Constants(t *testing.T) {
	tests := []struct {
		name   string
		target Vector3D
		want   Vector3D
	}{
		{
			name:   "ZeroVector는 0, 0, 0 벡터를 나타낸다.",
			target: ZeroVector,
			want:   Vector3D{X: 0, Y: 0, Z: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.target; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Constant() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector3D_Mul(t *testing.T) {
	tests := []struct {
		name   string
		fields Vector3D
		arg    float64
		want   Vector3D
	}{
		{
			name:   "1을 곱했을때, 그대로의 벡터를 반환한다.",
			fields: Vector3D{X: 1, Y: 2, Z: 3},
			arg:    1,
			want:   Vector3D{X: 1, Y: 2, Z: 3},
		},
		{
			name:   "스칼라값을 곱했을때, 벡터의 각 성분이 곱해진 벡터를 반환한다.",
			fields: Vector3D{X: 1, Y: 2, Z: 3},
			arg:    2,
			want:   Vector3D{X: 2, Y: 4, Z: 6},
		},
		{
			name:   "음수를 곱했을때, 벡터의 각 성분이 음수가 된 벡터를 반환한다.",
			fields: Vector3D{X: 1, Y: 2, Z: 3},
			arg:    -2,
			want:   Vector3D{X: -2, Y: -4, Z: -6},
		},
		{
			name:   "0을 곱했을때, 0벡터가 반환된다.",
			fields: Vector3D{X: 1, Y: 2, Z: 3},
			arg:    0,
			want:   ZeroVector,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector3D{
				X: tt.fields.X,
				Y: tt.fields.Y,
				Z: tt.fields.Z,
			}
			if got := v.Mul(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mul() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector3D_Add(t *testing.T) {
	tests := []struct {
		name   string
		fields Vector3D
		arg    Vector3D
		want   Vector3D
	}{
		{
			name:   "두 벡터를 더했을 때, 두 벡터의 각 성분이 더해진 벡터를 반환한다.",
			fields: Vector3D{X: 1, Y: 2, Z: 3},
			arg:    Vector3D{X: 4, Y: 5, Z: 6},
			want:   Vector3D{X: 5, Y: 7, Z: 9},
		},
		{
			name:   "양수인 벡터에 음수인 벡터를 더했을 때, 두 벡터의 각 성분이 더해진 벡터를 반환한다.",
			fields: Vector3D{X: 1, Y: 2, Z: 3},
			arg:    Vector3D{X: -4, Y: -5, Z: -6},
			want:   Vector3D{X: -3, Y: -3, Z: -3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector3D{
				X: tt.fields.X,
				Y: tt.fields.Y,
				Z: tt.fields.Z,
			}
			if got := v.Add(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector3D_Dot(t *testing.T) {
	tests := []struct {
		name   string
		fields Vector3D
		arg    Vector3D
		want   float64
	}{
		{
			name:   "두 벡터를 곱했을 때, 두 벡터의 점곱이 반환된다.",
			fields: Vector3D{X: 1, Y: 2, Z: 3},
			arg:    Vector3D{X: 4, Y: 5, Z: 6},
			want:   4 + 10 + 18,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Vector3D{
				X: tt.fields.X,
				Y: tt.fields.Y,
				Z: tt.fields.Z,
			}
			if got := v.Dot(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector3D_Length(t *testing.T) {
	tests := []struct {
		name   string
		fields Vector3D
		want   float64
	}{
		{
			name:   "벡터의 길이가 계산된다.",
			fields: Vector3D{X: 4, Y: 5, Z: 6},
			want:   math.Sqrt(4*4 + 5*5 + 6*6),
		},
		{
			name:   "벡터의 성분이 음수인 경우에, 길이가 계산된다.",
			fields: Vector3D{X: -4, Y: -5, Z: -6},
			want:   math.Sqrt(4*4 + 5*5 + 6*6),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.Length(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}
