package enemy

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Constants for enemy behavior
const (
	ENEMY_HEIGHT = 2.0
	ENEMY_RADIUS = 0.5
	HIT_FLASH_DURATION = 0.2
)

// Enemy represents an enemy entity
type Enemy struct {
	Position    rl.Vector3
	Health      float32
	HitTimer    float32
	Radius      float32
	Height      float32
}

// New creates a new enemy with default values
func New() *Enemy {
	return &Enemy{
		Position: rl.Vector3{X: 8, Y: 1, Z: 8},
		Health:   100.0,
		HitTimer: 0.0,
		Radius:   ENEMY_RADIUS,
		Height:   ENEMY_HEIGHT,
	}
}

// Update updates the enemy state
func (e *Enemy) Update(deltaTime float32) {
	// Update hit timer
	if e.HitTimer > 0 {
		e.HitTimer -= deltaTime
	}
	
	// Simple circular movement pattern
	time := float32(rl.GetTime())
	radius := float32(3.0)
	speed := float32(0.5)
	
	// Move in a circle around the center
	e.Position.X = float32(math.Cos(float64(time*speed))) * radius
	e.Position.Z = float32(math.Sin(float64(time*speed))) * radius
	e.Position.Y = 1.0 // Keep enemy at ground level
}

// TakeDamage applies damage to the enemy and starts hit effect
func (e *Enemy) TakeDamage(damage float32) {
	e.Health -= damage
	e.HitTimer = HIT_FLASH_DURATION
}

// IsAlive returns true if the enemy has health remaining
func (e *Enemy) IsAlive() bool {
	return e.Health > 0
}

// GetBoundingBox returns the enemy's collision bounding box
func (e *Enemy) GetBoundingBox() rl.BoundingBox {
	return rl.BoundingBox{
		Min: rl.Vector3{
			X: e.Position.X - e.Radius,
			Y: e.Position.Y,
			Z: e.Position.Z - e.Radius,
		},
		Max: rl.Vector3{
			X: e.Position.X + e.Radius,
			Y: e.Position.Y + e.Height,
			Z: e.Position.Z + e.Radius,
		},
	}
} 