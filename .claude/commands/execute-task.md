# Executing a Task with a Subagent

## Steps
1. Open the task file from `tasks/active/TASK-NNN.md`
2. Update status to `in-progress` and set `Updated` date
3. Choose the appropriate agent type based on the task
4. Launch the agent with full context from the task file
5. Update the task file with results
6. If completed, update status to `done`

## Agent Selection Guide
- **explore**: For codebase exploration, answering questions about code structure
- **task**: For running commands (tests, builds, linting) where only pass/fail matters
- **general-purpose**: For complex multi-step implementation tasks

## Execution Template

When launching an agent, include:
1. Full task description and acceptance criteria
2. Relevant file paths and context
3. Specific commands to run
4. Expected output format

## Updating Task Status

After execution, update the task file:
```markdown
## Status Update
- **Status**: done (or blocked/in-progress)
- **Completed**: YYYY-MM-DD
- **Result**: <Brief description of what was accomplished>
```

## Error Handling
- If agent fails, add notes under `## Notes` section
- Change status to `blocked` if external dependency prevents completion
- Re-assign to different agent type if initial approach fails
