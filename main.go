// FPS Camera Demo with Perfect Mouse Control using Raylib
package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// Initialize window
	rl.InitWindow(800, 600, "FPS Camera with Perfect Mouse Control - Raylib")
	defer rl.CloseWindow()

	// Set up FPS camera
	camera := rl.Camera3D{
		Position:   rl.Vector3{X: 0, Y: 2, Z: 10}, // Eye position
		Target:     rl.Vector3{X: 0, Y: 2, Z: 0},  // Target to look at
		Up:         rl.Vector3{X: 0, Y: 1, Z: 0},  // Up vector
		Fovy:       60,                            // Field of view
		Projection: rl.CameraPerspective,          // Perspective projection
	}

	// Enable cursor capture for FPS controls (this is what G3N is missing!)
	rl.DisableCursor()

	// Create a simple player object (conceptually - the camera represents the player)
	playerPos := rl.Vector3{X: 0, Y: 1, Z: 10}
	yaw := float32(0)
	pitch := float32(0)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		deltaTime := rl.GetFrameTime()

		// Get mouse movement (infinite because cursor is captured!)
		mouseDelta := rl.GetMouseDelta()
		sensitivity := float32(0.003)

		// Update rotation based on mouse movement
		yaw -= mouseDelta.X * sensitivity   // Horizontal rotation
		pitch -= mouseDelta.Y * sensitivity // Vertical rotation

		// Clamp pitch to prevent over-rotation
		maxPitch := float32(math.Pi/2 - 0.1)
		if pitch > maxPitch {
			pitch = maxPitch
		} else if pitch < -maxPitch {
			pitch = -maxPitch
		}

		// Handle WASD movement
		moveSpeed := float32(5.0) * deltaTime

		// Calculate forward and right vectors based on yaw
		forward := rl.Vector3{
			X: float32(math.Sin(float64(yaw))),
			Y: 0,
			Z: float32(math.Cos(float64(yaw))),
		}
		right := rl.Vector3{
			X: -float32(math.Cos(float64(yaw))),
			Y: 0,
			Z: float32(math.Sin(float64(yaw))),
		}

		// Movement input
		if rl.IsKeyDown(rl.KeyW) {
			playerPos = rl.Vector3Add(playerPos, rl.Vector3Scale(forward, moveSpeed))
		}
		if rl.IsKeyDown(rl.KeyS) {
			playerPos = rl.Vector3Subtract(playerPos, rl.Vector3Scale(forward, moveSpeed))
		}
		if rl.IsKeyDown(rl.KeyD) {
			playerPos = rl.Vector3Add(playerPos, rl.Vector3Scale(right, moveSpeed))
		}
		if rl.IsKeyDown(rl.KeyA) {
			playerPos = rl.Vector3Subtract(playerPos, rl.Vector3Scale(right, moveSpeed))
		}

		// Update camera position to player position (+ eye height)
		eyeHeight := float32(1.7)
		camera.Position = rl.Vector3{X: playerPos.X, Y: playerPos.Y + eyeHeight, Z: playerPos.Z}

		// Update camera target based on yaw and pitch
		// This is the key: camera looks in the direction determined by yaw and pitch
		targetDistance := float32(1.0) // Distance to look ahead
		camera.Target = rl.Vector3{
			X: camera.Position.X + float32(math.Sin(float64(yaw))*math.Cos(float64(pitch)))*targetDistance,
			Y: camera.Position.Y + float32(math.Sin(float64(pitch)))*targetDistance,
			Z: camera.Position.Z + float32(math.Cos(float64(yaw))*math.Cos(float64(pitch)))*targetDistance,
		}

		// Toggle cursor capture with Tab
		if rl.IsKeyPressed(rl.KeyTab) {
			if rl.IsCursorHidden() {
				rl.EnableCursor()
			} else {
				rl.DisableCursor()
			}
		}

		// Exit with Escape
		if rl.IsKeyPressed(rl.KeyEscape) {
			break
		}

		// Rendering
		rl.BeginDrawing()
		rl.ClearBackground(rl.NewColor(135, 206, 235, 255)) // Sky blue

		rl.BeginMode3D(camera)

		// Draw ground plane
		rl.DrawPlane(rl.Vector3{X: 0, Y: 0, Z: 0}, rl.Vector2{X: 20, Y: 20}, rl.Green)

		// Draw some reference cubes
		cubes := []rl.Vector3{
			{X: 5, Y: 0.5, Z: 5},
			{X: -5, Y: 0.5, Z: 5},
			{X: 5, Y: 0.5, Z: -5},
			{X: -5, Y: 0.5, Z: -5},
			{X: 0, Y: 0.5, Z: 10},
		}
		colors := []rl.Color{rl.Red, rl.Blue, rl.Yellow, rl.Purple, rl.Orange}

		for i, pos := range cubes {
			rl.DrawCube(pos, 1, 1, 1, colors[i%len(colors)])
			rl.DrawCubeWires(pos, 1, 1, 1, rl.Black)
		}

		// Draw a visual representation of the player (semi-transparent box)
		rl.DrawCube(playerPos, 0.5, eyeHeight, 0.5, rl.NewColor(255, 0, 0, 100))

		rl.EndMode3D()

		// UI
		rl.DrawText("FPS Camera with Perfect Mouse Control!", 10, 10, 20, rl.DarkGray)
		rl.DrawText("WASD: Move | Mouse: Look | Tab: Toggle cursor | ESC: Exit", 10, 35, 16, rl.DarkGray)

		if rl.IsCursorHidden() {
			rl.DrawText("Mouse captured - infinite rotation!", 10, 55, 16, rl.Green)
		} else {
			rl.DrawText("Mouse free - press Tab to capture", 10, 55, 16, rl.Red)
		}

		rl.EndDrawing()
	}
}
