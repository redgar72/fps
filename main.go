// FPS Camera Demo with Perfect Mouse Control using Raylib
package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Game constants
const (
	WINDOW_WIDTH        = 800
	WINDOW_HEIGHT       = 600
	TARGET_FPS          = 120
	MOVE_SPEED          = 5.0
	MOUSE_SENSITIVITY   = 0.003
	EYE_HEIGHT          = 1.7
	HIT_FLASH_DURATION  = 0.2
	CROSSHAIR_SIZE      = 10
	CROSSHAIR_THICKNESS = 2
	TRACER_DURATION     = 0.5 // Tracer lasts 0.5 seconds
	MAX_TRACERS         = 50  // Maximum number of active tracers

	// Gun positioning constants
	GUN_OFFSET_X      = 0.5  // Gun offset to the right
	GUN_OFFSET_Y      = -0.3 // Gun offset down
	GUN_OFFSET_Z      = 1.5  // Gun offset forward (away from camera)
	GUN_BARREL_LENGTH = 0.8  // Length of gun barrel
)

// Player represents the player state
type Player struct {
	Position rl.Vector3
	Yaw      float32
	Pitch    float32
}

// Tracer represents a bullet tracer line
type Tracer struct {
	Start    rl.Vector3
	End      rl.Vector3
	TimeLeft float32
	IsActive bool
}

// GameState holds all the game state
type GameState struct {
	Player         Player
	Camera         rl.Camera3D
	WeaponCamera   rl.Camera3D // Separate camera for weapon rendering
	WeaponModel    rl.Model    // 3D weapon model
	Cubes          []rl.Vector3
	Colors         []rl.Color
	OriginalColors []rl.Color
	HitTimers      []float32
	Tracers        []Tracer
}

// initializeGame sets up the initial game state
func initializeGame() *GameState {
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
	weaponModel := rl.LoadModel("ak47.glb")

	// Check if model loaded successfully
	if weaponModel.MeshCount == 0 {
		// Model failed to load - create a simple fallback cube
		weaponModel = rl.LoadModelFromMesh(rl.GenMeshCube(1.0, 1.0, 3.0))
	}

	return &GameState{
		Player: Player{
			Position: rl.Vector3{X: 0, Y: 1, Z: 10},
			Yaw:      0,
			Pitch:    0,
		},
		Camera: rl.Camera3D{
			Position:   rl.Vector3{X: 0, Y: 2, Z: 10},
			Target:     rl.Vector3{X: 0, Y: 2, Z: 0},
			Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
			Fovy:       60,
			Projection: rl.CameraPerspective,
		},
		WeaponCamera: rl.Camera3D{
			Position:   rl.Vector3{X: 0, Y: 0, Z: 5}, // Moved camera much further back
			Target:     rl.Vector3{X: 0, Y: 0, Z: 0}, // Looking at origin
			Up:         rl.Vector3{X: 0, Y: 1, Z: 0}, // Up vector
			Fovy:       60,                           // Standard FOV
			Projection: rl.CameraPerspective,
		},
		WeaponModel:    weaponModel,
		Cubes:          cubes,
		Colors:         colors,
		OriginalColors: originalColors,
		HitTimers:      make([]float32, len(cubes)),
		Tracers:        make([]Tracer, MAX_TRACERS),
	}
}

// handleMouseLook processes mouse movement for camera rotation
func handleMouseLook(player *Player, deltaTime float32) {
	mouseDelta := rl.GetMouseDelta()

	// Update rotation based on mouse movement
	player.Yaw -= mouseDelta.X * MOUSE_SENSITIVITY
	player.Pitch -= mouseDelta.Y * MOUSE_SENSITIVITY

	// Clamp pitch to prevent over-rotation
	maxPitch := float32(math.Pi/2 - 0.1)
	if player.Pitch > maxPitch {
		player.Pitch = maxPitch
	} else if player.Pitch < -maxPitch {
		player.Pitch = -maxPitch
	}
}

