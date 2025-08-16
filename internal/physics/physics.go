package physics

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"fps/internal/enemy"
)

// Constants for physics
const (
	HIT_FLASH_DURATION = 0.2
	TRACER_DURATION    = 0.5
	MAX_TRACERS        = 50
)

// Tracer represents a bullet tracer line
type Tracer struct {
	Start    rl.Vector3
	End      rl.Vector3
	TimeLeft float32
	IsActive bool
}

// TracerManager manages all active tracers
type TracerManager struct {
	Tracers []Tracer
}

// NewTracerManager creates a new tracer manager
func NewTracerManager() *TracerManager {
	return &TracerManager{
		Tracers: make([]Tracer, MAX_TRACERS),
	}
}

// Update updates all active tracers
func (tm *TracerManager) Update(deltaTime float32) {
	for i := range tm.Tracers {
		if tm.Tracers[i].IsActive {
			tm.Tracers[i].TimeLeft -= deltaTime
			if tm.Tracers[i].TimeLeft <= 0 {
				tm.Tracers[i].IsActive = false
			}
		}
	}
}

// AddTracer creates a new tracer line from start to end point
func (tm *TracerManager) AddTracer(start, end rl.Vector3) {
	// Find first inactive tracer slot
	for i := range tm.Tracers {
		if !tm.Tracers[i].IsActive {
			tm.Tracers[i] = Tracer{
				Start:    start,
				End:      end,
				TimeLeft: TRACER_DURATION,
				IsActive: true,
			}
			return
		}
	}
	// If no free slots, overwrite the oldest (first) tracer
	tm.Tracers[0] = Tracer{
		Start:    start,
		End:      end,
		TimeLeft: TRACER_DURATION,
		IsActive: true,
	}
}

// GetActiveTracers returns all active tracers
func (tm *TracerManager) GetActiveTracers() []Tracer {
	return tm.Tracers
}

// CheckCubeCollision checks if a ray hits any cube and returns hit info
func CheckCubeCollision(rayOrigin, rayDirection rl.Vector3, cubes []rl.Vector3, colors []rl.Color, hitTimers []float32) (bool, rl.Vector3, int) {
	// Check collision with each cube using camera raycast
	for i, cubePos := range cubes {
		cubeSize := rl.Vector3{X: 1, Y: 1, Z: 1}
		boundingBox := rl.BoundingBox{
			Min: rl.Vector3Subtract(cubePos, rl.Vector3Scale(cubeSize, 0.5)),
			Max: rl.Vector3Add(cubePos, rl.Vector3Scale(cubeSize, 0.5)),
		}

		// Perform ray-box intersection test
		collision := rl.GetRayCollisionBox(rl.Ray{Position: rayOrigin, Direction: rayDirection}, boundingBox)
		if collision.Hit {
			return true, collision.Point, i
		}
	}
	return false, rl.Vector3{}, -1
}

// CheckEnemyCollision checks if a ray hits the enemy and returns hit info
func CheckEnemyCollision(rayOrigin, rayDirection rl.Vector3, e *enemy.Enemy) (bool, rl.Vector3) {
	enemyBoundingBox := e.GetBoundingBox()
	enemyCollision := rl.GetRayCollisionBox(rl.Ray{Position: rayOrigin, Direction: rayDirection}, enemyBoundingBox)
	if enemyCollision.Hit {
		return true, enemyCollision.Point
	}
	return false, rl.Vector3{}
} 