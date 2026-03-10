# Creating a New Task

All task state lives in a single file: **`.claude/orchestrator.json`**.

## Steps

1. Open `.claude/orchestrator.json`.
2. Read `counters.nextTaskId` to get the next numeric ID (e.g. `3` → `TASK-003`).
3. Read `counters.nextPlanId` to get the next plan ID if you are also creating a new plan (e.g. `2` → `PLAN-002`).
4. Add the new task object to the `tasks` map using the template below.
5. If the task belongs to a new plan, add a plan object to the `plans` map and set `activePlanId`.
6. If the task belongs to an existing plan, append its ID to that plan's `taskIds` array.
7. Increment `counters.nextTaskId` (and `counters.nextPlanId` if a new plan was created).
8. Update `updatedAt` to the current ISO 8601 timestamp.

## Task Object Template

```json
"TASK-NNN": {
  "id": "TASK-NNN",
  "planId": "PLAN-NNN",
  "title": "<Short, imperative title>",
  "status": "todo",
  "priority": "medium",
  "dependencies": [],
  "estimatedEffort": "M",
  "createdAt": "YYYY-MM-DDTHH:MM:SS.000Z",
  "startedAt": null,
  "completedAt": null,
  "filesChanged": [],
  "buildStatus": null,
  "testStatus": null,
  "notes": ""
}
```

## Plan Object Template

```json
"PLAN-NNN": {
  "id": "PLAN-NNN",
  "goal": "<One-sentence description of what this plan achieves>",
  "status": "in-progress",
  "taskIds": ["TASK-NNN"],
  "criticalPath": ["TASK-NNN"],
  "createdAt": "YYYY-MM-DDTHH:MM:SS.000Z"
}
```

## Field Reference

### Task fields
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `id` | string | ✅ | `TASK-NNN` (zero-padded, 3+ digits) |
| `planId` | string | ✅ | Parent plan ID |
| `title` | string | ✅ | Imperative phrase, ≤ 80 chars |
| `status` | string | ✅ | `todo` · `in-progress` · `done` · `blocked` · `deferred` |
| `priority` | string | ✅ | `critical` · `high` · `medium` · `low` |
| `dependencies` | string[] | ✅ | `[]` if none |
| `estimatedEffort` | string | ✅ | `XS` · `S` · `M` · `L` · `XL` |
| `createdAt` | ISO string | ✅ | |
| `startedAt` | ISO string \| null | ✅ | Set when status → `in-progress` |
| `completedAt` | ISO string \| null | ✅ | Set when status → `done` |
| `filesChanged` | string[] | ✅ | Relative repo paths; fill after completion |
| `buildStatus` | string \| null | ✅ | `pass` · `fail` · `n/a` · `null` |
| `testStatus` | string \| null | ✅ | `pass` · `fail` · `n/a` · `null` |
| `notes` | string | ✅ | Any extra context; use `""` if empty |
| `role` | string | ❌ | Optional agent role hint (e.g. `frontend`, `qa`, `verify`) |
| `deferralReason` | string | ❌ | Required when `status` is `deferred` |
