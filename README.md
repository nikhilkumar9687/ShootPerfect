# ShootPerfect Core

ShootPerfect Core is the backend and video-analysis engine for ShootPerfect, a computer-vision-based training platform for **Olympic-style 10m Air Pistol** shooting.

The purpose of Core is to process shooting-session videos, manage session data, synchronize multiple camera views, define shot windows, and generate structured technique metrics.

Core should eventually support AI-assisted coaching, but the first versions will focus on simple, reliable, non-AI video analysis.

---

## Project Goal

ShootPerfect Core is being built to objectively analyze shooting technique using synchronized smartphone camera videos.

The first practical goal is to analyze a two-camera setup:

* Side-view camera
* Rear-view camera

Core should help detect and measure things like:

* body stability
* arm movement
* wrist movement
* pistol movement
* hold stability
* follow-through stability
* body sway
* pistol cant
* consistency across shots

The long-term goal is to evolve ShootPerfect into a digital shooting coach, but the initial focus is to build a clean and usable Core engine.

---

## Initial Camera Setup

The first version is designed around two smartphone cameras.

### Camera 1: Side View

The side-view camera is the primary analysis camera.

It is used for:

* posture
* arm stability
* wrist movement
* trigger hand movement
* hold duration
* follow-through
* head and upper-body movement

### Camera 2: Rear View

The rear-view camera is placed behind the shooter, slightly offset if required.

It is used for:

* Natural Point of Aim observation
* body alignment
* shoulder and hip alignment
* body sway
* pistol cant
* weight shift

The system should not be hardcoded only for two cameras. The architecture should allow additional camera views later, such as front view, target view, or close-up trigger-hand view.

---

## Core Responsibilities

ShootPerfect Core owns the backend and analysis engine.

Core is responsible for:

* shooting session management
* camera and video registration
* video metadata extraction
* multi-camera synchronization
* shot timestamp and shot-window management
* frame extraction from video
* non-AI video-analysis pipelines
* technique metric generation
* analysis result storage
* exposing APIs later for ShootPerfect Studio

The early focus is local analysis. APIs will remain minimal until Studio development begins.

---

## What Core Does Not Own

ShootPerfect Core does not own the main user interface.

The following belong to ShootPerfect Studio or another client application:

* video playback UI
* timeline viewer
* visual overlays
* charts and reports
* user-facing coaching screen
* manual annotation UI
* dashboard experience

Core produces structured data. Studio displays and interacts with that data.

---

## Core-First Development Approach

The project will be built Core-first.

The initial workflow should be:

```text
Load session metadata
        |
        v
Register side-view and rear-view videos
        |
        v
Apply manual synchronization
        |
        v
Define shot timestamps and shot windows
        |
        v
Extract frames for each shot
        |
        v
Run basic non-AI analysis
        |
        v
Generate structured metrics
        |
        v
Write analysis results
```

The API layer is not the center of the early project. It is only a boundary layer that will become important once ShootPerfect Studio starts consuming Core functionality.

---

## Suggested Repository Structure

```text
shootperfect-core/
  cmd/
    shootperfect-core/
      main.go

    shootperfect-analyze/
      main.go

  internal/
    session/
    camera/
    video/
    sync/
    shots/
    analysis/
    metrics/
    storage/
    config/
    api/

  docs/
    architecture.md
    api.md
    roadmap.md

  configs/
  testdata/

  go.mod
  go.sum
  README.md
```

---

## Main Modules

### Session Module

Manages shooting sessions.

A session represents one training session, match simulation, or recorded shooting attempt.

Example session information:

* session ID
* discipline
* date/time
* shooter/profile reference
* notes
* associated camera videos
* shot markers
* analysis status

---

### Camera Module

Defines camera roles and camera placement metadata.

Initial camera roles:

* side_view
* rear_view

Future camera roles may include:

* front_view
* target_view
* trigger_hand_view
* top_view

The camera model should remain flexible so more cameras can be added later.

---

### Video Module

Handles video registration and video metadata.

Example video metadata:

* video ID
* session ID
* camera role
* file path
* duration
* frame rate
* resolution
* codec
* sync offset
* import status

---

### Sync Module

Handles synchronization between multiple camera videos.

Initial synchronization will be manual.

Example:

```text
side_view video starts at 0 ms
rear_view video starts at +420 ms
```

Future synchronization may use:

* clap/marker detection
* audio spike detection
* shot sound detection
* visible movement detection
* scoring-system integration

---

### Shots Module

Manages shot timestamps and shot windows.

A shot window may include:

* pre-shot hold
* trigger moment
* follow-through
* recovery

Example:

```text
Shot 1:
  trigger_time: 72.500s
  pre_shot_window: 6s
  follow_through_window: 2s
  recovery_window: 2s
```

This allows analysis to run per shot instead of processing the full video blindly.

---

### Analysis Module

Runs the actual video-analysis pipeline.

Initial analysis will be non-AI and rule-based.

Early analysis may include:

