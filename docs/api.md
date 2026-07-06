# ShootPerfect Core API

This document describes the planned API for ShootPerfect Core.

The API is not the main focus in the early versions of the project. Core should first work locally through a command-line runner and session files.

The API will become important later when ShootPerfect Studio starts using Core.

---

## 1. API Philosophy

ShootPerfect Core is the analysis engine.

The API is only a thin access layer.

In simple terms:

```text
Studio calls API.
API calls Core services.
Core services do the actual work.
API returns the result.
```

The API should not contain business logic or analysis logic.

Good:

```text
API -> Session Service -> Storage
API -> Analysis Service -> Metrics -> Storage
```

Bad:

```text
API handler directly does video analysis
API handler directly calculates metrics
API handler directly manages too many details
```

---

## 2. API Status in Early Versions

In early versions, the API can remain very small.

Initial endpoint:

```text
GET /health
```

That is enough to confirm that the Core service starts.

Most early work should happen through:

```bash
shootperfect-analyze --session ./testdata/session-001/session.yaml
```

The larger API surface should be added when Studio development begins.

---

## 3. Planned API Versions

The API should be versioned.

Initial version:

```text
/api/v1
```

Example:

```text
GET /api/v1/health
```

This gives us room to change the API later without breaking Studio.

---

## 4. Common Response Format

For simple responses, use JSON.

Example success response:

```json
{
  "status": "ok",
  "message": "request completed"
}
```

Example error response:

```json
{
  "status": "error",
  "error": {
    "code": "SESSION_NOT_FOUND",
    "message": "session not found"
  }
}
```

Keep error messages clear and human-readable.

---

## 5. Health API

### GET `/api/v1/health`

Checks whether ShootPerfect Core is running.

### Response

```json
{
  "status": "ok",
  "service": "shootperfect-core"
}
```

### Purpose

Used by developers, Studio, or deployment tools to confirm that Core is alive.

---

## 6. Session APIs

These APIs will be useful once Studio starts.

### POST `/api/v1/sessions`

Creates a new shooting session.

### Request

```json
{
  "discipline": "10m_air_pistol",
  "shooter_id": "default",
  "notes": "Evening training session"
}
```

### Response

```json
{
  "session_id": "session-001",
  "discipline": "10m_air_pistol",
  "status": "created"
}
```

---

### GET `/api/v1/sessions`

Lists available sessions.

### Response

```json
{
  "sessions": [
    {
      "session_id": "session-001",
      "discipline": "10m_air_pistol",
      "created_at": "2026-07-06T18:30:00+05:30",
      "status": "created"
    }
  ]
}
```

---

### GET `/api/v1/sessions/{session_id}`

Gets details of one session.

### Response

```json
{
  "session_id": "session-001",
  "discipline": "10m_air_pistol",
  "videos": [
    {
      "video_id": "video-side-001",
      "camera_role": "side_view"
    },
    {
      "video_id": "video-rear-001",
      "camera_role": "rear_view"
    }
  ],
  "shots": [
    {
      "shot_id": "shot-001",
      "trigger_time_ms": 72500
    }
  ],
  "analysis_status": "not_started"
}
```

---

## 7. Video APIs

### POST `/api/v1/sessions/{session_id}/videos`

Registers a video against a session.

In early versions, this may register a local file path. Later, this can support file upload.

### Request

```json
{
  "camera_role": "side_view",
  "file_path": "./testdata/session-001/side.mp4"
}
```

### Response

```json
{
  "video_id": "video-side-001",
  "session_id": "session-001",
  "camera_role": "side_view",
  "status": "registered"
}
```

---

### GET `/api/v1/sessions/{session_id}/videos`

Lists videos registered for a session.

### Response

```json
{
  "session_id": "session-001",
  "videos": [
    {
      "video_id": "video-side-001",
      "camera_role": "side_view",
      "file_path": "./testdata/session-001/side.mp4",
      "fps": 60,
      "duration_ms": 420000,
      "width": 1920,
      "height": 1080
    },
    {
      "video_id": "video-rear-001",
      "camera_role": "rear_view",
      "file_path": "./testdata/session-001/rear.mp4",
      "fps": 60,
      "duration_ms": 420000,
      "width": 1920,
      "height": 1080
    }
  ]
}
```

---

## 8. Sync APIs

### POST `/api/v1/sessions/{session_id}/sync`

Stores sync offset information for videos in a session.

### Request

```json
{
  "offsets": [
    {
      "video_id": "video-side-001",
      "sync_offset_ms": 0
    },
    {
      "video_id": "video-rear-001",
      "sync_offset_ms": 420
    }
  ]
}
```

### Response

```json
{
  "session_id": "session-001",
  "status": "sync_updated"
}
```

---

### GET `/api/v1/sessions/{session_id}/sync`

Gets sync offset information for a session.

### Response

```json
{
  "session_id": "session-001",
  "offsets": [
    {
      "video_id": "video-side-001",
      "camera_role": "side_view",
      "sync_offset_ms": 0
    },
    {
      "video_id": "video-rear-001",
      "camera_role": "rear_view",
      "sync_offset_ms": 420
    }
  ]
}
```

---

## 9. Shot APIs

### POST `/api/v1/sessions/{session_id}/shots`

Adds a shot marker.

### Request

```json
{
  "trigger_time_ms": 72500,
  "pre_shot_window_ms": 6000,
  "follow_through_window_ms": 2000,
  "recovery_window_ms": 2000
}
```

