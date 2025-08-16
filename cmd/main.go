// FPS Camera Demo with Perfect Mouse Control using Raylib
package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"fps/internal/game"
	"fps/internal/input"
	"fps/internal/rendering"
)

// Game constants
const (
	WINDOW_WIDTH  = 800
	WINDOW_HEIGHT = 600
	TARGET_FPS    = 120

	// Anti-aliasing configuration
	FLAG_MSAA_4X_HINT   = 0x00000020 // Enable 4x MSAA for smoother edges
	FLAG_WINDOW_HIGHDPI = 0x00002000 // Enable high DPI support for better rendering
)

func main() {
	// Set MSAA 4x hint for smoother anti-aliasing
	// This will significantly reduce the stairstepping/aliasing on cube edges
	// Also enable high DPI support for better rendering on high-resolution displays
	rl.SetConfigFlags(FLAG_MSAA_4X_HINT | FLAG_WINDOW_HIGHDPI)
	
	// Initialize window
	rl.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, "FPS Camera with Perfect Mouse Control - Raylib")
	defer rl.CloseWindow()

	// Initialize game state
	gameState := game.New()

	// Enable cursor capture for FPS controls
	rl.DisableCursor()
	rl.SetTargetFPS(TARGET_FPS)

	// Main game loop
	for !rl.WindowShouldClose() {
		deltaTime := rl.GetFrameTime()

		// Handle all input
		input.HandleMouseLook(gameState.Player, deltaTime)
		input.HandleMovement(gameState.Player, deltaTime)
		if input.HandleSystemInput() {
			break // Exit requested
		}

		// Update game systems
		gameState.Update(deltaTime)
		gameState.HandleShooting()

		// Render everything
		rl.BeginDrawing()
		rl.ClearBackground(rl.NewColor(135, 206, 235, 255)) // Sky blue

		// Render the 3D world
		rendering.RenderWorld(
			gameState.Player,
			gameState.Enemy,
			gameState.GetCubes(),
			gameState.GetColors(),
			gameState.GetHitTimers(),
			gameState.TracerManager,
			gameState.Camera,
		)

		// Render weapon viewport
		rendering.RenderWeaponViewport(gameState.WeaponCamera, gameState.WeaponModel)

		// Render UI elements
		rendering.RenderUI(gameState.Player, gameState.Enemy)
		rendering.RenderCrosshair()

		rl.EndDrawing()
	}
} 