// handleMovement processes WASD movement input
func handleMovement(player *Player, deltaTime float32) {
	moveSpeed := MOVE_SPEED * deltaTime

	// Calculate forward and right vectors based on yaw
	forward := rl.Vector3{
		X: float32(math.Sin(float64(player.Yaw))),
		Y: 0,
		Z: float32(math.Cos(float64(player.Yaw))),
	}
	right := rl.Vector3{
		X: -float32(math.Cos(float64(player.Yaw))),
		Y: 0,
		Z: float32(math.Sin(float64(player.Yaw))),
	}

	// Movement input
	if rl.IsKeyDown(rl.KeyW) {
		player.Position = rl.Vector3Add(player.Position, rl.Vector3Scale(forward, moveSpeed))
	}
	if rl.IsKeyDown(rl.KeyS) {
		player.Position = rl.Vector3Subtract(player.Position, rl.Vector3Scale(forward, moveSpeed))
	}
	if rl.IsKeyDown(rl.KeyD) {
		player.Position = rl.Vector3Add(player.Position, rl.Vector3Scale(right, moveSpeed))
	}
	if rl.IsKeyDown(rl.KeyA) {
		player.Position = rl.Vector3Subtract(player.Position, rl.Vector3Scale(right, moveSpeed))
	}
}

// updateCamera updates the camera position and target based on player state
func updateCamera(game *GameState) {
	// Update camera position to player position (+ eye height)
	game.Camera.Position = rl.Vector3{
		X: game.Player.Position.X,
		Y: game.Player.Position.Y + EYE_HEIGHT,
		Z: game.Player.Position.Z,
	}

	// Update camera target based on yaw and pitch
	targetDistance := float32(1.0)
	game.Camera.Target = rl.Vector3{
		X: game.Camera.Position.X + float32(math.Sin(float64(game.Player.Yaw))*math.Cos(float64(game.Player.Pitch)))*targetDistance,
		Y: game.Camera.Position.Y + float32(math.Sin(float64(game.Player.Pitch)))*targetDistance,
		Z: game.Camera.Position.Z + float32(math.Cos(float64(game.Player.Yaw))*math.Cos(float64(game.Player.Pitch)))*targetDistance,
	}
}

// handleSystemInput processes system-level input (cursor toggle, exit)
func handleSystemInput() bool {
	// Toggle cursor capture with Tab
	if rl.IsKeyPressed(rl.KeyTab) {
		if rl.IsCursorHidden() {
			rl.EnableCursor()
		} else {
			rl.DisableCursor()
		}
	}

	// Exit with Escape
	return rl.IsKeyPressed(rl.KeyEscape)
}

// updateHitTimers decrements hit timers and resets colors when they expire
func updateHitTimers(game *GameState, deltaTime float32) {
	for i := range game.HitTimers {
		if game.HitTimers[i] > 0 {
			game.HitTimers[i] -= deltaTime
			if game.HitTimers[i] <= 0 {
				// Reset to original color when timer expires
				game.Colors[i] = game.OriginalColors[i]
			}
		}
	}
}

// updateTracers decrements tracer timers and deactivates expired tracers
func updateTracers(game *GameState, deltaTime float32) {
	for i := range game.Tracers {
		if game.Tracers[i].IsActive {
			game.Tracers[i].TimeLeft -= deltaTime
			if game.Tracers[i].TimeLeft <= 0 {
				game.Tracers[i].IsActive = false
			}
		}
	}
}

// addTracer creates a new tracer line from start to end point
func addTracer(game *GameState, start, end rl.Vector3) {
	// Find first inactive tracer slot
	for i := range game.Tracers {
		if !game.Tracers[i].IsActive {
			game.Tracers[i] = Tracer{
				Start:    start,
				End:      end,
				TimeLeft: TRACER_DURATION,
				IsActive: true,
			}
			return
		}
	}
	// If no free slots, overwrite the oldest (first) tracer
	game.Tracers[0] = Tracer{
		Start:    start,
		End:      end,
		TimeLeft: TRACER_DURATION,
		IsActive: true,
	}
}

