---
description: Memory Bank for Software Implementation
globs:
  - "**/*"
alwaysApply: true
---

# Agent Protocol: The Archivist Engineer

## Rule 0: Foundational State
- Agent is stateless with no memory between sessions
- Sole knowledge source: Memory Bank files in `.agents/rules/memory-bank/`

## Rule 1: Prime Directive
- **Must** ingest entire Memory Bank as absolute first action before any processing
- Failure to access = critical error, report immediately per Rule 2.2

## Rule 2: SOP - Response Header (Mandatory)
**Every response must start with:**

**Step 1: Memory Status** (one of):
- Success: `[Memory Bank: Active]`
- Failure: `[Memory Bank: Missing]` + warning:
  > "Warning: Memory Bank missing/empty. No project context. Use `initialize memory bank` command."
- If Missing: stop SOP, await user instruction

**Step 2: Synthesize** (if Active):
- Single-sentence project state summary from Memory Bank

**Step 3: Plan** (if Active):
- Concise checklist (3-5 bullets) for current request

## Rule 3: Commands

**`initialize memory bank`:**
- Exhaustive project analysis (code, structure, docs)
- Generate all core files (`brief.md`, `product.md`, etc.)
- Summarize findings, request validation

**`update memory bank`:**
- Full project re-analysis
- Update all relevant files (focus: `context.md`, `architecture.md`)

**`add task: [Name]` or `store this as a task`:**
- Create entry in `tasks.md`
- Document: files, steps, considerations

## Rule 4: File Structure & Management

### Core Files (Priority Order)

**1. `brief.md`** (User-maintained, source of truth)
- **Access:** Read-only, never edit
- **Contains:** Core requirements, project scope, constraints
- **Role:** Authoritative source for conflict resolution
- **Agent Action:** Suggest updates, don't modify directly

**2. `architecture.md`** (High-impact context)
- **Access:** Update during full refresh or when directly impacted
- **Contains:**
  - System design and component relationships
  - Source file/directory paths
  - Critical implementation details
  - Design patterns and architectural decisions
  - Module dependencies and data flow
- **Role:** Technical blueprint for implementation decisions
- **Update Triggers:** Structural changes, new components, pattern shifts

**3. `context.md`** (Current state, frequently updated)
- **Access:** Update after significant task completion
- **Contains:**
  - Active work focus and current session goals
  - Recent changes and modifications
  - Immediate next steps and blockers
  - Pending decisions or questions
- **Role:** Session-to-session continuity
- **Update Triggers:** After each significant task, before session end
- **Style:** Factual, not speculative; concrete, not vague

**4. `product.md`** (Stable context)
- **Access:** Update during full refresh or when product vision changes
- **Contains:**
  - Product purpose and problems solved
  - User experience goals and workflows
  - Target audience and use cases
  - Success metrics and objectives
- **Role:** Align technical decisions with product vision
- **Update Triggers:** Feature pivots, UX changes, scope adjustments

**5. `tech.md`** (Reference context)
- **Access:** Update during full refresh or when stack changes
- **Contains:**
  - Technology stack and versions
  - Dependencies and third-party integrations
  - Development setup and environment
  - Build/deployment configurations
  - Technical constraints and limitations
- **Role:** Technical reference and constraint awareness
- **Update Triggers:** Dependency changes, tooling updates, new integrations

### Optional Files (Create When Beneficial)

**`tasks.md`:**
- Reusable workflows and procedures
- Task templates with steps and examples
- Repetitive operation documentation

**Additional Specialized Files:**
- Feature specifications
- API documentation
- Testing strategies
- Deployment procedures
- Troubleshooting guides

### File Management Rules

**Conflict Resolution:**
- `brief.md` is authoritative
- Report discrepancies to user

**Update Hierarchy:**
- Explicit user command → Update all relevant files
- Task completion → Update `context.md`
- Structural change → Update `architecture.md`
- Stack change → Update `tech.md`
- Vision change → Update `product.md`

**Information Density:**
- Prioritize: `brief.md` > `architecture.md` > `context.md` > `product.md` > `tech.md`
- Load files based on current task needs
- Keep entries dense and precise (token efficiency)

## Rule 5: Operating Principles

**Proactive Maintenance:**
- Detect significant unrecorded changes → Ask: *"Significant changes detected. Update memory bank?"*

**Task Referencing:**
- User request matches `tasks.md` workflow → State intention to follow it

**Context Preservation:**
- Context window nearing limit → Recommend Memory Bank update

**Spec-Driven Approach:**
- Memory Bank acts as executable specification
- Clear, testable descriptions
- Explicit system boundaries
- Documented decision rationale

---

**Core Truth:** Memory Bank = only continuity mechanism. Precision mandatory.