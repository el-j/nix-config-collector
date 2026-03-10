# Archiving Completed Plans

All state lives in **`.claude/orchestrator.json`**. "Archiving" means marking a plan `done` and optionally trimming its completed tasks from the hot section — the JSON file is the single source of truth; no files are moved.

## When to Archive a Plan

Archive a plan when **all** of the following are true:
- Every task in `plan.taskIds` has `status: "done"` or `status: "deferred"`.
- No open follow-up tasks reference this plan.

## Steps

### 1. Verify completeness

Check every `taskId` in the plan:
```jsonc
// All of these must be "done" or "deferred":
"tasks": {
  "TASK-NNN": { "status": "done", ... },
  ...
}
```

### 2. Mark the plan done

```json
"PLAN-NNN": {
  ...
  "status": "done"
}
```

### 3. Update `activePlanId`

Set `activePlanId` to the next in-progress plan's ID, or `null` if there are none.

### 4. Update root `updatedAt`

```json
{
  "updatedAt": "YYYY-MM-DDTHH:MM:SS.000Z",
  ...
}
```

### 5. Optionally prune old task data

Once a plan is `done`, its tasks may be removed from the `tasks` map to keep the file small. Before doing so, update the root `notes` field to record what was removed:

```json
"notes": "Tasks TASK-001..TASK-NNN from PLAN-001 removed after archival. Full history in git log."
```

## Example — before archiving

```json
{
  "activePlanId": "PLAN-001",
  "plans": {
    "PLAN-001": { "status": "in-progress", "taskIds": ["TASK-001"] }
  },
  "tasks": {
    "TASK-001": { "status": "done", ... }
  }
}
```

## Example — after archiving

```json
{
  "updatedAt": "YYYY-MM-DDTHH:MM:SS.000Z",
  "activePlanId": null,
  "notes": "PLAN-001 archived. Tasks TASK-001 removed; see git history.",
  "plans": {
    "PLAN-001": { "status": "done", "taskIds": ["TASK-001"] }
  },
  "tasks": {}
}
```

## Bulk Archive

To close out multiple plans at once, repeat steps 1–5 for each plan in sequence, then do a single `updatedAt` update at the end.

## Important

- **Never delete plan objects** from the `plans` map — they serve as a permanent index.
- Removing tasks from the `tasks` map is optional; always document removed ranges in `notes`.
- All history is recoverable via `git log -- .claude/orchestrator.json`.
