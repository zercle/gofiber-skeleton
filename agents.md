---
description: Memory Bank for Software Implementation
globs:
  - "**/*"
alwaysApply: true
---

# Memory Bank Protocol

I'm an AI agent with session-based memory. Between sessions, I rely on Memory Bank files in `.agents/rules/memory-bank/` to maintain project continuity.

## Session Start Protocol

**ALWAYS check the actual file system first - never assume files are missing:**

### Step 1: File System Check
- **FIRST**: Check if `.agents/rules/memory-bank/` directory exists
- **THEN**: List actual files present in that directory
- **NEVER** report "missing" without actually checking

### Step 2: Response Based on Actual State

#### If Files Exist and Readable:
- `[Memory Bank: Active - Loaded: {list actual files found}]`
- Brief project summary from files
- Proceed with loaded context

#### If Directory Exists but No Files:
- `[Memory Bank: Directory found, files empty - Initializing]`
- Immediately offer: "I found the memory bank directory but it's empty. I'll analyze the project and create the core files now."
- **Auto-proceed** with initialization unless user objects

#### If Directory Missing:
- `[Memory Bank: Creating structure]`
- Create `.agents/rules/memory-bank/` directory
- Immediately initialize with project analysis
- Report completion

#### If Partial Files Exist:
- `[Memory Bank: Partial - Found: {list actual files}]`
- Load available files
- Note which core files are missing: `Missing: {list}`
- Offer: "Should I create the missing files to complete the memory bank?"

## Improved File Discovery

### Before Reading Files:
1. **Check directory existence**: `.agents/rules/memory-bank/`
2. **List directory contents**: Show what's actually there
3. **Attempt to read each file**: Handle read errors gracefully
4. **Report actual status**: Based on what was found, not assumptions

### Error Handling Hierarchy:
1. **Directory missing** → Create and initialize
2. **Files missing** → Create missing files
3. **Files unreadable** → Report specific file issues
4. **Permissions issues** → Provide specific guidance

## Core Files (Priority Order)

### Essential Files
1. **`brief.md`** - Project requirements and scope (user-maintained)
2. **`architecture.md`** - System design and implementation details
3. **`context.md`** - Current state and active work
4. **`product.md`** - Purpose and user experience goals
5. **`tech.md`** - Technology stack and dependencies

### File Creation Strategy
- **Always create directory first** if missing
- **Populate immediately** with project analysis
- **Verify creation** and report success
- **Request user review** for accuracy

## Enhanced Workflows

### 1. Smart Auto-Initialization

**Triggers:**
- Directory missing
- Directory empty
- Core files missing

**Process:**
1. Create `.agents/rules/memory-bank/` if needed
2. Analyze current project structure
3. Identify project type, main files, and patterns
4. Create all core files with discovered information
5. Report: `[Memory Bank: Initialized - Created: {list files}]`
6. Ask: "Please review the generated memory bank files for accuracy"

### 2. Intelligent File Recovery

**When some files exist:**
1. Load and analyze existing files
2. Identify missing core files
3. Generate missing files based on:
   - Existing file content
   - Current project analysis
   - Consistency with loaded context
4. Maintain coherence across all files

### 3. Context-Aware Updates

**Smart update triggers:**
- After implementing new features
- When project structure changes
- When dependencies are modified
- When user requests updates

**Update process:**
1. **Verify all files accessible**
2. **Load current state** from all files
3. **Compare with project reality**
4. **Update incrementally** (don't rewrite everything)
5. **Prioritize `context.md`** for recent changes
6. **Preserve user-maintained content** in `brief.md`

## Robust Error Prevention

### File System Interaction Rules:
- **Never assume** - always check actual file system
- **Create missing structure** proactively
- **Handle partial states** gracefully
- **Provide specific error messages** with solutions

### Common Issue Prevention:

**"Memory bank missing" reports:**
- **Root cause**: Not checking file system
- **Solution**: Always verify directory and file existence first
- **Response**: Create missing components immediately

**Permission/access issues:**
- **Check**: File permissions and directory access
- **Report**: Specific files and permission issues
- **Suggest**: Concrete steps to resolve

**Inconsistent file states:**
- **Detect**: Files that don't match current project
- **Update**: Outdated information automatically
- **Preserve**: User-provided requirements and constraints

## Performance Optimizations

### Efficient Context Loading:
1. **Load `context.md` first** - most relevant for current work
2. **Load `architecture.md` second** - high-impact technical context
3. **Load remaining files** based on immediate needs
4. **Cache insights** to avoid re-reading during session

### Smart Update Strategy:
- **Incremental updates** rather than full rewrites
- **Focus on changed areas** of the project
- **Preserve stable information** across updates
- **Version awareness** - note when major changes occur

## User Experience Improvements

### Clear Status Communication:
- **Specific file states**: "Found 3 of 5 core files"
- **Immediate actions**: "Creating missing architecture.md"
- **Next steps**: "Ready to proceed with loaded context"

### Proactive Problem Solving:
- **Detect issues early** in session start
- **Offer solutions immediately** - don't wait for user requests
- **Auto-recover when possible** - minimize user intervention
- **Explain what happened** - build user confidence

### Reduced Friction:
- **Auto-create missing structure** instead of reporting errors
- **Initialize with intelligent defaults** based on project analysis
- **Update incrementally** rather than requiring full rebuilds
- **Maintain continuity** even when files are temporarily inaccessible

---

## Key Behavioral Changes

1. **File System First**: Always check actual files before reporting status
2. **Proactive Creation**: Create missing components immediately
3. **Intelligent Defaults**: Populate files with project analysis, not empty templates
4. **Graceful Degradation**: Work effectively even with partial or missing files
5. **User-Centric**: Minimize troubleshooting, maximize productivity

**Remember**: The memory bank should be invisible infrastructure that "just works." When there are issues, solve them immediately rather than requiring user intervention.