### Response

```json
{
  "shot_id": "shot-001",
  "session_id": "session-001",
  "status": "created"
}
```

---

### GET `/api/v1/sessions/{session_id}/shots`

Lists shot markers for a session.

### Response

```json
{
  "session_id": "session-001",
  "shots": [
    {
      "shot_id": "shot-001",
      "trigger_time_ms": 72500,
      "pre_shot_window_ms": 6000,
      "follow_through_window_ms": 2000,
      "recovery_window_ms": 2000
    }
  ]
}
```

---

## 10. Analysis APIs

These APIs should be added only after the local analysis runner works.

### POST `/api/v1/sessions/{session_id}/analysis`

Starts analysis for a session.

### Request

```json
{
  "algorithm": "basic-motion-v1"
}
```

### Response

```json
{
  "analysis_id": "analysis-001",
  "session_id": "session-001",
  "algorithm": "basic-motion-v1",
  "status": "queued"
}
```

In early versions, analysis may run synchronously.

Later, it can become an async job.

---

### GET `/api/v1/sessions/{session_id}/analysis`

Gets the latest analysis result for a session.

### Response

```json
{
  "session_id": "session-001",
  "analysis_id": "analysis-001",
  "algorithm": "basic-motion-v1",
  "status": "completed",
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

---

### GET `/api/v1/sessions/{session_id}/analysis/{analysis_id}`

Gets a specific analysis result.

### Response

```json
{
  "session_id": "session-001",
  "analysis_id": "analysis-001",
  "algorithm": "basic-motion-v1",
  "status": "completed",
  "created_at": "2026-07-06T18:30:00+05:30",
  "completed_at": "2026-07-06T18:31:20+05:30",
  "shots": [
    {
      "shot_id": "shot-001",
      "metrics": {
        "hold_stability_score": 78,
        "follow_through_score": 64
      }
    }
  ]
}
```

---

## 11. Suggested Endpoint Summary

Early API:

```text
GET    /api/v1/health
```

Later Studio API:

```text
POST   /api/v1/sessions
GET    /api/v1/sessions
GET    /api/v1/sessions/{session_id}

POST   /api/v1/sessions/{session_id}/videos
GET    /api/v1/sessions/{session_id}/videos

POST   /api/v1/sessions/{session_id}/sync
GET    /api/v1/sessions/{session_id}/sync

POST   /api/v1/sessions/{session_id}/shots
GET    /api/v1/sessions/{session_id}/shots

POST   /api/v1/sessions/{session_id}/analysis
GET    /api/v1/sessions/{session_id}/analysis
GET    /api/v1/sessions/{session_id}/analysis/{analysis_id}
```

---

## 12. API and Core Service Mapping

The API should map to internal modules clearly.

```text
/api/v1/sessions
    -> internal/session

/api/v1/sessions/{id}/videos
    -> internal/video
    -> internal/camera

/api/v1/sessions/{id}/sync
    -> internal/sync

/api/v1/sessions/{id}/shots
    -> internal/shots

/api/v1/sessions/{id}/analysis
    -> internal/analysis
    -> internal/metrics
    -> internal/storage
```

---

## 13. Error Codes

Suggested error codes:

```text
SESSION_NOT_FOUND
VIDEO_NOT_FOUND
INVALID_CAMERA_ROLE
INVALID_VIDEO_PATH
VIDEO_METADATA_FAILED
SYNC_OFFSET_INVALID
SHOT_NOT_FOUND
SHOT_WINDOW_INVALID
ANALYSIS_NOT_FOUND
ANALYSIS_FAILED
INTERNAL_ERROR
```

Example:

```json
{
  "status": "error",
  "error": {
    "code": "INVALID_CAMERA_ROLE",
    "message": "camera role must be side_view or rear_view"
  }
}
```

---

## 14. API Design Rules

### Rule 1: Keep API Thin

The API should only receive requests, validate input, call services, and return responses.

---

### Rule 2: Do Not Put Analysis Logic in Handlers

Wrong:

```text
handler extracts video frames
handler compares frames
handler calculates scores
```

Correct:

```text
handler calls analysis service
analysis service runs pipeline
metrics module defines result
storage saves result
```

---

### Rule 3: Keep IDs Clear

Use clear IDs such as:

```text
session-001
video-side-001
video-rear-001
shot-001
analysis-001
```

Later these can become UUIDs.

---

### Rule 4: Keep Studio Needs in Mind, But Do Not Build Everything Early

The API should eventually support Studio, but early Core work should not be blocked by API design.

---

### Rule 5: Version Analysis Results

Every analysis result should include an algorithm name or version.

Example:

```json
{
  "algorithm": "basic-motion-v1"
}
```

This is important because the analysis logic will improve over time.

---

## 15. First API Implementation Plan

When we start coding, the first API implementation should be very small.

Step 1:

```text
GET /api/v1/health
```

Step 2:

```text
POST /api/v1/sessions
GET  /api/v1/sessions/{session_id}
```

Step 3:

```text
POST /api/v1/sessions/{session_id}/videos
GET  /api/v1/sessions/{session_id}/videos
```

The remaining APIs can come later.

---

## 16. Final Note

The API is important, but it is not the heart of ShootPerfect Core.

The heart of Core is:

```text
session
videos
sync
shots
analysis
metrics
results
```

The API is just one way for Studio or another client to access that engine.
