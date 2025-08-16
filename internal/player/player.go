package player

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Constants for player behavior
const (
	EYE_HEIGHT = 1.7
)

// Player represents the player state
type Player struct {
	Position rl.Vector3
	Yaw      float32
	Pitch    float32
}

// New creates a new player with default values
func New() *Player {
	return &Player{
		Position: rl.Vector3{X: 0, Y: 1, Z: 10},
		Yaw:      0,
		Pitch:    0,
	}
}

// GetEyePosition returns the camera position (player position + eye height)
func (p *Player) GetEyePosition() rl.Vector3 {
	return rl.Vector3{
		X: p.Position.X,
		Y: p.Position.Y + EYE_HEIGHT,
		Z: p.Position.Z,
	}
}

// GetForwardVector returns the forward direction vector based on yaw and pitch
func (p *Player) GetForwardVector() rl.Vector3 {
	return rl.Vector3{
		X: float32(math.Sin(float64(p.Yaw)) * math.Cos(float64(p.Pitch))),
		Y: float32(math.Sin(float64(p.Pitch))),
		Z: float32(math.Cos(float64(p.Yaw)) * math.Cos(float64(p.Pitch))),
	}
}

// GetRightVector returns the right direction vector based on yaw
func (p *Player) GetRightVector() rl.Vector3 {
	return rl.Vector3{
		X: -float32(math.Cos(float64(p.Yaw))),
		Y: 0,
		Z: float32(math.Sin(float64(p.Yaw))),
	}
}

// GetUpVector returns the up direction vector
func (p *Player) GetUpVector() rl.Vector3 {
	return rl.Vector3{X: 0, Y: 1, Z: 0}
} 