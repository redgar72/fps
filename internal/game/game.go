package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"fps/internal/player"
	"fps/internal/enemy"
	"fps/internal/physics"
)

// GameState holds all the game state
type GameState struct {
	Player         *player.Player
	Camera         rl.Camera3D
	WeaponCamera   rl.Camera3D
	WeaponModel    rl.Model
	Cubes          []rl.Vector3
	Colors         []rl.Color
	OriginalColors []rl.Color
	HitTimers      []float32
	TracerManager  *physics.TracerManager
	Enemy          *enemy.Enemy
}

// New creates a new game state
func New() *GameState {
	cubes := []rl.Vector3{
		{X: 5, Y: 0.5, Z: 5},
		{X: -5, Y: 0.5, Z: 5},
		{X: 5, Y: 0.5, Z: -5},
		{X: -5, Y: 0.5, Z: -5},
		{X: 0, Y: 0.5, Z: 10},
	}

	originalColors := []rl.Color{rl.Red, rl.Blue, rl.Yellow, rl.Purple, rl.Orange}
	colors := make([]rl.Color, len(originalColors))
	copy(colors, originalColors)

	// Load weapon model
	weaponModel := rl.LoadModel("assets/ak47.glb")

	// Check if model loaded successfully
	if weaponModel.MeshCount == 0 {
		// Model failed to load - create a simple fallback cube
		weaponModel = rl.LoadModelFromMesh(rl.GenMeshCube(1.0, 1.0, 3.0))
	}

	return &GameState{
		Player:         player.New(),
		Camera:         rl.Camera3D{
			Position:   rl.Vector3{X: 0, Y: 2, Z: 10},
			Target:     rl.Vector3{X: 0, Y: 2, Z: 0},
			Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
			Fovy:       60,
			Projection: rl.CameraPerspective,
		},
		WeaponCamera:   rl.Camera3D{
			Position:   rl.Vector3{X: 0, Y: 0, Z: 5},
			Target:     rl.Vector3{X: 0, Y: 0, Z: 0},
			Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
			Fovy:       60,
			Projection: rl.CameraPerspective,
		},
		WeaponModel:    weaponModel,
		Cubes:          cubes,
		Colors:         colors,
		OriginalColors: originalColors,
		HitTimers:      make([]float32, len(cubes)),
		TracerManager:  physics.NewTracerManager(),
		Enemy:          enemy.New(),
	}
}

// UpdateCamera updates the camera position and target based on player state
func (g *GameState) UpdateCamera() {
	// Update camera position to player eye position
	g.Camera.Position = g.Player.GetEyePosition()

	// Update camera target based on yaw and pitch
	targetDistance := float32(1.0)
	forward := g.Player.GetForwardVector()
	g.Camera.Target = rl.Vector3Add(g.Camera.Position, rl.Vector3Scale(forward, targetDistance))
}

// UpdateHitTimers decrements hit timers and resets colors when they expire
func (g *GameState) UpdateHitTimers(deltaTime float32) {
	for i := range g.HitTimers {
		if g.HitTimers[i] > 0 {
			g.HitTimers[i] -= deltaTime
			if g.HitTimers[i] <= 0 {
				// Reset to original color when timer expires
				g.Colors[i] = g.OriginalColors[i]
			}
		}
	}
}

// Update updates all game systems
func (g *GameState) Update(deltaTime float32) {
	g.UpdateCamera()
	g.UpdateHitTimers(deltaTime)
	g.TracerManager.Update(deltaTime)
	g.Enemy.Update(deltaTime)
}

// HandleShooting processes shooting input and raycast collision
func (g *GameState) HandleShooting() {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && rl.IsCursorHidden() {
		// Create ray from camera position in the direction camera is looking
		rayOrigin := g.Camera.Position
		rayDirection := rl.Vector3Normalize(rl.Vector3Subtract(g.Camera.Target, g.Camera.Position))

		// Default end point for tracer (if no hit, draw a long line)
		maxRange := float32(100.0)
		tracerEnd := rl.Vector3Add(rayOrigin, rl.Vector3Scale(rayDirection, maxRange))

		// Check collision with cubes
		hit, hitPoint, cubeIndex := physics.CheckCubeCollision(rayOrigin, rayDirection, g.Cubes, g.Colors, g.HitTimers)
		if hit {
			// Hit detected! Flash white and start timer
			g.Colors[cubeIndex] = rl.White
			g.HitTimers[cubeIndex] = physics.HIT_FLASH_DURATION
			tracerEnd = hitPoint
		}

		// Check collision with enemy
		enemyHit, enemyHitPoint := physics.CheckEnemyCollision(rayOrigin, rayDirection, g.Enemy)
		if enemyHit {
			// Enemy hit! Apply damage and flash effect
			g.Enemy.TakeDamage(25.0)
			tracerEnd = enemyHitPoint
		}

		// Create tracer from camera to hit point
		g.TracerManager.AddTracer(rayOrigin, tracerEnd)
	}
}

// GetCubes returns the cubes slice
func (g *GameState) GetCubes() []rl.Vector3 {
	return g.Cubes
}

// GetColors returns the colors slice
func (g *GameState) GetColors() []rl.Color {
	return g.Colors
}

// GetHitTimers returns the hit timers slice
func (g *GameState) GetHitTimers() []float32 {
	return g.HitTimers
} 