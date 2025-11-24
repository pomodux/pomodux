---
status: approved
---

# Session Label Default Behavior

## 1. Context / Background

### 1.1 Problem Statement

Every timer session needs a label for history tracking and session identification. The question is whether labels should be:
- **Strictly required** (user must provide)
- **Optional with default** (fallback if omitted)
- **Optional with no default** (empty string allowed)

This decision affects command ergonomics and session data quality.

### 1.2 User Context

**Use Cases:**
- **Quick Timers**: User wants to start timer ASAP without typing label
- **Tracked Work**: User wants meaningful labels for time tracking
- **Retrospective Review**: User looks at history to understand past work
- **Client Billing**: Labels help identify billable time

**User Behavior:**
- Some users want detailed labels
- Some users prefer minimal friction
- Some users establish labeling conventions
- Some users forget to add labels

## 2. Decision

**Selected Solution:** Option B - Optional with smart defaults

### 2.1 Behavior

**With Preset:**
```bash
pomodux start work
# Label: "Work"  (prettified preset name)

pomodux start work "Implementing auth"
# Label: "Implementing auth"  (user-provided)

pomodux start longbreak
# Label: "Long Break"  (prettified preset name)
```

**With Custom Duration:**
```bash
pomodux start 45m
# Label: "Generic timer session"  (generic default)

pomodux start 45m "Client meeting"
# Label: "Client meeting"  (user-provided)
```

### 2.2 Prettification Rules

**Preset Name → Default Label:**
- `work` → `"Work"`
- `break` → `"Break"`
- `longbreak` → `"Long Break"`
- `meeting` → `"Meeting"`
- `custom_preset` → `"Custom Preset"`

**Algorithm:**
```
1. Split on underscores/hyphens
2. Capitalize first letter of each word
3. Join with spaces
```

### 2.3 Rationale

**Why Allow Omitting Label?**

1. **Ergonomic**: Quick start without typing
2. **Friction Reduction**: Fewer keystrokes for common timers
3. **User Choice**: Power users can add labels, casual users can skip
4. **Workflow Speed**: Fast timer chains

**Why Provide Default?**

1. **History Quality**: Every session has some context
2. **Retrospective Value**: Can distinguish preset types in history
3. **Not Empty**: No confusing blank labels in stats
4. **Searchable**: Can filter by preset type later

**Why Prettify Preset Name?**

1. **Human Readable**: "Work" > "work"
2. **Professional**: Better for client-facing reports
3. **Conventional**: Multi-word presets look normal ("Long Break")
4. **Minimal Surprise**: Matches user's mental model of preset

## 3. Alternatives Considered

### 3.1 Option A: Strictly Required Labels

**Approach:**
```bash
pomodux start work
# Error: Label required
# Usage: pomodux start work <label>

pomodux start work "My task"
# OK
```

**Pros:**
- Forces deliberate labeling
- High-quality session history
- Clear intent for every session
- Better for time tracking
- No ambiguous sessions

**Cons:**
- Friction for quick timers
- Annoying for simple breaks
- Breaks workflow for rapid starts
- Users might use placeholder labels ("foo", "timer", "work")
- Defeats purpose of presets (why preset if still typing?)

**Rejected:** Too much friction, reduces ergonomics

---

### 3.2 Option B: Optional with Smart Defaults (Selected)

**Approach:**
```bash
pomodux start work
# Label: "Work" (from preset)

pomodux start 25m
# Label: "Generic timer session"

pomodux start work "Real task"
# Label: "Real task"
```

**Pros:**
- Fast quick starts
- Meaningful defaults from presets
- User choice (can add label)
- No empty labels
- Good history quality
- Ergonomic and professional

**Cons:**
- Generic default for custom durations
- Users might rely on defaults too much
- Preset-based sessions harder to distinguish without label

**Selected:** Best balance of ergonomics and data quality

---

### 3.3 Option C: Optional with No Default

**Approach:**
```bash
pomodux start work
# Label: "" (empty string)

pomodux start work "My task"
# Label: "My task"
```

**Pros:**
- Simple implementation
- User fully controls labels
- No "magic" behavior

**Cons:**
- Empty labels in history
- Confusing stats ("5 sessions labeled: ''")
- Hard to distinguish sessions
- Poor user experience reviewing history
- Wasted opportunity (preset name available)