// handleShooting processes shooting input and raycast collision
func handleShooting(game *GameState) {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && rl.IsCursorHidden() {
		// Create ray from camera position in the direction camera is looking (for accurate hit detection)
		rayOrigin := game.Camera.Position
		rayDirection := rl.Vector3Normalize(rl.Vector3Subtract(game.Camera.Target, game.Camera.Position))

		// Get gun barrel tip for visual tracer start (world position)
		gunBarrelTip := calculateGunBarrelTip(game)

		// Default end point for tracer (if no hit, draw a long line)
		maxRange := float32(100.0)
		tracerEnd := rl.Vector3Add(rayOrigin, rl.Vector3Scale(rayDirection, maxRange))

		// Check collision with each cube using camera raycast (accurate)
		for i, cubePos := range game.Cubes {
			cubeSize := rl.Vector3{X: 1, Y: 1, Z: 1}
			boundingBox := rl.BoundingBox{
				Min: rl.Vector3Subtract(cubePos, rl.Vector3Scale(cubeSize, 0.5)),
				Max: rl.Vector3Add(cubePos, rl.Vector3Scale(cubeSize, 0.5)),
			}

			// Perform ray-box intersection test
			collision := rl.GetRayCollisionBox(rl.Ray{Position: rayOrigin, Direction: rayDirection}, boundingBox)
			if collision.Hit {
				// Hit detected! Flash white and start timer
				game.Colors[i] = rl.White
				game.HitTimers[i] = HIT_FLASH_DURATION

				// Use actual hit point for tracer end
				tracerEnd = collision.Point
				break // Only hit the first cube in the ray path
			}
		}

		// Create tracer from gun barrel to hit point (visual effect)
		addTracer(game, gunBarrelTip, tracerEnd)
	}
}

// renderWorld draws the 3D world elements
func renderWorld(game *GameState) {
	rl.BeginMode3D(game.Camera)

	// Draw ground plane
	rl.DrawPlane(rl.Vector3{X: 0, Y: 0, Z: 0}, rl.Vector2{X: 20, Y: 20}, rl.Green)

	// Draw cubes
	for i, pos := range game.Cubes {
		rl.DrawCube(pos, 1, 1, 1, game.Colors[i%len(game.Colors)])
		rl.DrawCubeWires(pos, 1, 1, 1, rl.Black)
	}

	// Draw player representation (semi-transparent box)
	rl.DrawCube(game.Player.Position, 0.5, EYE_HEIGHT, 0.5, rl.NewColor(255, 0, 0, 100))

	// Draw tracers
	renderTracers(game)

	rl.EndMode3D()

	// Draw weapon in separate viewport (after ending world 3D mode)
	renderWeaponViewport(game)
}

// calculateGunViewPosition returns the gun position relative to camera (viewport space)
func calculateGunViewPosition() rl.Vector3 {
	// Fixed position relative to camera viewport (right, down, forward)
	return rl.Vector3{
		X: GUN_OFFSET_X, // Right offset
		Y: GUN_OFFSET_Y, // Down offset
		Z: GUN_OFFSET_Z, // Forward offset
	}
}

// calculateGunWorldPosition converts viewport position to world position for tracers
func calculateGunWorldPosition(game *GameState) rl.Vector3 {
	// Get camera orientation vectors
	yaw := game.Player.Yaw
	pitch := game.Player.Pitch

	// Right vector (for X offset)
	right := rl.Vector3{
		X: -float32(math.Cos(float64(yaw))),
		Y: 0,
		Z: float32(math.Sin(float64(yaw))),
	}

	// Up vector (for Y offset)
	up := rl.Vector3{X: 0, Y: 1, Z: 0}

	// Forward vector (for Z offset)
	forward := rl.Vector3{
		X: float32(math.Sin(float64(yaw)) * math.Cos(float64(pitch))),
		Y: float32(math.Sin(float64(pitch))),
		Z: float32(math.Cos(float64(yaw)) * math.Cos(float64(pitch))),
	}

	// Start from camera position
	gunPos := game.Camera.Position

	// Apply offsets in world space
	gunPos = rl.Vector3Add(gunPos, rl.Vector3Scale(right, GUN_OFFSET_X))
	gunPos = rl.Vector3Add(gunPos, rl.Vector3Scale(up, GUN_OFFSET_Y))
	gunPos = rl.Vector3Add(gunPos, rl.Vector3Scale(forward, GUN_OFFSET_Z))

	return gunPos
}

