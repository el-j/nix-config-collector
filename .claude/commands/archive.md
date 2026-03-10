# Archiving Completed Tasks

## When to Archive
- Task status is `done`
- All acceptance criteria are checked off
- Any follow-up tasks have been created

## Steps
1. Verify the task is fully complete (all criteria checked)
2. Move the task file: `mv tasks/active/TASK-NNN.md tasks/archived/TASK-NNN.md`
3. Update `orchestrator.md`:
   - Update the task registry table row status to `done`
   - Optionally move the row to an "Archived Tasks" section

## Archive Format
When archiving, add a completion summary to the task file:

```markdown
## Completion Summary
- **Archived**: YYYY-MM-DD
- **Duration**: X days
- **Outcome**: <What was delivered>
- **Artifacts**: <List of files created/modified>
```

## Bulk Archive
To archive all completed tasks at once:
1. Run: `ls tasks/active/` to see active tasks
2. For each task with status `done`, follow the steps above
3. Update the registry table in `orchestrator.md`
