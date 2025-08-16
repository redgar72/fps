package rendering

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"fps/internal/player"
)

// Constants for weapon rendering
const (
	GUN_OFFSET_X      = 0.5  // Gun offset to the right
	GUN_OFFSET_Y      = -0.3 // Gun offset down
	GUN_OFFSET_Z      = 1.5  // Gun offset forward (away from camera)
	GUN_BARREL_LENGTH = 0.8  // Length of gun barrel
)

// RenderWeaponViewport draws the weapon in a separate static viewport
func RenderWeaponViewport(weaponCamera rl.Camera3D, weaponModel rl.Model) {
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
	rl.BeginMode3D(weaponCamera)

	// Position weapon at center of weapon viewport
	weaponPos := rl.Vector3{X: 0.0, Y: 0.0, Z: 0.0}
	weaponScale := float32(0.5)

	// Draw the weapon model
	rl.DrawModel(weaponModel, weaponPos, weaponScale, rl.White)
	rl.DrawModelWires(weaponModel, weaponPos, weaponScale, rl.Black)

	rl.EndMode3D()

	// End scissor mode
	rl.EndScissorMode()
}

// RenderWeaponInWorld draws a 3D weapon model in the main world (for testing)
func RenderWeaponInWorld(p *player.Player, weaponModel rl.Model, cameraPos rl.Vector3) {
	// Right vector (for positioning weapon to the right)
	right := p.GetRightVector()

	// Forward vector (for positioning weapon forward)
	forward := p.GetForwardVector()

	// Position weapon relative to camera with proper rotation
	weaponOffset := rl.Vector3{X: 0, Y: 0, Z: 0}
	weaponOffset = rl.Vector3Add(weaponOffset, rl.Vector3Scale(right, 0.15))
	weaponOffset = rl.Vector3Add(weaponOffset, rl.Vector3{X: 0, Y: -0.2, Z: 0})
	weaponOffset = rl.Vector3Add(weaponOffset, rl.Vector3Scale(forward, 0.25))

	weaponWorldPos := rl.Vector3Add(cameraPos, weaponOffset)
	weaponScale := float32(0.2) // Smaller scale

	// Draw weapon
	rl.DrawModelWires(weaponModel, weaponWorldPos, weaponScale, rl.Yellow)
	rl.DrawModel(weaponModel, weaponWorldPos, weaponScale, rl.Gray)
}

// CalculateGunBarrelTip returns the world position of the gun barrel tip for tracers
func CalculateGunBarrelTip(p *player.Player, cameraPos rl.Vector3) rl.Vector3 {
	gunPos := calculateGunWorldPosition(p, cameraPos)

	// Forward vector for barrel extension
	forward := p.GetForwardVector()

	// Extend gun position by barrel length
	return rl.Vector3Add(gunPos, rl.Vector3Scale(forward, GUN_BARREL_LENGTH))
}

// calculateGunWorldPosition converts viewport position to world position for tracers
func calculateGunWorldPosition(p *player.Player, cameraPos rl.Vector3) rl.Vector3 {
	// Right vector (for X offset)
	right := p.GetRightVector()

	// Up vector (for Y offset)
	up := p.GetUpVector()

	// Forward vector (for Z offset)
	forward := p.GetForwardVector()

	// Start from camera position
	gunPos := cameraPos

	// Apply offsets in world space
	gunPos = rl.Vector3Add(gunPos, rl.Vector3Scale(right, GUN_OFFSET_X))
	gunPos = rl.Vector3Add(gunPos, rl.Vector3Scale(up, GUN_OFFSET_Y))
	gunPos = rl.Vector3Add(gunPos, rl.Vector3Scale(forward, GUN_OFFSET_Z))

	return gunPos
} 