// calculateGunBarrelTip returns the world position of the gun barrel tip for tracers
func calculateGunBarrelTip(game *GameState) rl.Vector3 {
	gunPos := calculateGunWorldPosition(game)

	// Forward vector for barrel extension
	yaw := game.Player.Yaw
	pitch := game.Player.Pitch
	forward := rl.Vector3{
		X: float32(math.Sin(float64(yaw)) * math.Cos(float64(pitch))),
		Y: float32(math.Sin(float64(pitch))),
		Z: float32(math.Cos(float64(yaw)) * math.Cos(float64(pitch))),
	}

	// Extend gun position by barrel length
	return rl.Vector3Add(gunPos, rl.Vector3Scale(forward, GUN_BARREL_LENGTH))
}

// renderWeapon draws a 3D weapon model in the main world (for testing)
func renderWeapon(game *GameState) {
	// Calculate camera-relative position that rotates with camera
	yaw := game.Player.Yaw

	// Right vector (for positioning weapon to the right)
	right := rl.Vector3{
		X: -float32(math.Cos(float64(yaw))),
		Y: 0,
		Z: float32(math.Sin(float64(yaw))),
	}

	// Forward vector (for positioning weapon forward)
	forward := rl.Vector3{
		X: float32(math.Sin(float64(yaw))),
		Y: 0,
		Z: float32(math.Cos(float64(yaw))),
	}

	// Position weapon relative to camera with proper rotation
	weaponOffset := rl.Vector3{X: 0, Y: 0, Z: 0}
	weaponOffset = rl.Vector3Add(weaponOffset, rl.Vector3Scale(right, 0.15))    // Closer to center (was 0.3)
	weaponOffset = rl.Vector3Add(weaponOffset, rl.Vector3{X: 0, Y: -0.2, Z: 0}) // Less down (was -0.3)
	weaponOffset = rl.Vector3Add(weaponOffset, rl.Vector3Scale(forward, 0.25))  // Much closer forward (was 0.8)

	weaponWorldPos := rl.Vector3Add(game.Camera.Position, weaponOffset)
	weaponScale := float32(0.2) // Smaller scale

	// Remove test cube, just draw weapon
	rl.DrawModelWires(game.WeaponModel, weaponWorldPos, weaponScale, rl.Yellow)

	// Optional: Draw solid model too to see if it works
	rl.DrawModel(game.WeaponModel, weaponWorldPos, weaponScale, rl.Gray)
}

// renderWeaponViewport draws the weapon in a separate static viewport
func renderWeaponViewport(game *GameState) {
	// Set up weapon viewport (bottom-right corner)
	screenWidth := rl.GetScreenWidth()
	screenHeight := rl.GetScreenHeight()
	weaponViewportWidth := screenWidth / 3
	weaponViewportHeight := screenHeight / 3
	weaponViewportX := screenWidth - weaponViewportWidth
	weaponViewportY := screenHeight - weaponViewportHeight

	// Clear the weapon viewport area with a transparent background
	rl.DrawRectangle(int32(weaponViewportX), int32(weaponViewportY), int32(weaponViewportWidth), int32(weaponViewportHeight), rl.ColorAlpha(rl.Black, 0.1))

	// Begin scissor mode to clip rendering to the weapon viewport
	rl.BeginScissorMode(int32(weaponViewportX), int32(weaponViewportY), int32(weaponViewportWidth), int32(weaponViewportHeight))

	// Begin weapon camera 3D mode
	rl.BeginMode3D(game.WeaponCamera)

	// Position weapon at center of weapon viewport
	weaponPos := rl.Vector3{X: 0.0, Y: 0.0, Z: 0.0}
	weaponScale := float32(0.5)

	// Draw the weapon model
	rl.DrawModel(game.WeaponModel, weaponPos, weaponScale, rl.White)
	rl.DrawModelWires(game.WeaponModel, weaponPos, weaponScale, rl.Black)

	rl.EndMode3D()

	// End scissor mode
	rl.EndScissorMode()
}