* frame extraction
* frame-to-frame motion detection
* movement scoring
* region-based movement tracking
* stability calculation
* follow-through movement comparison

AI/ML-based analysis can be added later after the basic pipeline is stable.

---

### Metrics Module

Defines structured shooting metrics produced by analysis.

Example metrics:

* hold stability score
* follow-through stability score
* visible movement before shot
* visible movement after shot
* body sway estimate
* pistol movement estimate
* movement timeline per shot
* session consistency summary

Metrics should be versioned over time because analysis algorithms will evolve.

---

### Storage Module

Handles persistence.

Early options:

* JSON files for simple local development
* SQLite once the data model stabilizes

Core should store:

* sessions
* videos
* camera metadata
* synchronization offsets
* shot markers
* analysis jobs
* analysis results

---

### API Module

Provides an external interface to Core.

In early versions, this module should remain minimal.

Initial API may only include:

```text
GET /health
```

Later, when Studio development begins, APIs can expose:

```text
POST   /sessions
GET    /sessions
GET    /sessions/{id}
POST   /sessions/{id}/videos
GET    /sessions/{id}/videos
POST   /sessions/{id}/sync
POST   /sessions/{id}/shots
POST   /sessions/{id}/analysis
GET    /sessions/{id}/analysis
```

The API module should stay thin. It should not contain business logic or analysis logic.

---

## Initial Technology Choices

The initial implementation will use:

* Go
* local CLI runner
* local filesystem for videos
* JSON or SQLite for metadata
* REST API later for Studio integration

Possible future additions:

* OpenCV
* FFmpeg
* Python-based CV/AI worker
* PostgreSQL
* object storage
* desktop app integration
* cloud deployment
* Kubernetes deployment

---

## Development Milestones

### Milestone 1: Core Skeleton

Goal: create the basic Go project foundation.

Tasks:

* initialize Go module
* create project structure
* add config loading
* add basic session model
* add basic camera model
* add basic video metadata model
* add minimal health endpoint only
* add simple local runner/CLI entry point

---

### Milestone 2: Session + Video Registration

Goal: allow Core to understand a shooting session and its camera videos.

Tasks:

* create/load session metadata
* register side-view video
* register rear-view video
* define camera roles
* store video path, camera role, FPS, duration, and resolution
* validate supported video formats
* validate that required camera views exist

---

### Milestone 3: Local Persistence

Goal: make sessions and video metadata survive application restart.

Tasks:

* add JSON or SQLite storage
* persist sessions
* persist camera metadata
* persist video metadata
* add repository layer
* add basic tests

---

### Milestone 4: Manual Synchronization

Goal: align side-view and rear-view videos.

Tasks:

* allow manual sync offset between videos
* store per-video offset
* load sync metadata during analysis
* prepare synchronized timelines for analysis

---

### Milestone 5: Shot Marking

Goal: identify individual shot windows inside the session.

Tasks:

* add shot model
* allow manual shot timestamp entry
* define pre-shot window
* define trigger moment
* define follow-through window
* define recovery window
* store shot markers

---

### Milestone 6: First Analysis Pipeline

Goal: run basic analysis on each shot window.

Tasks:

* define analysis job model
* define analysis result schema
* open video files
* extract frames from shot windows
* sample frames at fixed intervals
* produce placeholder/basic movement result
* write result JSON

---

### Milestone 7: First Non-AI Computer Vision Metrics

Goal: start extracting useful shooting metrics without AI.

Tasks:

* frame-to-frame motion detection
* hand/pistol region movement tracking
* hold stability score
* follow-through stability score
* body sway approximation
* per-shot movement timeline
* session-level consistency summary

---

## Design Principles

### Core Is the Engine

Core should perform the actual analysis.

Studio should display results, but Core should produce the metrics.

---

### Keep the First Version Simple

The first version should avoid unnecessary complexity.

Do not start with:

* full UI
* cloud architecture
* AI coaching
* automatic pose estimation
* complex authentication
* large API surface

Start with:

* sessions
* videos
* synchronization metadata
* shot windows
* basic analysis
* result files

---

### API Is Secondary Initially

The API layer is useful, but it should not drive early development.

Core should first work through local runners and test data. Once useful results are generated, APIs can be added for Studio.

---

### Make Analysis Repeatable

The same input videos and shot timestamps should produce repeatable results.

Analysis results should include:

* input session ID
* video IDs
* shot IDs
* algorithm name/version
* generated metrics
* timestamps
* errors or warnings

---

### Keep Modules Independent

Each module should have a clear responsibility.

For example:

* session should not know OpenCV details
* analysis should not own video storage
* API should not contain business logic
* metrics should be structured and stable
* storage should be replaceable later

---

### Design for More Cameras Later

The initial system uses two cameras, but the data model should support more.

Avoid hardcoding logic that assumes only side and rear views forever.

---

## Current Status

Project starting point.

Immediate next steps:

1. Finalize README
2. Create architecture document
3. Create Go project skeleton
4. Add basic models
5. Add local session runner
6. Add sample test session data

---

## License

To be decided.