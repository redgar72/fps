# FPS Camera with Perfect Mouse Control

A complete first-person shooter camera implementation in Go using Raylib, featuring proper cursor capture and infinite mouse rotation.

## Features

- ✅ **Perfect cursor capture**: Hidden cursor with infinite rotation
- ✅ **FPS camera tied to player**: Camera follows player position + eye height
- ✅ **Smooth mouse look**: Separate yaw/pitch with proper clamping
- ✅ **WASD movement**: Movement relative to camera direction
- ✅ **Real-time rendering**: 60 FPS with 3D environment

## Setup

### Prerequisites (macOS)
```bash
# Install Raylib
brew install raylib
```

### Build and Run
```bash
# Clone/download the project
cd fps-demo

# Install Go dependencies
go mod tidy

# Run the demo
go run .
```

## Controls

- **Mouse**: Look around infinitely (cursor captured and hidden by default)
- **W/A/S/D**: Move forward/left/backward/right
- **Tab**: Toggle cursor capture on/off
- **ESC**: Exit application

## Core Concepts

### Camera Binding to Player Object

The key to FPS camera implementation is tying the camera to a player object:

```go
// Update camera position to player position (+ eye height)
eyeHeight := float32(1.7)
camera.Position = rl.NewVector3(playerPos.X, playerPos.Y + eyeHeight, playerPos.Z)

// Update camera target based on yaw and pitch rotation
targetDistance := float32(1.0)
camera.Target = rl.NewVector3(
    camera.Position.X + float32(math.Sin(float64(yaw))*math.Cos(float64(pitch)))*targetDistance,
    camera.Position.Y + float32(math.Sin(float64(pitch)))*targetDistance,
    camera.Position.Z + float32(math.Cos(float64(yaw))*math.Cos(float64(pitch)))*targetDistance,
)
```

### Perfect Mouse Control

Raylib provides proper cursor capture that enables infinite rotation:

```go
// Enable cursor capture for FPS controls
rl.DisableCursor()

// Get raw mouse movement (infinite because cursor is captured!)
mouseDelta := rl.GetMouseDelta()
sensitivity := float32(0.003)

// Update rotation based on mouse movement
yaw -= mouseDelta.X * sensitivity    // Horizontal rotation (unlimited)
pitch -= mouseDelta.Y * sensitivity  // Vertical rotation (clamped)
```

### Movement System

Movement is calculated relative to the camera's current orientation:

```go
// Calculate forward and right vectors based on yaw
forward := rl.NewVector3(
    float32(math.Sin(float64(yaw))),
    0,
    float32(math.Cos(float64(yaw))),
)
right := rl.NewVector3(
    float32(math.Cos(float64(yaw))),
    0,
    -float32(math.Sin(float64(yaw))),
)

// Apply movement based on input
if rl.IsKeyDown(rl.KeyW) {
    playerPos = rl.Vector3Add(playerPos, rl.Vector3Scale(forward, moveSpeed))
}
```

## Why Raylib?

This implementation uses Raylib instead of other Go 3D engines because:

1. **Proper cursor capture**: `rl.DisableCursor()` provides true FPS-style mouse control
2. **Cross-platform**: Works consistently across Windows, macOS, and Linux
3. **Simple API**: Easy to understand and extend
4. **Active development**: Well-maintained with regular updates
5. **Performance**: Efficient C library with Go bindings

## Architecture

```
Player Object (position + rotation)
        ↓
Camera (follows player position + eye offset)
        ↓  
Rendering (3D scene from camera perspective)
```

The camera is essentially "attached" to an invisible player object, creating the FPS experience where the camera represents the player's eyes.

## Extending the Demo

- Add physics and collision detection
- Load 3D models for more complex environments  
- Implement weapon systems and UI
- Add multiplayer networking
- Create level loading from files

This foundation provides everything needed for a complete FPS game engine!