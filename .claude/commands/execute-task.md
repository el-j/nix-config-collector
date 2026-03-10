# Executing a Task with a Subagent

All task state lives in **`.claude/orchestrator.json`**. Update it in place — never create separate task files.

## Steps

1. **Claim the task**: In `orchestrator.json`, set the task's `status` to `"in-progress"`, set `startedAt` to the current ISO timestamp, and update `updatedAt`.
2. **Gather context**: Read the task's `title`, `notes`, `dependencies`, and look up any dependency tasks in the same `tasks` map.
3. **Select an agent** using the guide below.
4. **Launch the agent** with full context (task description, acceptance criteria, relevant file paths, commands to run).
5. **Record results**: After the agent completes, update the task object:
   - `status` → `"done"` (or `"blocked"` if blocked)
   - `completedAt` → current ISO timestamp
   - `filesChanged` → list of files created or modified
   - `buildStatus` → `"pass"` / `"fail"` / `"n/a"`
   - `testStatus` → `"pass"` / `"fail"` / `"n/a"`
   - `notes` → brief summary of what was accomplished
6. Update `updatedAt` on the root object.

## Agent Selection Guide

| Agent type | When to use |
|------------|-------------|
| `explore` | Codebase exploration, answering questions about structure or logic |
| `task` | Running commands (tests, builds, lint) where only pass/fail matters |
| `general-purpose` | Complex multi-step implementation requiring reasoning and file edits |

## Minimal JSON Diff — Marking a Task In-Progress

```json
"TASK-NNN": {
  ...
  "status": "in-progress",
  "startedAt": "YYYY-MM-DDTHH:MM:SS.000Z",
  ...
}
```

## Minimal JSON Diff — Marking a Task Done

```json
"TASK-NNN": {
  ...
  "status": "done",
  "completedAt": "YYYY-MM-DDTHH:MM:SS.000Z",
  "filesChanged": ["path/to/file.go", "path/to/other.go"],
  "buildStatus": "pass",
  "testStatus": "pass",
  "notes": "Brief description of outcome."
}
```

## Blocking a Task

If a task cannot proceed due to an external dependency or upstream blocker:

```json
"TASK-NNN": {
  ...
  "status": "blocked",
  "notes": "Blocked by: <reason>. Waiting on: <TASK-MMM or external>."
}
```

## Deferring a Task

If a task is intentionally postponed:

```json
"TASK-NNN": {
  ...
  "status": "deferred",
  "deferralReason": "<Explanation of why deferred and any conditions for re-activation>"
}
```
