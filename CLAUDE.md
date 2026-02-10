# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
go build -o unicorn-toots .
./unicorn-toots
```

## Architecture

Single-file Go application (`main.go`) using the [gopxl/pixel](https://github.com/gopxl/pixel) v2 OpenGL game library. Renders an 800x600 window with a movable square controlled by arrow keys.

Key dependencies:
- `gopxl/pixel/v2` — 2D game library (window, input, drawing)
- `gopxl/pixel/v2/backends/opengl` — OpenGL backend (requires `opengl.Run()` on main thread)
- `golang.org/x/image/colornames` — named colors

The entry point uses `opengl.Run(run)` which is required by pixel to execute the game loop on the main goroutine.
