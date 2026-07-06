# ShootPerfect Core Roadmap

This document defines the development roadmap for ShootPerfect Core.

The goal is to build the system in small, practical stages. Each milestone should produce something useful and testable.

The project should stay Core-first. That means we should first prove that the analysis engine works locally before spending too much effort on Studio, UI, cloud, or large APIs.

---

## Roadmap Philosophy

ShootPerfect Core should not start as a big platform.

It should start as a simple local engine that can:

```text
load videos
understand a session
apply sync
read shot timestamps
extract frames
run basic analysis
write results
```

Once that works, we can add better storage, APIs, Studio integration, computer vision improvements, and eventually AI-based analysis.

---

## Milestone 1: Core Skeleton

### Goal

Create the basic Go project foundation.

This milestone is about setting up the project cleanly, not doing real analysis yet.

### Tasks

```text
- initialize Go module
- create folder structure
- add basic config loading
- add basic session model
- add basic camera model
- add basic video metadata model
- add minimal health endpoint
- add simple CLI/local runner entry point
```

### Expected Output

The project should build successfully.

Example:

```bash
go build ./...
```

A basic command should run without doing much.

Example:

```bash
shootperfect-analyze --help
```

The API may only expose:

```text
GET /health
```

### Done When

```text
- Go module exists
- project structure is created
- basic models compile
- health endpoint works
- local runner starts
```

---

## Milestone 2: Session + Video Registration

### Goal

Allow Core to understand a shooting session and the videos attached to it.

At this stage, Core should know:

```text
This is one shooting session.
This session has a side-view video.
This session has a rear-view video.
```

### Tasks

```text
- create/load session metadata
- define camera roles
- support side_view camera
- support rear_view camera
- register video files against a session
- store video path
- store camera role
- store FPS
- store duration
- store resolution
- validate required camera views
- validate supported video formats
```

### Expected Input

A simple session file.

Example:

```yaml
session_id: session-001
discipline: 10m_air_pistol

videos:
  - id: video-side-001
    camera_role: side_view
    path: ./testdata/session-001/side.mp4

  - id: video-rear-001
    camera_role: rear_view
    path: ./testdata/session-001/rear.mp4
```

### Expected Output

Core should be able to load the session and print something like:

```text
Loaded session: session-001
Found side-view video: ./testdata/session-001/side.mp4
Found rear-view video: ./testdata/session-001/rear.mp4
Session validation passed
```

### Done When

```text
- session file can be loaded
- videos can be registered
- camera roles are validated
- missing side/rear videos are detected
- invalid video paths are reported clearly
```

---

## Milestone 3: Local Persistence

### Goal

Make sessions and video metadata survive application restart.

Early persistence can be simple. JSON or YAML is fine in the beginning. SQLite can be added once the model stabilizes.

### Tasks

```text
- add JSON or SQLite storage
- persist sessions
- persist video metadata
- persist camera metadata
- add repository layer
- add basic tests
```

### Recommended First Approach

Start with file-based storage.

Example:

```text
testdata/
  session-001/
    session.yaml
    result.json
```

Later, move metadata into SQLite.

### Why Not Start With Full Database Design?

Because the data model will change while we learn what the analysis engine needs.

It is better to keep early data human-readable and easy to edit.

### Done When

```text
- session data can be saved
- session data can be loaded again
- video metadata is persisted
- tests cover basic save/load behavior
```

---

## Milestone 4: Manual Synchronization

### Goal

Align side-view and rear-view videos.

Different phones may start recording at slightly different times. Core needs to know the offset between the videos.

### Tasks

```text
- allow manual sync offset between videos
- store sync offset per video
- load sync metadata during analysis
- prepare synchronized shot windows
```

### Example

```yaml
videos:
  - id: video-side-001
    camera_role: side_view
    path: ./testdata/session-001/side.mp4
    sync_offset_ms: 0

  - id: video-rear-001
    camera_role: rear_view
    path: ./testdata/session-001/rear.mp4
    sync_offset_ms: 420
```

This means the rear-view video starts 420 milliseconds later than the side-view video.

### Done When

```text
- each video can have a sync offset
- sync offset is saved and loaded
- analysis can calculate the correct video timestamp per camera
```

---

## Milestone 5: Shot Marking

### Goal

Identify individual shot windows inside a session.

Core should not analyze the full video blindly. It should analyze useful windows around each shot.

### Tasks

