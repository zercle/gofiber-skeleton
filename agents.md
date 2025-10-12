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
- **Critical**: When uncertain about Memory Bank status, attempt to read files first
- Default assumption: if `.agents/rules/memory-bank/` directory exists, Memory Bank is available

## Rule 1: Prime Directive - Memory Bank Discovery
**First action sequence:**
1. **Attempt to read** core Memory Bank files from `.agents/rules/memory-bank/`:
   - `brief.md` (highest priority - user requirements)
   - `architecture.md` (system design and structure)
   - `context.md` (current state and session continuity)
2. **Count readable files** (any .md files in memory-bank directory)
3. **Determine status** based on actual file access, not directory existence

## Rule 2: SOP - Response Header (Mandatory)
**Every response must start with:**

**Step 1: Memory Status** (determined by actual file reads):
- Success: `[Memory Bank: Active - {X} files loaded]` (1+ .md files successfully read)
- Empty: `[Memory Bank: Empty - Directory exists, no readable .md files]`
- Missing: `[Memory Bank: Not Found - Directory does not exist]`
- Fallback: `[Memory Bank: Access Issue - Proceeding with available context]`

**Step 2a: If Memory Bank Missing/Empty:**
- State: "No existing project context found in Memory Bank."
- Offer: "I can initialize the Memory Bank by analyzing your project structure and requirements. Would you like me to run `initialize memory bank`?"
- **Continue processing** the user's request using available context (visible code files, user input, project structure)

**Step 2b: If Memory Bank Active:**
- **Synthesize**: One-sentence project state summary from Memory Bank content
- **Highlight**: Current focus area from `context.md` (if available)
- **Plan**: Concise action checklist (3-5 bullets) for current request based on Memory Bank insights

**Step 2c: If Access Issue:**
- State: "Memory Bank directory detected but files not immediately accessible."
- **Continue processing** request with available context
- Note: "Will attempt Memory Bank integration as files become available."

## Rule 3: Enhanced Commands

**`initialize memory bank`:**
- Create `.agents/rules/memory-bank/` directory structure
- Perform comprehensive project analysis:
  - Code structure and file organization
  - Dependencies and technology stack
  - User requirements from conversation history
  - Existing documentation and comments
- Generate complete Memory Bank with all core files
- **Validation step**: Present findings summary and ask for user confirmation
- Confirm Memory Bank activation with file count

**`update memory bank`:**
- Re-analyze project state against current Memory Bank
- **Smart updates**: Only modify files with detected changes
- Priority order: `context.md` → `architecture.md` → `tech.md` → others
- Report specific changes made and reasoning
- **Incremental approach**: Preserve user customizations in `brief.md`

**`refresh context`:**
- Quick update focused on `context.md` only
- Capture recent changes, current session progress
- Update next steps and blockers
- Faster alternative to full Memory Bank update

**`add task: [Name]` or `store this as task`:**
- Append to `tasks.md` with structured format:
  - Task name and description
  - Affected files and components
  - Step-by-step procedure
  - Dependencies and considerations

## Rule 4: File Structure & Management

### Core Files (Priority Order)

**1. `brief.md`** (Sacred - User Authority)
- **Access:** Read-only for agent, user-maintained
- **Contains:** Authoritative requirements, project scope, constraints, success criteria
- **Agent Role:** Reference for all decisions, suggest updates via comments, never edit directly
- **Conflict Resolution:** This file wins all disputes

**2. `architecture.md`** (Technical Blueprint)
- **Update Triggers:** Structural changes, new components, design pattern shifts
- **Contains:**
  - System architecture and component relationships
  - File/directory structure with paths
  - Core implementation patterns and conventions
  - Module dependencies and data flow diagrams
  - Integration points and external interfaces
- **Maintenance:** Keep synchronized with actual codebase structure

**3. `context.md`** (Session Continuity)
- **Update Triggers:** After significant task completion, session transitions
- **Contains:**
  - Current work focus and active development area
  - Recent changes and their impact
  - Immediate next steps and priorities
  - Open questions and pending decisions
  - Session-specific notes and discoveries
- **Role:** Bridge between sessions, maintain momentum

**4. `product.md`** (User Experience)
- **Update Triggers:** Feature changes, UX pivots, user feedback integration
- **Contains:** Product vision, user workflows, feature specifications, success metrics
- **Stability:** Changes less frequently than technical files

**5. `tech.md`** (Technical Environment)
- **Update Triggers:** Dependency changes, tooling updates, environment shifts
- **Contains:** Technology stack, build systems, deployment processes, constraints
- **Role:** Technical decision reference

### Optional Specialized Files
- `tasks.md`: Reusable procedures and workflows
- `features.md`: Feature specifications and status
- `api.md`: API documentation and contracts
- `testing.md`: Testing strategies and procedures
- `deployment.md`: Deployment processes and environments

### File Management Principles
- **Atomic Updates**: Complete file rewrites, not patches
- **Size Limits**: Target 300-500 lines per file for optimal context window usage
- **Density**: Prefer structured, scannable content over prose
- **Consistency**: Maintain consistent formatting and structure across files
- **Auto-Recovery**: Recreate missing core files during updates when possible

## Rule 5: Enhanced Operating Principles

### Graceful Degradation Strategy
- **Missing Memory Bank**: Offer initialization, continue with available context
- **Partial Memory Bank**: Use available files, note missing components
- **Access Issues**: Work with visible project files, attempt Memory Bank integration later
- **Always Deliver Value**: Provide useful output regardless of Memory Bank status

### Proactive Context Management
- **Change Detection**: Monitor for significant architectural or requirement changes
- **Update Recommendations**: Suggest Memory Bank updates when context drift detected
- **Context Window Optimization**: Recommend Memory Bank refresh when approaching limits
- **Smart Loading**: Prioritize most relevant Memory Bank files for current task

### Performance Optimizations
- **Selective Reading**: Load only relevant Memory Bank files for focused tasks
- **Incremental Updates**: Update only changed sections during Memory Bank maintenance
- **Context Prioritization**: Weight recent changes and current focus areas higher
- **Batch Operations**: Combine multiple small updates into single Memory Bank refresh

### Error Recovery
- **File Corruption**: Attempt reconstruction from project analysis
- **Missing Dependencies**: Document in `context.md` for user resolution
- **Conflicting Information**: Defer to `brief.md`, document conflicts for user review
- **Incomplete Initialization**: Support partial Memory Bank creation and gradual completion

## Rule 6: Quality Assurance

### Memory Bank Validation
- **Consistency Checks**: Ensure alignment between Memory Bank files and actual codebase
- **Completeness Verification**: Identify gaps in project coverage
- **Accuracy Monitoring**: Flag potential outdated information
- **User Feedback Integration**: Incorporate corrections and clarifications

### Continuous Improvement
- **Learning from Sessions**: Capture patterns in successful project interactions
- **Template Refinement**: Improve Memory Bank file templates based on usage
- **Process Optimization**: Streamline common workflows and procedures
- **User Experience**: Minimize friction in Memory Bank interactions

---

**Core Philosophy:** Memory Bank enables intelligent, context-aware assistance. When available, leverage it fully. When missing or incomplete, build it collaboratively while delivering immediate value.