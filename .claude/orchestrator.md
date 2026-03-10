# Claude Orchestration System

## Overview
This orchestration system tracks tasks, coordinates subagent execution, and maintains workflow state for the nix-config-collector project.

## Task Registry

| ID | Title | Status | Priority | Agent | Created |
|----|-------|--------|----------|-------|---------|
| TASK-001 | Initial project scaffolding | done | high | general-purpose | 2024-01-01 |

## Task States
- `pending` - Task defined but not started
- `in-progress` - Task currently being executed
- `done` - Task completed successfully
- `blocked` - Task blocked by a dependency

## Workflow

1. **Create Task**: Use `.claude/commands/new-task.md` template
2. **Execute Task**: Assign to agent via `.claude/commands/execute-task.md`
3. **Archive Task**: Move completed tasks via `.claude/commands/archive.md`

## Directory Structure
```
.claude/
  orchestrator.md          # This file - main task registry
  commands/
    new-task.md            # How to create new tasks
    execute-task.md        # How to execute tasks with agents
    archive.md             # How to archive completed tasks
  tasks/
    active/                # Active task files
    archived/              # Completed task files
```

## Active Tasks
See `tasks/active/` directory for individual task files.

## Conventions
- Task IDs: `TASK-NNN` (zero-padded 3 digits)
- Priority levels: `critical`, `high`, `medium`, `low`
- Each task has its own `.md` file in `tasks/active/`
