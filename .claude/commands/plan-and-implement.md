Read through the full prompt. Use the provided pattern steps to implement the following feature:

$ARGUMENTS

**CRITICAL**: This pattern contains mandatory approval gates that MUST be respected:

- **NEVER** skip phases or approval gates
- **ALWAYS** wait for explicit approval before proceeding  
- **STOP** implementation if approval is not granted
- **FORMAT** all "Report back" sections clearly and comprehensively
- **ASK** explicit questions when seeking approval

## Pattern Steps

### 1. Requirements Clarification Phase
**Objective:** Establish clear understanding of business requirements through iterative clarification

**Activities:**
- [ ] Initial requirements gathering from stakeholder request
- [ ] Identify ambiguous or unclear requirements
- [ ] Ask clarifying questions about:
  - Business objectives and success criteria
  - User experience expectations
  - Functional requirements and edge cases
  - Non-functional requirements (performance, accessibility, etc.) if applicable
  - Integration points and dependencies
  - Constraints and limitations
- [ ] Iterate through multiple rounds of clarification as needed
- [ ] **Report back**: Summarize understood requirements and ask for confirmation
- [ ] **Seek approval**: Explicitly ask stakeholder to approve understanding before proceeding

**MANDATORY APPROVAL GATE**: You MUST receive explicit approval (e.g., "approved", "yes", "proceed") before continuing to Phase 2. Do not proceed without clear confirmation.

**Key Questions to Address:**
- What problem is this solving for users?
- What does success look like?
- Are there any constraints or limitations?
- How should this integrate with existing features?
- What are the expected user workflows?
- Are there any performance requirements?

**Outcomes:**
- Clear understanding of business requirements (summarized in response)
- Stakeholder confirmation of understanding
- Clear understanding of defined success criteria
- Clear understanding of identified constraints and dependencies

### 2. Technical Analysis Phase
**Objective:** Understand the scope and impact of the feature/request

**Activities:**
- [ ] Perform an audit of the existing code base
- [ ] Analyze relevant documentation (technical specs, design system, @docs/adr/)
- [ ] Assess impact on existing functionality
- [ ] Identify potential conflicts with established patterns
- [ ] Identify files that need modification
- [ ] **Report back**: Present technical analysis findings and implementation approach

**MANDATORY CHECKPOINT**: Present analysis in structured format before proceeding to Phase 3.

**Outcomes:**
- A clear understanding of impacted files and components
- A clear understanding of identified dependencies and integration points
- Assessment of potential risks or conflicts
- A clear understanding of alignment with existing architectural decisions

### 3. Checklist Creation
**Objective:** Create a comprehensive implementation checklist

**Activities:**
- [ ] Break down the feature into discrete tasks
- [ ] Identify subtasks for each major component
- [ ] Define acceptance criteria for each task
- [ ] Estimate complexity and dependencies
- [ ] Identify testing requirements
- [ ] Audit any existing test scripts
  - Plan to create any new scripts (as required)
  - Identify any scripts that need to be updated (as required)
  - Identify any scripts that need to be deleted (as required)
- [ ] Document any new architectural decisions needed
- [ ] **Report back**: Present comprehensive implementation plan and ask for feedback

**MANDATORY APPROVAL GATE**: You MUST receive explicit approval before proceeding to Phase 4.

**Outcomes:**
- Comprehensive task breakdown (can be summarized in response)
- Defined subtasks with acceptance criteria
- Testing strategy outline
- Risk assessment summary

### 4. Review and Clarification
**Objective:** Ensure alignment and gather feedback before implementation

**Activities:**
- [ ] **Report back**: Present checklist to stakeholders for review
- [ ] Clarify any ambiguous requirements
- [ ] Validate architectural alignment
- [ ] Finalize implementation approach
- [ ] **Seek approval**: Explicitly request approval to proceed with implementation

**MANDATORY APPROVAL GATE**: You MUST receive explicit approval (e.g., "approved", "proceed") before starting Phase 5.

**Outcomes:**
- Stakeholder approval of implementation approach
- Any additional clarified requirements
- Confirmed architectural approach

### 5. Implementation Plan Creation
**Objective:** Create detailed implementation documentation

**Activities:**
- [ ] Create implementation plan in ~/Documents/pomodux/ (per Documentation Guidelines in CLAUDE.md)
- [ ] Use standardized template if available
- [ ] Document all tasks and subtasks
- [ ] Define success criteria
- [ ] Set up tracking structure
- [ ] **Report back**: Present created implementation plan and confirm it's ready for execution

**MANDATORY CHECKPOINT**: Confirm plan is complete and ready before execution.

**Outcomes:**
- Implementation plan document created in ~/Documents/pomodux/
- Task tracking structure established
- Success criteria documented

### 6. Implementation Execution
**Objective:** Execute the implementation with ongoing documentation

**Activities:**
- [ ] Implement features according to the plan
- [ ] Update implementation plan with progress
- [ ] Mark tasks as complete with details
- [ ] Document any deviations or challenges
- [ ] Update status throughout implementation
- [ ] **Report back**: Provide regular progress updates during implementation
- [ ] **Final report**: Present completed implementation and request final acceptance

**MANDATORY APPROVAL GATE**: You MUST receive explicit acceptance before proceeding to Phase 7.

**Outcomes:**
- Working feature implementation
- Implementation plan updated with progress
- Issues and changes documented

### 7. Final Documentation Audit and Acceptance
**Objective:** Final audit and update of documentation with confirmation and acceptance

**Activities:**
- [ ] Audit all documentation in the @docs folder for accuracy and completeness
- [ ] Update documentation only if required (avoid unnecessary changes)
- [ ] Confirm documentation aligns with implemented features
- [ ] **Seek Approval**: Request acceptance of any updated documentation
- [ ] **Final report**: Present documentation status and request final acceptance

**MANDATORY APPROVAL GATE**: You MUST receive final approval before considering the feature complete.

**Outcomes:**
- Documentation accurately reflects implementation
- Only necessary documentation updates completed
- Documentation acceptance confirmed