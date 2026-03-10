# Creating a New Task

## Steps
1. Determine the next task ID from the registry in `orchestrator.md`
2. Create a new file `tasks/active/TASK-NNN.md` using the template below
3. Add an entry to the task registry table in `orchestrator.md`

## Task File Template

```markdown
# TASK-NNN: <Title>

## Metadata
- **ID**: TASK-NNN
- **Status**: pending
- **Priority**: medium
- **Dependencies**: none
- **Assigned Agent**: unassigned
- **Created**: YYYY-MM-DD
- **Updated**: YYYY-MM-DD

## Description
<Detailed description of what needs to be done>

## Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2

## Notes
<Any additional context or notes>
```

## Field Descriptions
- **ID**: Unique task identifier (TASK-NNN format)
- **Status**: Current state: `pending`, `in-progress`, `done`, `blocked`
- **Priority**: `critical`, `high`, `medium`, `low`
- **Dependencies**: List of TASK-NNN IDs this task depends on, or `none`
- **Assigned Agent**: Agent type: `general-purpose`, `explore`, `task`, or specific agent name
- **Created**: ISO date when task was created
- **Updated**: ISO date when task was last modified