```text
- add shot model
- allow manual shot timestamp entry
- define pre-shot window
- define trigger moment
- define follow-through window
- define recovery window
- store shot markers
```

### Example Shot

```yaml
shots:
  - id: shot-001
    trigger_time_ms: 72500
    pre_shot_window_ms: 6000
    follow_through_window_ms: 2000
    recovery_window_ms: 2000
```

This means:

```text
Trigger moment: 72.5 seconds
Analyze from:   66.5 seconds
Analyze until:  76.5 seconds
```

### Why Manual First?

Manual shot marking is good enough for early versions.

Automatic shot detection can come later after the basic pipeline is working.

### Done When

```text
- shots can be defined in session file
- shot windows can be calculated
- invalid shot windows are detected
- shot data can be saved and loaded
```

---

## Milestone 6: First Analysis Pipeline

### Goal

Run basic analysis on each shot window.

This is where ShootPerfect Core starts becoming an actual analysis engine.

### Tasks

```text
- define analysis job model
- define analysis result schema
- open video files
- extract frames from shot windows
- sample frames at fixed intervals
- produce basic placeholder movement result
- write result JSON
```

### First Pipeline Flow

```text
Load session
    |
    v
Load videos
    |
    v
Load sync offsets
    |
    v
Load shot markers
    |
    v
Build shot windows
    |
    v
Extract frames
    |
    v
Run basic analysis
    |
    v
Write result JSON
```

### Expected Output

Example:

```json
{
  "session_id": "session-001",
  "algorithm": "basic-motion-v1",
  "status": "completed",
  "shots": [
    {
      "shot_id": "shot-001",
      "status": "analyzed"
    }
  ]
}
```

### Done When

```text
- analysis command can run against a session file
- frames are extracted from shot windows
- each shot gets a result entry
- result JSON is written
- errors are reported clearly
```

---

## Milestone 7: First Non-AI Computer Vision Metrics

### Goal

Extract useful shooting metrics without AI.

This milestone should use simple computer-vision or frame-difference techniques before moving to AI/ML.

### Tasks

```text
- frame-to-frame motion detection
- region-based movement tracking
- hand/pistol movement estimate
- hold stability score
- follow-through stability score
- body sway approximation
- per-shot movement timeline
- session-level consistency summary
```

### Example Result

```json
{
  "session_id": "session-001",
  "algorithm": "basic-motion-v1",
  "shots": [
    {
      "shot_id": "shot-001",
      "metrics": {
        "hold_stability_score": 78,
        "follow_through_score": 64,
        "movement_before_shot": "medium",
        "movement_after_shot": "high",
        "body_sway_estimate": "low"
      },
      "warnings": [
        "Visible movement increased after trigger"
      ]
    }
  ]
}
```

### Done When

```text
- Core produces actual movement-based metrics
- metrics are stored in result JSON
- results are understandable without Studio
- algorithm name/version is included
```

---

## Later Milestones

The first 7 milestones are enough to build the first usable Core.

After that, future milestones may include:

```text
- better OpenCV-based tracking
- automatic shot detection
- automatic video synchronization
- pose estimation
- pistol cant estimation
- Studio API expansion
- Studio timeline integration
- report generation
- cloud storage
- user profiles
- AI-assisted coaching feedback
```

---

## Suggested Version Mapping

This is only a rough version plan.

```text
v0.1  Core skeleton
v0.2  Session and video registration
v0.3  Local persistence
v0.4  Manual synchronization
v0.5  Shot marking
v0.6  First analysis pipeline
v0.7  First non-AI CV metrics
v0.8  Improved metrics
v0.9  More test sessions and validation
v1.0  First usable local Core release
```

Studio work should ideally start after Core can produce meaningful analysis results.

A practical point to start Studio would be around:

```text
v1.0 or later
```

At that stage, Studio will have real data to display.

---

## What We Should Avoid Early

Avoid spending too much time early on:

```text
- full UI
- authentication
- cloud deployment
- complex REST APIs
- Kubernetes
- AI coaching
- mobile app
- multi-user dashboards
```

These are useful later, but they should not block the first working Core engine.

---

## Current Priority

The current priority is:

```text
1. README.md
2. docs/architecture.md
3. docs/roadmap.md
4. docs/api.md
5. Go project skeleton
6. basic models
7. sample session file
8. local analysis runner
```

The first real success point is simple:

```text
Given a session file and two videos, Core should produce a result JSON.
```
