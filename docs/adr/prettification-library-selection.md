---
status: approved
---

# ADR 007: Preset Name Prettification with strcase

## 1. Context / Background

### 1.1 Problem Statement

Pomodux allows users to define timer presets in their config file (e.g., `work`, `break`, `long_break`). When a user starts a timer without providing an explicit label, the preset name should be used as the session label. However, preset names follow programming conventions (lowercase, underscores) which don't look polished in session history and statistics displays.

**Examples of the problem:**
```bash
# User runs:
pomodux start work

# Session recorded with label: "work" (not pretty)
# Desired label: "Work"

# User runs:
pomodux start long_break

# Session recorded with label: "long_break" (not pretty)
# Desired label: "Long Break"
```

### 1.2 Requirements

- **Automatic Prettification**: Convert preset names to human-readable labels
- **Convention Support**: Handle common naming conventions (snake_case, camelCase, kebab-case)
- **Consistency**: Predictable output for users
- **No Magic**: Simple, understandable transformation rules
- **User Override**: Users can always provide explicit labels to override prettification
- **Minimal Dependencies**: Prefer standard library, but accept well-maintained packages if needed

### 1.3 Design Principles

Following Go best practices, preset names should use proper conventions:
- `work` - single word (lowercase)
- `long_break` - compound word (snake_case)
- `short_break` - compound word (snake_case)
- `code_review` - compound word (snake_case)

Prettification should then handle the formatting for display purposes.

## 2. Decision

**Selected Solution:** Use `github.com/iancoleman/strcase` for preset name prettification

### 2.1 Rationale

**Why strcase?**
1. **Well-Maintained**: 2.3k+ GitHub stars, actively maintained, last updated recently
2. **Zero Dependencies**: Pure Go implementation, no transitive dependencies
3. **Comprehensive**: Handles all common case conventions (snake_case, camelCase, kebab-case, etc.)
4. **Battle-Tested**: Widely used in Go ecosystem
5. **Simple API**: Easy to use and understand
6. **Proper Title Casing**: Converts to space-delimited and handles capitalization correctly

**Why not alternatives?**
- **Standard library `strings.Title()`**: Deprecated in Go 1.18+, doesn't split words
- **`golang.org/x/text/cases`**: Doesn't handle word splitting (e.g., "longbreak" → "Longbreak" not "Long Break")
- **Custom implementation**: Reinventing the wheel, more maintenance burden, edge cases

### 2.2 Usage Examples

**In Configuration:**
```yaml
timers:
  work: 25m
  break: 5m
  long_break: 15m
  code_review: 15m
  deep_focus: 50m
```

**Prettification Results:**
```
work         → "Work"
break        → "Break"
long_break   → "Long Break"
code_review  → "Code Review"
deep_focus   → "Deep Focus"
```

**In Practice:**
```bash
# Without explicit label - uses prettified preset name
$ pomodux start long_break
→ Session label: "Long Break"

# With explicit label - user override
$ pomodux start long_break "Coffee and email"
→ Session label: "Coffee and email"

# Custom duration without label - uses generic fallback
$ pomodux start 45m
→ Session label: "Generic timer session"
```

## 3. Alternatives Considered

### 3.1 Standard Library Only (golang.org/x/text/cases)

**Approach:**
```go
import "golang.org/x/text/cases"
import "golang.org/x/text/language"

caser := cases.Title(language.English)
caser.String("work")      // "Work"
caser.String("long_break") // "Long_break" ❌
```

**Pros:**
- No external dependencies
- Official Go extended package

**Cons:**
- Doesn't handle word splitting
- Produces incorrect output: "Long_break" instead of "Long Break"
- Would require custom splitting logic anyway

**Rejected:** Insufficient functionality for requirements

---

### 3.2 Custom Implementation

**Approach:**
```go
func PrettifyPresetName(preset string) string {
    // Manual word splitting
    words := strings.Split(preset, "_")
    for i, word := range words {
        if len(word) > 0 {
            words[i] = strings.ToUpper(word[0:1]) + word[1:]
        }
    }
    return strings.Join(words, " ")
}
```

**Pros:**
- No dependencies
- Full control over behavior

**Cons:**
- Only handles snake_case (not camelCase, kebab-case, etc.)
- Requires significant testing for edge cases
- Maintenance burden
- Doesn't handle special cases (acronyms, etc.)

**Rejected:** Reinventing the wheel, limited functionality

---

### 3.3 github.com/stoewer/go-strcase

**Approach:**
Similar to iancoleman/strcase

**Pros:**
- Well-maintained
- Similar functionality

**Cons:**
- Less popular (fewer stars/usage)
- Less comprehensive documentation
- Smaller community

**Rejected:** iancoleman/strcase is more widely adopted

---

### 3.4 Configuration-Based Prettification

**Approach:**
Allow users to define display names in config:

```yaml
timers:
  long_break:
    duration: 15m
    display_name: "Long Break"
```

**Pros:**
- Maximum user control
- No dependencies needed

**Cons:**
- More verbose configuration
- Users must define display names for every preset
- Breaks simple config format
- Still need fallback for missing display names

**Rejected:** Over-engineered for MVP, can add later if needed

## 4. Consequences

### 4.1 Positive

- **User Experience**: Session labels look polished and professional
- **Convention Support**: Users can follow Go naming conventions in config
- **Consistency**: Predictable output across all case conventions
- **Maintainability**: Well-tested library handles edge cases
- **Flexibility**: Users can override with explicit labels anytime

### 4.2 Negative

- **External Dependency**: Adds one external package to dependency tree
- **Binary Size**: Minimal increase (~10-20KB)
- **Learning Curve**: Developers need to understand strcase behavior

### 4.3 Risks and Mitigations

**Risk: Library becomes unmaintained**
- **Likelihood**: Low (2.3k stars, active development)
- **Impact**: Low (simple, stable functionality)
- **Mitigation**: Library code is small and simple, could vendor or fork if needed

**Risk: Unexpected prettification behavior**
- **Likelihood**: Medium (edge cases with unusual preset names)
- **Impact**: Low (user can always provide explicit label)
- **Mitigation**: Document expected behavior, comprehensive testing

**Risk: Dependency vulnerability**
- **Likelihood**: Very low (zero dependencies, minimal attack surface)
- **Impact**: Low (display formatting only, no security implications)
- **Mitigation**: Regular dependency scanning, simple to replace if needed

