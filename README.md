# FPS Game Engine

A first-person shooter game engine built in Go using Raylib, featuring modular architecture, enemy AI, and smooth anti-aliased graphics.


## Quick Start

```bash
# Install dependencies
go mod tidy

# Run the game
go run ./cmd

# Or build and run
go build ./cmd
./cmd/fps
```

## Controls

- **Mouse**: Look around
- **WASD**: Move
- **Left Click**: Shoot
- **Tab**: Toggle cursor capture
- **ESC**: Exit

## Architecture

```
fps/
├── cmd/           # Main entry point
├── internal/      # Game packages
│   ├── game/     # Game state
│   ├── player/   # Player system
│   ├── enemy/    # Enemy AI
│   ├── input/    # Input handling
│   ├── physics/  # Collision
│   └── rendering/ # Visual systems
└── assets/       # Game assets
```

## Development

The modular architecture makes it easy to extend:

```go
// Add new enemy types
func (e *CustomEnemy) Update(deltaTime float32) {
    // Custom AI logic
}

// Add new weapons
func (w *Weapon) Fire(player *player.Player) {
    // Custom firing logic
}
```

## Customization

Edit constants in respective packages:
```go
// internal/input/input.go
const MOVE_SPEED = 5.0        // Adjust movement speed
const MOUSE_SENSITIVITY = 0.003 // Adjust mouse sensitivity
```

## Troubleshooting

- **"No Go files in directory"**: Use `go run ./cmd` instead of `go run .`
- **Import errors**: Run `go mod tidy`
- **Model loading fails**: Check assets are in `assets/` directory

---

**Built with Go + Raylib** 🎮