// renderTracers draws all active tracer lines with fade effect
func renderTracers(game *GameState) {
	for _, tracer := range game.Tracers {
		if tracer.IsActive {
			// Calculate fade alpha based on remaining time
			fadeRatio := tracer.TimeLeft / TRACER_DURATION
			alpha := uint8(fadeRatio * 255)

			// Create color with fade effect (bright yellow/orange tracer)
			tracerColor := rl.NewColor(255, 255, 0, alpha) // Yellow with fade

			// Draw the tracer line
			rl.DrawLine3D(tracer.Start, tracer.End, tracerColor)
		}
	}
}

// renderUI draws all UI elements
func renderUI() {
	// Title and controls
	rl.DrawText("FPS Camera with Perfect Mouse Control!", 10, 10, 20, rl.DarkGray)
	rl.DrawText("WASD: Move | Mouse: Look | Left Click: Shoot | Tab: Toggle cursor | ESC: Exit", 10, 35, 16, rl.DarkGray)

	// Cursor status
	if rl.IsCursorHidden() {
		rl.DrawText("Mouse captured - infinite rotation!", 10, 55, 16, rl.Green)
	} else {
		rl.DrawText("Mouse free - press Tab to capture", 10, 55, 16, rl.Red)
	}

	// FPS counter in top right
	fps := rl.GetFPS()
	fpsText := fmt.Sprintf("FPS: %d", fps)
	fpsTextWidth := rl.MeasureText(fpsText, 20)
	screenWidth := int32(rl.GetScreenWidth())
	rl.DrawText(fpsText, screenWidth-fpsTextWidth-10, 10, 20, rl.White)
}

// renderCrosshair draws the crosshair when in FPS mode
func renderCrosshair() {
	if rl.IsCursorHidden() {
		centerX := int32(rl.GetScreenWidth() / 2)
		centerY := int32(rl.GetScreenHeight() / 2)

		// Draw horizontal line
		rl.DrawRectangle(centerX-CROSSHAIR_SIZE, centerY-CROSSHAIR_THICKNESS/2, CROSSHAIR_SIZE*2, CROSSHAIR_THICKNESS, rl.White)
		// Draw vertical line
		rl.DrawRectangle(centerX-CROSSHAIR_THICKNESS/2, centerY-CROSSHAIR_SIZE, CROSSHAIR_THICKNESS, CROSSHAIR_SIZE*2, rl.White)
	}
}

func main() {
	// Initialize window
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "FPS Camera with Perfect Mouse Control - Raylib")
	defer rl.CloseWindow()

	// Initialize game state
	game := initializeGame()

	// Enable cursor capture for FPS controls
	rl.DisableCursor()
	rl.SetTargetFPS(TARGET_FPS)

	// Main game loop
	for !rl.WindowShouldClose() {
		deltaTime := rl.GetFrameTime()

		// Handle all input
		handleMouseLook(&game.Player, deltaTime)
		handleMovement(&game.Player, deltaTime)
		if handleSystemInput() {
			break // Exit requested
		}

		// Update game systems
		updateCamera(game)
		updateHitTimers(game, deltaTime)
		updateTracers(game, deltaTime)
		handleShooting(game)

		// Render everything
		rl.BeginDrawing()
		rl.ClearBackground(rl.NewColor(135, 206, 235, 255)) // Sky blue

		renderWorld(game)
		renderUI()
		renderCrosshair()

		rl.EndDrawing()
	}
}
