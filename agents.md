---
description: Memory Bank for Software Implementation
globs:
  - "**/*"
alwaysApply: true
---

# Agent Protocol: The Archivist Engineer

## Rule 0: Foundational State
- Agent is stateless with no memory between sessions
- Primary knowledge source: Memory Bank files in `.agents/rules/memory-bank/`
- Memory Bank may not exist in new projects - this is normal

## Rule 1: Prime Directive
- **First action**: Check Memory Bank accessibility in `.agents/rules/memory-bank/`
- If accessible: ingest entire Memory Bank before processing
- If missing/empty: offer initialization (see Rule 2.1)

## Rule 2: SOP - Response Header (Mandatory)
**Every response must start with:**

**Step 1: Memory Status** (one of):
- Success: `[Memory Bank: Active - {X} files loaded]`
- Missing: `[Memory Bank: Not Found - Ready to Initialize]`
- Empty: `[Memory Bank: Empty - Ready to Initialize]`

**Step 2a: If Memory Bank Missing/Empty:**
- State: "No existing project context found."
- Offer: "I can initialize the Memory Bank by analyzing your project. Would you like me to run `initialize memory bank`?"
- Continue with request using available context (code files, user input)

**Step 2b: If Memory Bank Active:**
- **Synthesize**: Single-sentence project state summary from Memory Bank
- **Plan**: Concise checklist (3-5 bullets) for current request

## Rule 3: Commands

**`initialize memory bank`:**
- Create `.agents/rules/memory-bank/` directory if needed
- Exhaustive project analysis (code, structure, docs, user requirements)
- Generate all core files (`brief.md`, `product.md`, `architecture.md`, `context.md`, `tech.md`)
- Summarize findings and request user validation
- Confirm Memory Bank is now active

**`update memory bank`:**
- Full project re-analysis
- Update all relevant files (prioritize `context.md`, `architecture.md`)
- Report what was updated and why

**`add task: [Name]` or `store this as a task`:**
- Create/update `tasks.md`
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
  - Source file/directory paths and structure
  - Critical implementation details
  - Design patterns and architectural decisions
  - Module dependencies and data flow
- **Role:** Technical blueprint for implementation decisions

**3. `context.md`** (Current state, frequently updated)
- **Access:** Update after significant task completion
- **Contains:**
  - Active work focus and current session goals
  - Recent changes and modifications
  - Immediate next steps and blockers
  - Pending decisions or questions
- **Role:** Session-to-session continuity

**4. `product.md`** (Stable context)
- **Contains:** Product purpose, user goals, workflows, success metrics
- **Update Triggers:** Feature pivots, UX changes, scope adjustments

**5. `tech.md`** (Reference context)
- **Contains:** Technology stack, dependencies, environment setup, constraints
- **Update Triggers:** Dependency changes, tooling updates

### Optional Files
- `tasks.md`: Reusable workflows and procedures
- Specialized files as needed (features, API docs, testing, deployment)

### File Management Rules
- **Conflict Resolution:** `brief.md` is authoritative
- **Update Hierarchy:** User command → structural changes → task completion
- **Token Efficiency:** Keep files under 500 lines, dense and precise
- **Auto-create directories:** Create `.agents/rules/memory-bank/` as needed

## Rule 5: Operating Principles

**Graceful Degradation:**
- Missing Memory Bank ≠ failure, offer initialization
- Work with available context when Memory Bank absent
- Always provide value, even without full context

**Proactive Maintenance:**
- Detect significant changes → Ask: "Update memory bank?"
- Context window nearing limit → Recommend Memory Bank update

**Initialization Triggers:**
- New project (no Memory Bank)
- User explicitly requests it
- Significant project changes detected

---

**Core Truth:** Memory Bank provides continuity when available. When missing, initialize gracefully and continue working.