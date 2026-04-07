# UX Validator CLI

Command-line tool for [UX Validator](https://deltix.ai) — automated UX quality validation for mobile apps.

## Install

```bash
brew install whitingdeltix/tap/uxvalidator
```

Or download the binary from [Releases](https://github.com/whitingdeltix/deltix-cli/releases).

## Quick Start

```bash
# 1. Login
uxvalidator login

# 2. List your apps
uxvalidator apps

# 3. Run full UX validation (uses LLM)
uxvalidator validate --app <app-id>

# 4. Run playbook regression (no LLM, fast, free)
uxvalidator replay --app <app-id>
```

## Commands

### Authentication

```bash
uxvalidator login
```

Prompts for email and password. Saves token to `~/.uxvalidator/config.json`.

### List Apps

```bash
uxvalidator apps
```

```
ID          NAME      BUNDLE                 TASKS   SCORE
c919f385    Fitness   com.apple.Fitness      8       75
d90e4feb    Settings  com.apple.Preferences  7       66
```

### List Tasks

```bash
uxvalidator tasks --app <app-id>
```

### List Playbooks

```bash
uxvalidator playbooks --app <app-id>
```

```
ID          NAME                        STEPS   DIFFICULTY   PLATFORM
70712d60    turn_off_all_notifs         6       medium       ios
0b77fb7a    turn_off_activity_sharing   3       easy         ios
```

### Run Validation

Full UX validation using LLM (Claude). Discovers new UX issues and scores quality.

```bash
uxvalidator validate --app <app-id>
uxvalidator validate --app <app-id> --threshold 70
uxvalidator validate --app <app-id> --branch main --commit abc123
uxvalidator validate --app <app-id> --tasks <task-id-1>,<task-id-2>
uxvalidator validate --app <app-id> --device real_device
```

| Flag | Default | Description |
|------|---------|-------------|
| `--app` | required | App ID |
| `--threshold` | 0 | Minimum score — exits with code 1 if below |
| `--branch` | — | Branch name (stored on run) |
| `--commit` | — | Commit hash (stored on run) |
| `--tasks` | all | Comma-separated task IDs |
| `--device` | simulator | `simulator` or `real_device` |
| `--wait` | true | Wait for completion and show results |

### Run Playbooks

Replay saved playbooks deterministically without LLM. Fast and free.

```bash
# Replay all playbooks for an app
uxvalidator replay --app <app-id>

# Replay a specific playbook
uxvalidator replay --playbook <playbook-id>
```

```
UX Validator — Playbook Regression (no LLM)
App: Fitness (com.apple.Fitness)
Playbooks: 2

  ✓ turn_off_all_notifs      6 steps  1.2s  PASS
  ✗ turn_off_activity_share  3 steps  0.8s  FAIL
    Step 3: element 'Activity Sharing' not found

──────────────────────────────────────────────────
Result: 1/2 passed
```

### Check Status

```bash
uxvalidator status --run <run-id>
```

### View Results

```bash
uxvalidator results --run <run-id>
uxvalidator results --run <run-id> --json
```

## CI Integration

### GitHub Actions

```yaml
name: UX Validation

on:
  pull_request:
    branches: [main]

jobs:
  ux-validate:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4

      - name: Build app
        run: |
          xcodebuild build-for-testing \
            -scheme MyApp \
            -sdk iphonesimulator \
            -derivedDataPath ./build

      - name: Install UX Validator
        run: brew install whitingdeltix/tap/uxvalidator

      - name: Run playbook regression
        env:
          UX_VALIDATOR_KEY: ${{ secrets.UX_VALIDATOR_KEY }}
        run: |
          uxvalidator login
          uxvalidator replay \
            --app ${{ vars.UX_APP_ID }}

      - name: Run full validation
        env:
          UX_VALIDATOR_KEY: ${{ secrets.UX_VALIDATOR_KEY }}
        run: |
          uxvalidator validate \
            --app ${{ vars.UX_APP_ID }} \
            --branch ${{ github.head_ref }} \
            --commit ${{ github.sha }} \
            --threshold 70
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All passed / score above threshold |
| 1 | Failed / score below threshold |

## Configuration

Config is stored at `~/.uxvalidator/config.json`:

```json
{
  "api_url": "https://api.deltix.ai",
  "token": "...",
  "username": "your_username"
}
```

Override the API URL:

```bash
uxvalidator apps --api-url http://localhost:8080
```

## Two-Tier Testing

| | Validation | Playbooks |
|---|---|---|
| **Command** | `uxvalidator validate` | `uxvalidator replay` |
| **Uses LLM** | Yes (Claude) | No |
| **Speed** | Minutes | Seconds |
| **Cost** | Paid (LLM tokens) | Free |
| **Purpose** | Discover new UX issues | Catch regressions |
| **When to run** | Weekly / on demand | Every push |
