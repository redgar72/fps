package input

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"fps/internal/player"
)

// Constants for input behavior
const (
	MOVE_SPEED        = 5.0
	MOUSE_SENSITIVITY = 0.003
)

// HandleMouseLook processes mouse movement for camera rotation
func HandleMouseLook(p *player.Player, deltaTime float32) {
	mouseDelta := rl.GetMouseDelta()

	// Update rotation based on mouse movement
	p.Yaw -= mouseDelta.X * MOUSE_SENSITIVITY
	p.Pitch -= mouseDelta.Y * MOUSE_SENSITIVITY

	// Clamp pitch to prevent over-rotation
	maxPitch := float32(math.Pi/2 - 0.1)
	if p.Pitch > maxPitch {
		p.Pitch = maxPitch
	} else if p.Pitch < -maxPitch {
		p.Pitch = -maxPitch
	}
}

// HandleMovement processes WASD movement input
func HandleMovement(p *player.Player, deltaTime float32) {
	moveSpeed := MOVE_SPEED * deltaTime

	// Get movement vectors
	forward := p.GetForwardVector()
	right := p.GetRightVector()

	// Movement input
	if rl.IsKeyDown(rl.KeyW) {
		p.Position = rl.Vector3Add(p.Position, rl.Vector3Scale(forward, moveSpeed))
	}
	if rl.IsKeyDown(rl.KeyS) {
		p.Position = rl.Vector3Subtract(p.Position, rl.Vector3Scale(forward, moveSpeed))
	}
	if rl.IsKeyDown(rl.KeyD) {
		p.Position = rl.Vector3Add(p.Position, rl.Vector3Scale(right, moveSpeed))
	}
	if rl.IsKeyDown(rl.KeyA) {
		p.Position = rl.Vector3Subtract(p.Position, rl.Vector3Scale(right, moveSpeed))
	}

	// Keep player anchored to the ground plane (Y = 1)
	p.Position.Y = 1.0

	// Keep player within world boundaries
	worldSize := float32(10.0)
	if p.Position.X > worldSize {
		p.Position.X = worldSize
	}
	if p.Position.X < -worldSize {
		p.Position.X = -worldSize
	}
	if p.Position.Z > worldSize {
		p.Position.Z = worldSize
	}
	if p.Position.Z < -worldSize {
		p.Position.Z = -worldSize
	}
}

// HandleSystemInput processes system-level input (cursor toggle, exit)
func HandleSystemInput() bool {
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