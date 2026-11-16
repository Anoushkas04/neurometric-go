# NeuroMetric

A web platform for repeated cognitive screening. Users register, complete a battery of game-style cognitive tasks, and track their results over multiple sessions with charts and written interpretations.

This was a small early prototype that fed into [Elaris Island](https://github.com/Anoushkas04/elaris-v3), which takes the same idea of measuring cognition through games and builds it into a narrative game.

**Stack:** Go · Gin · GORM · SQLite · JWT + bcrypt · vanilla JavaScript · Tailwind CSS · Chart.js · Docker

## Features

- **Assessment battery (`merged_games.html`):** ten timed tasks run in sequence:
  1. Tile Memory: remember and reproduce a visual sequence
  2. Delayed Recall (encode): items to be recalled later
  3. Attention Vigilance: respond to targets in a stream
  4. Pathfinder: connect nodes in order while switching rules
  5. Recognition Memory: recall the items from task 2
  6. Spatial Reasoning
  7. Clock Drawing
  8. Kitchen Rush, in three stages of increasing load: conveyor focus, ticket (word) counting, and a combined multitasking round
- **Metrics:** each session records an overall score, accuracy, total errors, completion time and average reaction time, plus a per-task breakdown stored as JSON.
- **Reports:** a trend line across sessions, a radar chart of the cognitive profile, per-task bar charts, a plain-language interpretation page, and recommendations based on weaker areas.
- **Auth:** registration and login with bcrypt-hashed passwords and JWTs. Assessment endpoints are scoped to the signed-in user.

## API

| Method | Endpoint | Auth |
|---|---|---|
| POST | `/api/auth/register` | — |
| POST | `/api/auth/login` | — |
| GET | `/api/user/profile` | JWT |
| POST | `/api/assessments` | JWT |
| GET | `/api/assessments` | JWT |

## Voice task prototype

`public/ai_processor.py` is a standalone script for a sentence-repetition task. It records from the microphone, transcribes with OpenAI Whisper (`base`), and checks that the key words of the target sentence are present. It isn't connected to the Go server yet.

## Running locally

```bash
go run main.go           # http://localhost:8080, database at data/neurometric.db
```

Or with Docker:

```bash
docker build -t neurometric .
docker run -p 8080:8080 -v $(pwd)/data:/app/data neurometric
```