**Rejected:** Poor history quality, confusing UX

---

### 3.4 Option D: Interactive Prompt

**Approach:**
```bash
pomodux start work
# Prompts: "Enter label (default: Work): _"
# User can press Enter for default or type custom
```

**Pros:**
- Explicit choice every time
- Encourages thoughtful labeling
- Clear what default will be

**Cons:**
- Interrupts flow
- Extra interaction required
- Annoying for quick timers
- Breaks scriptability
- Against Unix philosophy (non-interactive CLI)

**Rejected:** Too interactive, breaks scripting

## 4. Consequences

### 4.1 Positive

**User Experience:**
- Quick starts without typing
- Meaningful defaults from presets
- Freedom to add detailed labels when needed
- No empty labels in history

**History Quality:**
- All sessions labeled
- Preset type visible in history
- Can filter by preset ("show all Work sessions")
- Professional appearance

**Ergonomics:**
- Fast workflow: `pomodux start work`
- Still precise when needed: `pomodux start work "Bug #123"`
- Natural progressive enhancement

### 4.2 Negative

**Generic Custom Timers:**
- Custom durations get generic label if omitted
- History might have many "Generic timer session" entries
- Hard to distinguish these sessions later

**Potential Over-Reliance:**
- Users might always omit labels
- History less detailed than it could be
- Harder to track specific tasks

**Mitigations:**
- Documentation encourages labeling for tracked work
- Stats can show "unlabeled" session count
- Future: Prompt after N generic sessions?

### 5.3 Session History Example

```
Recent Sessions
─────────────────────────────────────────────────────────
Start Time          Duration  Status      Label
─────────────────────────────────────────────────────────
2025-01-15 14:00    25m       completed   Auth module refactoring
2025-01-15 14:30    5m        completed   Break
2025-01-15 14:40    25m       completed   Work
2025-01-15 15:10    45m       stopped     Generic timer session
2025-01-15 16:00    15m       completed   Long Break
─────────────────────────────────────────────────────────
```

Notice:
- Custom labels are detailed
- Preset defaults are clear ("Work", "Break")
- Generic default appears for custom duration without label
- Mix of labeling styles based on user choice

## 6. User Guidance

### 6.1 Best Practices Documentation

**README:**
```markdown
### Session Labels

Labels are optional but recommended for work tracking:

\`\`\`bash
# Quick start with preset (uses "Work" as label)
pomodux start work

# Detailed label for specific task
pomodux start work "Implement user authentication"

# Custom duration without label (uses "Generic timer session")
pomodux start 45m

# Custom duration with label
pomodux start 45m "Client meeting with ACME Corp"
\`\`\`

**Tip:** Add labels for work you'll bill or want to track specifically.
```

### 6.2 Config File Comments

```yaml
# Timer presets
# When you start a timer with a preset but no label,
# the preset name is prettified and used as the label.
# Example: "work" → "Work", "long_break" → "Long Break"
timers:
  work: 25m
  break: 5m
  long_break: 15m
```

### 6.3 CLI Help

```
EXAMPLES:
  # Start work session (label: "Work")
  pomodux start work

  # Start with specific label
  pomodux start work "Implement login feature"

  # Custom duration without label (label: "Generic timer session")
  pomodux start 45m

  # Custom duration with label
  pomodux start 1h30m "Team meeting"
```

## 7. Future Enhancements

### 7.1 Smart Label Suggestions

**Potential:**
```bash
pomodux start work
# After 3 generic "Work" sessions:
# Tip: Consider adding specific labels for better tracking
# Example: pomodux start work "Your task description"
```

**Evaluation:**
- Track frequency of default labels
- Only suggest if pattern detected
- Don't nag power users who prefer defaults

### 7.2 Label Templates

**Potential:**
```yaml
timers:
  work:
    duration: 25m
    default_label_template: "Work - {prompt}"
    # Prompts: "Work - _" on start
```

**Concerns:**
- Adds complexity
- Interactive prompts break scriptability
- Better to encourage manual labels
- Not needed if docs are clear

### 7.3 Session Tagging (Post-MVP)

**Instead of complex labels:**
```bash
pomodux start work "Auth" --tags project:backend,priority:high
```

**Benefits:**
- Better than cramming tags into labels
- Structured metadata
- Better filtering/reporting

**Deferred:** Post-MVP feature

