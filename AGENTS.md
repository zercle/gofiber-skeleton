---
description: Agentic AI Coder - Optimized Ruleset
globs:
  - "**/*"
alwaysApply: true
---

## Core Philosophy
Spec-Driven Development emphasizes intent-driven development where specifications define the "what" before the "how." Uses multi-step refinement rather than one-shot generation, with Memory Bank providing persistent context across sessions.

## Memory Bank Structure
Location: `${workspaceFolder}/.agents/rules/memory-bank`

### Critical Workflow Rules
1. **MANDATORY**: Read ALL memory bank files at start of EVERY task
2. Indicate status with `[Memory Bank: Active]` or `[Memory Bank: Missing]`
3. Briefly summarize project understanding for confirmation
4. For significant changes: suggest updating memory bank
5. At task completion: update `context.md` with current state

### Core Files (Required)
1. `brief.md` - Foundation document (user-maintained, don't edit directly)
   - Core requirements and goals
   - Source of truth for project scope
2. `product.md` - Problem statement, UX goals, user personas
3. `context.md` - Current work focus, recent changes, next steps, development phase
4. `architecture.md` - System structure, component relationships, ADR format for decisions
5. `tech.md` - Technology stack, dependencies, setup requirements

### Optional Files
- `specs.md` - Feature specifications, acceptance criteria, business rules
- `experiments.md` - Alternative implementations and comparisons
- `tasks.md` - Repetitive task workflows
- Additional specialized files as needed

## Development Phases
1. **0-to-1 Development**: Generate from scratch using high-level requirements
2. **Creative Exploration**: Parallel implementations and diverse solutions
3. **Iterative Enhancement**: Feature additions and legacy modernization

## Key Workflows

### Memory Bank Initialization (CRITICAL)
When user requests `initialize memory bank`:
- Perform exhaustive analysis of all project files
- Create comprehensive initial memory bank
- Each memory bank core file limit to 300-500 lines
- Provide summary for user verification
- Emphasize this foundation affects all future interactions

### Memory Bank Update
Triggered by:
1. User explicitly requests `update memory bank` (MUST review ALL files)
2. Discovering new project patterns
3. After implementing significant changes
4. When context needs clarification
5. When moving between development phases

Process:
1. Review ALL project files
2. Document current state and insights
3. If requested with context source (e.g., "@/Makefile"), focus special attention
4. Ensure alignment between specs.md and implementation

### Task Documentation
When user requests `add task` or for repetitive workflows:
- Create/update `tasks.md` in memory bank folder
- Document with:
  - Task name and description
  - Files to modify
  - Step-by-step workflow
  - Important considerations and gotchas
  - Example implementation
- Suggest documentation for potentially recurring tasks

## Implementation Process
1. **Specification Creation**:
   - Start with requirements in brief.md
   - Expand details in specs.md
   - Define acceptance criteria
   - Get confirmation before implementation

2. **Implementation Planning**:
   - Review specs.md and brief.md
   - Plan steps in context.md
   - Identify files to modify based on architecture.md
   - Document new technical decisions in architecture.md

3. **Validation and Refinement**:
   - Validate against specs.md requirements
   - Update context.md with current state
   - Document any deviations from specifications
   - Suggest memory bank update for significant changes

## Context Window Management
When context window fills:
- Suggest updating memory bank to preserve state
- Recommend starting fresh conversation
- New sessions automatically reload memory bank for continuity

## Priority Rules
- Prioritize brief.md over conflicting information
- If inconsistencies detected: note discrepancies to user
- Memory Bank is only link between sessions - maintain precision