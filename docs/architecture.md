# ShootPerfect Core Architecture

This document explains the architecture of **ShootPerfect Core** in a simple way.

The goal is that even someone new to the project should be able to read this document, open the code, and understand what each part is supposed to do.

---

## 1. What Is ShootPerfect Core?

ShootPerfect Core is the backend and video-analysis engine of the ShootPerfect system.

Its job is to take shooting-session videos, understand them, process them, and generate useful shooting technique metrics.

In simple words:

```text
Core takes videos in.
Core analyzes them.
Core gives structured results out.
```

ShootPerfect Studio will later use those results to show timelines, charts, overlays, and coaching feedback.

---

## 2. Big Picture

The first version of ShootPerfect Core is focused on a two-camera setup:

* one side-view camera
* one rear-view camera

The videos are recorded during a shooting session. Core will load these videos, synchronize them, identify shot windows, analyze movement, and produce results.

```text
+-------------------+
| Side View Video   |
+-------------------+
          |
          |
          v
+-------------------+
|                   |
| ShootPerfect Core |
|                   |
+-------------------+
          ^
          |
          |
+-------------------+
| Rear View Video   |
+-------------------+
          |
          v
+-------------------+
| Analysis Results  |
+-------------------+
```

---

## 3. Core and Studio Separation

ShootPerfect is split into two main parts:

```text
+------------------------+        +-------------------------+
| ShootPerfect Studio    |        | ShootPerfect Core       |
|------------------------|        |-------------------------|
| UI                     |        | Session management      |
| Video playback         |        | Video metadata          |
| Timeline viewer        | <----> | Synchronization         |
| Charts and reports     |        | Shot windows            |
| Coaching screen        |        | Video analysis          |
| Visual overlays        |        | Metrics generation      |
+------------------------+        +-------------------------+
```

### Core owns

* session data
* video metadata
* camera roles
* synchronization offsets
* shot timestamps
* frame extraction
* analysis logic
* metric generation
* result storage

### Studio owns

* user interface
* video player
* timeline display
* overlays
* charts
* reports
* coaching dashboard

The rule is simple:

```text
Core produces the data.
Studio displays the data.
```

---

## 4. First Workflow

The first useful version of Core should work without Studio.

It should run locally from a command or simple runner.

Example workflow:

```text
+----------------------+
| Load Session Config  |
+----------+-----------+
           |
           v
+----------------------+
| Register Videos      |
| - side view          |
| - rear view          |
+----------+-----------+
           |
           v
+----------------------+
| Apply Sync Offset    |
+----------+-----------+
           |
           v
+----------------------+
| Load Shot Markers    |
+----------+-----------+
           |
           v
+----------------------+
| Extract Shot Frames  |
+----------+-----------+
           |
           v
+----------------------+
| Run Basic Analysis   |
+----------+-----------+
           |
           v
+----------------------+
| Generate Metrics     |
+----------+-----------+
           |
           v
+----------------------+
| Write Result JSON    |
+----------------------+
```

This is the main Core-first approach.

The API can come later when Studio starts.

---

## 5. Initial Repository Structure

The planned structure is:

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

## 6. Why This Structure?

The project is divided into small modules so that each module has one clear responsibility.

This keeps the code easier to understand and easier to change later.

```text
session   -> What shooting session is this?
camera    -> Which camera view is this?
video     -> Which video file belongs to which camera?
sync      -> How are multiple videos aligned?
shots     -> Where is each shot in the video?
analysis  -> What processing should run?
metrics   -> What result did we calculate?
storage   -> Where do we save/load data?
api       -> How will Studio talk to Core later?
config    -> How do we load settings?
```

---

## 7. Module Responsibilities

### 7.1 Session Module

Location:

```text
internal/session/
```

The session module manages shooting sessions.

A session means one recorded training or match-like attempt.

Example:

```text
Session:
  ID: session-001
  Discipline: 10m Air Pistol
  Date: 2026-07-06
  Videos:
    - side-view video
    - rear-view video
  Shots:
    - shot 1
    - shot 2
    - shot 3
```

The session module should answer questions like:

```text
What session is this?
Which videos belong to it?
How many shots are marked?
Has analysis been run?
```

It should not do video analysis directly.

---

### 7.2 Camera Module

Location:

```text
internal/camera/
```

The camera module defines camera roles.

For the first version, we need:

```text
side_view
rear_view
```

Later we may add:

```text
front_view
target_view
trigger_hand_view
top_view
```

The camera module should help the system understand what each camera is supposed to capture.

Example:

```text
side_view -> used for arm, wrist, hold, follow-through
rear_view -> used for alignment, sway, pistol cant
```

---

### 7.3 Video Module

Location:

```text
internal/video/
```

The video module manages video files and video metadata.

It does not decide shooting technique. It only understands the video file.

Example video metadata:

```text
Video:
  ID: video-001
  SessionID: session-001
  CameraRole: side_view
  FilePath: ./testdata/session-001/side.mp4
  FPS: 60
  Duration: 420 seconds
  Width: 1920
  Height: 1080
```

The video module should answer:

```text
Where is the file?
Which camera role does it belong to?
What is the FPS?
What is the duration?
What is the resolution?
```

---

### 7.4 Sync Module

Location:

```text
internal/sync/
```

The sync module aligns videos from different cameras.

Since two phones may not start recording at the exact same time, we need a sync offset.

Example:

```text
side_view starts at 0 ms
rear_view starts 420 ms later
```

So the sync data may look like:

```text
Sync:
  side_view offset: 0 ms
  rear_view offset: 420 ms
```

Simple diagram:

```text
Side View:
0s -------- 10s -------- 20s -------- 30s

Rear View:
     0s -------- 10s -------- 20s -------- 30s
     ^
     starts 420 ms later
```

After sync:

```text
Side View:
0s -------- 10s -------- 20s -------- 30s

Rear View:
0s -------- 10s -------- 20s -------- 30s
```

In the first version, this offset can be entered manually.

Later, automatic sync can be added using audio spikes, clap markers, shot sound, or visible movement.

---

### 7.5 Shots Module

Location:

```text
internal/shots/
```

The shots module defines where each shot happens in the video.

A shot should not be treated as just one timestamp. For analysis, we need a time window around the shot.

Example:

```text
Shot 1:
  Trigger time: 72.500 seconds
  Pre-shot window: 6 seconds
  Follow-through window: 2 seconds
  Recovery window: 2 seconds
```

This means Core will analyze:

```text
from 66.500 seconds to 76.500 seconds
```

Simple diagram:

```text
66.5s              72.5s              74.5s       76.5s
 |------------------|------------------|-----------|
 pre-shot hold      trigger moment     follow      recovery
```

The shots module should answer:

```text
When did the shot happen?
What time window should be analyzed?
What part is pre-shot?
What part is follow-through?
```

---

### 7.6 Analysis Module

Location:

```text
internal/analysis/
```

The analysis module runs the actual analysis.

This is the heart of Core.

In the first version, analysis should be simple and non-AI.

Example early analysis:

```text
For each shot:
  1. Open side-view video
  2. Open rear-view video
  3. Extract frames from shot window
  4. Compare frame-to-frame movement
  5. Calculate basic movement scores
  6. Generate result
```

Simple analysis flow:

```text
+-------------------+
| Shot Window       |
+---------+---------+
          |
          v
+-------------------+
| Extract Frames    |
+---------+---------+
          |
          v
+-------------------+
| Compare Frames    |
+---------+---------+
          |
          v
+-------------------+
| Calculate Scores  |
+---------+---------+
          |
          v
+-------------------+
| Return Metrics    |
+-------------------+
```

The analysis module should not care about UI.

It should only produce structured results.

---

### 7.7 Metrics Module

Location:

```text
internal/metrics/
```

The metrics module defines the output of analysis.

Example metrics:

```text
Shot 1:
  hold_stability_score: 78
  follow_through_score: 64
  movement_before_shot: medium
  movement_after_shot: high
  body_sway_estimate: low
```

The exact scoring logic can change over time, but the structure should be clear.

Possible metric groups:

```text
Hold Metrics:
  - hold stability score
  - movement during final hold
  - stillness before trigger

Follow-through Metrics:
  - follow-through stability score
  - movement after shot
  - recovery movement

Body Metrics:
  - body sway
  - shoulder movement
  - head movement

Pistol Metrics:
  - pistol movement
  - pistol cant
  - muzzle movement estimate
```

Metrics should be versioned because our algorithm will improve over time.

Example:

```text
algorithm: basic-motion-v1
```

Later:

```text
algorithm: opencv-region-tracking-v2
```

---

### 7.8 Storage Module

Location:

```text
internal/storage/
```

The storage module saves and loads data.

In the first version, JSON files are enough.

Later, SQLite can be added.

The storage module may store:

```text
sessions.json
videos.json
shots.json
analysis-results.json
```

Or one file per session:

```text
testdata/
  session-001/
    session.yaml
    side.mp4
    rear.mp4
    result.json
```

The storage module should hide the storage details from other modules.

For example, the session module should not care whether data is saved in JSON or SQLite.

---

### 7.9 Config Module

Location:

```text
internal/config/
```

The config module loads application settings.

Example config:

```text
video_root: ./testdata
output_root: ./output
storage_type: json
default_pre_shot_window_seconds: 6
default_follow_through_window_seconds: 2
```

The config module should keep hardcoded values out of the main logic.

---

### 7.10 API Module

Location:

```text
internal/api/
```

The API module is a thin layer that will allow Studio to talk to Core later.

In the early versions, the API may only expose:

```text
GET /health
```

Later, when Studio starts, APIs can be added for:

```text
POST /sessions
POST /sessions/{id}/videos
POST /sessions/{id}/sync
POST /sessions/{id}/shots
POST /sessions/{id}/analysis
GET  /sessions/{id}/analysis
```

Important rule:

```text
API should not contain analysis logic.
API should call the correct Core services.
```

Example:

```text
Studio calls API
      |
      v
API handler receives request
      |
      v
Analysis service runs analysis
      |
      v
Storage saves result
      |
      v
API returns result status
```

---

## 8. Data Flow

This section explains how data moves through Core.

### 8.1 Session Creation Flow

```text
+---------------------+
| User creates config |
+----------+----------+
           |
           v
+---------------------+
| Session module      |
| validates session   |
+----------+----------+
           |
           v
+---------------------+
| Storage module      |
| saves session       |
+---------------------+
```

---

### 8.2 Video Registration Flow

```text
+---------------------+
| Video file provided |
+----------+----------+
           |
           v
+---------------------+
| Video module        |
| reads metadata      |
+----------+----------+
           |
           v
+---------------------+
| Camera module       |
| validates role      |
+----------+----------+
           |
           v
+---------------------+
| Storage module      |
| saves video data    |
+---------------------+
```

---

### 8.3 Analysis Flow

```text
+---------------------+
| Load session        |
+----------+----------+
           |
           v
+---------------------+
| Load videos         |
+----------+----------+
           |
           v
+---------------------+
| Load sync offsets   |
+----------+----------+
           |
           v
+---------------------+
| Load shot markers   |
+----------+----------+
           |
           v
+---------------------+
| Build shot windows  |
+----------+----------+
           |
           v
+---------------------+
| Extract frames      |
+----------+----------+
           |
           v
+---------------------+
| Run analysis        |
+----------+----------+
           |
           v
+---------------------+
| Generate metrics    |
+----------+----------+
           |
           v
+---------------------+
| Save result         |
+---------------------+
```

---

## 9. Example Session File

For early development, a simple YAML or JSON session file can be used.

Example:

```yaml
session_id: session-001
discipline: 10m_air_pistol
date: 2026-07-06

videos:
  - id: video-side-001
    camera_role: side_view
    path: ./testdata/session-001/side.mp4
    sync_offset_ms: 0

  - id: video-rear-001
    camera_role: rear_view
    path: ./testdata/session-001/rear.mp4
    sync_offset_ms: 420

shots:
  - id: shot-001
    trigger_time_ms: 72500
    pre_shot_window_ms: 6000
    follow_through_window_ms: 2000
    recovery_window_ms: 2000

  - id: shot-002
    trigger_time_ms: 93500
    pre_shot_window_ms: 6000
    follow_through_window_ms: 2000
    recovery_window_ms: 2000
```

This type of file is easy to understand and easy to test before building Studio.

---

## 10. Example Result File

Core should produce structured output.

Example:

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

This result can later be used by Studio to show charts, feedback, and overlays.

---

## 11. Initial Command-Line Flow

Before Studio exists, Core should be usable from the command line.

Example:

```bash
shootperfect-analyze --session ./testdata/session-001/session.yaml
```

Expected output:

```text
Loading session: session-001
Found side-view video
Found rear-view video
Loaded 2 shot markers
Running analysis: basic-motion-v1
Analysis complete
Result written to: ./output/session-001-result.json
```

