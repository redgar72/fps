package rendering

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"fps/internal/player"
	"fps/internal/enemy"
	"fps/internal/physics"
)

// Constants for rendering
const (
	CROSSHAIR_SIZE      = 10
	CROSSHAIR_THICKNESS = 2
	TRACER_DURATION     = 0.5
)

// RenderWorld draws the 3D world elements
func RenderWorld(p *player.Player, e *enemy.Enemy, cubes []rl.Vector3, colors []rl.Color, hitTimers []float32, tm *physics.TracerManager, camera rl.Camera3D) {
	rl.BeginMode3D(camera)

	// Draw ground plane with improved visual quality
	groundColor := rl.NewColor(34, 139, 34, 255) // Forest green
	rl.DrawPlane(rl.Vector3{X: 0, Y: 0, Z: 0}, rl.Vector2{X: 20, Y: 20}, groundColor)
	
	// Draw subtle grid lines on the ground for better depth perception
	gridColor := rl.NewColor(0, 100, 0, 100) // Darker green, semi-transparent
	for i := -10; i <= 10; i++ {
		// Vertical lines
		start := rl.Vector3{X: float32(i), Y: 0.01, Z: -10}
		end := rl.Vector3{X: float32(i), Y: 0.01, Z: 10}
		rl.DrawLine3D(start, end, gridColor)
		
		// Horizontal lines
		start = rl.Vector3{X: -10, Y: 0.01, Z: float32(i)}
		end = rl.Vector3{X: 10, Y: 0.01, Z: float32(i)}
		rl.DrawLine3D(start, end, gridColor)
	}

	// Draw cubes with improved visual quality
	for i, pos := range cubes {
		// Draw main cube with current color
		rl.DrawCube(pos, 1, 1, 1, colors[i%len(colors)])
		
		// Draw wireframe with softer color for less harsh edges
		wireframeColor := rl.NewColor(50, 50, 50, 255) // Dark gray instead of pure black
		rl.DrawCubeWires(pos, 1, 1, 1, wireframeColor)
	}

	// Draw enemy capsule
	enemyColor := rl.NewColor(255, 0, 100, 255) // Bright pink/magenta for visibility
	
	// Change color based on hit status
	if e.HitTimer > 0 {
		enemyColor = rl.White // Flash white when hit
	} else if e.Health < 50 {
		enemyColor = rl.Red // Red when low health
	}
	
	enemyHeight := float32(2.0) // Height of the capsule
	enemyRadius := float32(0.5) // Radius of the capsule
	
	// Draw the enemy capsule (cylinder with rounded ends)
	rl.DrawCylinder(e.Position, enemyRadius, enemyRadius, enemyHeight, 8, enemyColor)
	
	// Draw wireframe for the enemy
	enemyWireframeColor := rl.NewColor(100, 0, 50, 255) // Darker pink for wireframe
	rl.DrawCylinderWires(e.Position, enemyRadius, enemyRadius, enemyHeight, 8, enemyWireframeColor)

	// Draw player representation (semi-transparent box)
	rl.DrawCube(p.Position, 0.5, player.EYE_HEIGHT, 0.5, rl.NewColor(255, 0, 0, 100))

	// Draw tracers
	renderTracers(tm)

	rl.EndMode3D()
}

// renderTracers draws all active tracer lines with fade effect
func renderTracers(tm *physics.TracerManager) {
	for _, tracer := range tm.GetActiveTracers() {
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

// RenderUI draws all UI elements
func RenderUI(p *player.Player, e *enemy.Enemy) {
	// Title and controls
	rl.DrawText("FPS Camera with Perfect Mouse Control!", 10, 10, 20, rl.DarkGray)
	rl.DrawText("WASD: Move | Mouse: Look | Left Click: Shoot | Tab: Toggle cursor | ESC: Exit", 10, 35, 16, rl.DarkGray)

	// Cursor status
	if rl.IsCursorHidden() {
		rl.DrawText("Mouse captured - infinite rotation!", 10, 55, 16, rl.Green)
	} else {
		rl.DrawText("Mouse free - press Tab to capture", 10, 55, 16, rl.Red)
	}

	// Enemy health display
	enemyHealthText := fmt.Sprintf("Enemy Health: %.0f", e.Health)
	rl.DrawText(enemyHealthText, 10, 75, 16, rl.White)

	// FPS counter in top right
	fps := rl.GetFPS()
	fpsText := fmt.Sprintf("FPS: %d", fps)
	fpsTextWidth := rl.MeasureText(fpsText, 20)
	screenWidth := int32(rl.GetScreenWidth())
	rl.DrawText(fpsText, screenWidth-fpsTextWidth-10, 10, 20, rl.White)
}

// RenderCrosshair draws the crosshair when in FPS mode
func RenderCrosshair() {
	if rl.IsCursorHidden() {
		centerX := int32(rl.GetScreenWidth() / 2)
		centerY := int32(rl.GetScreenHeight() / 2)

		// Draw horizontal line
	rl.DrawRectangle(centerX-CROSSHAIR_SIZE, centerY-CROSSHAIR_THICKNESS/2, CROSSHAIR_SIZE*2, CROSSHAIR_THICKNESS, rl.White)
		// Draw vertical line
	rl.DrawRectangle(centerX-CROSSHAIR_THICKNESS/2, centerY-CROSSHAIR_SIZE, CROSSHAIR_THICKNESS, CROSSHAIR_SIZE*2, rl.White)
	}
} 