This proves that Core works even before UI is built.

---

## 12. API Flow Later

When Studio development starts, Studio will call Core APIs instead of using local files directly.

```text
+----------------------+
| ShootPerfect Studio  |
+----------+-----------+
           |
           | HTTP REST API
           v
+----------------------+
| Core API Module      |
+----------+-----------+
           |
           v
+----------------------+
| Core Services        |
+----------+-----------+
           |
           v
+----------------------+
| Storage / Analysis   |
+----------------------+
```

The API should only translate external requests into internal service calls.

It should not become the main logic.

---

## 13. First Analysis Strategy

The first analysis should be simple.

Do not start with AI.

A basic first version can use frame-to-frame difference.

Simple idea:

```text
Frame 1 compared with Frame 2
Frame 2 compared with Frame 3
Frame 3 compared with Frame 4
...
```

If the difference is high, movement is high.

If the difference is low, movement is low.

Simple flow:

```text
+----------+       +----------+
| Frame 1  | ----> | Frame 2  |
+----------+       +----------+
      \             /
       \           /
        v         v
     Compare pixel/region difference
              |
              v
       Movement score
```

This can help estimate:

* whether the shooter was still before the shot
* whether there was visible movement during trigger execution
* whether follow-through was stable
* whether movement increased after the shot

This is not perfect, but it is enough to prove the pipeline.

---

## 14. How Side View and Rear View Are Used

### Side View

Side view is useful for:

```text
arm stability
wrist movement
trigger hand movement
head movement
follow-through
```

### Rear View

Rear view is useful for:

```text
body alignment
body sway
pistol cant
shoulder movement
weight shift
```

In the early version, both views can produce simple movement scores.

Later, the system can calculate more specific metrics for each view.

---

## 15. Basic Analysis Output Per Shot

For every shot, Core should eventually produce something like:

```text
Shot ID: shot-001

Side View:
  hold movement: low
  trigger-time movement: medium
  follow-through movement: high

Rear View:
  body sway: low
  alignment movement: medium
  pistol cant change: unknown

Overall:
  hold stability score: 78
  follow-through score: 64
  warning: movement increased after trigger
```

This is easier to understand than raw frame data.

---

## 16. Important Design Rules

### Rule 1: Keep Core Independent

Core should work without Studio.

Studio is a client. Core is the engine.

---

### Rule 2: Keep API Thin

API should not contain business logic.

Good:

```text
API -> Analysis Service -> Metrics -> Storage
```

Bad:

```text
API handler does everything
```

---

### Rule 3: Start With Manual Inputs

Manual sync and manual shot timestamps are acceptable in early versions.

Automatic detection can come later.

---

### Rule 4: Use Simple Analysis First

Do not start with AI.

Start with:

```text
frame extraction
frame comparison
basic movement score
JSON result
```

---

### Rule 5: Make Data Easy to Inspect

Early files should be human-readable.

Prefer simple YAML or JSON for test sessions and results.

This helps debugging.

---

### Rule 6: Do Not Hardcode Two Cameras Everywhere

The first version uses two cameras, but the model should allow more cameras later.

Good:

```text
videos: list of videos
camera_role: side_view / rear_view / etc.
```

Bad:

```text
sideVideoPath
rearVideoPath
thirdVideoPath
```

---

### Rule 7: Keep Analysis Versioned

Every result should mention which algorithm created it.

Example:

```text
basic-motion-v1
```

This is important because the algorithm will improve over time.

---

## 17. Development Order

Recommended order:

```text
1. Create Go module and folders
2. Add session, camera, and video models
3. Add sample session file
4. Add config loader
5. Add local runner
6. Load session and videos
7. Add shot model
8. Add sync offsets
9. Extract frames from shot windows
10. Add basic movement analysis
11. Write result JSON
12. Add minimal API later
```

This keeps the project practical and avoids building UI/backend layers before the actual engine works.

---

## 18. Final Simple Mental Model

The easiest way to understand ShootPerfect Core is this:

```text
Session tells us what happened.
Videos show us what happened.
Sync aligns the videos.
Shots tell us where to look.
Analysis studies those windows.
Metrics explain what was found.
Storage saves everything.
API exposes it later.
Studio displays it later.
```

That is the whole architecture.
