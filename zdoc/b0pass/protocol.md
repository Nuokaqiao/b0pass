# Project Knowledge Protocol

Version: 1.1

---

# 1. Purpose

This Project Knowledge Base is a local, tool-agnostic knowledge system
for understanding and maintaining a software project.

It is designed to be used by:

- Humans
- Claude Code
- Codex
- Cursor
- Gemini
- Other AI coding agents
- Future AI tools and agent harnesses

The Knowledge Base must not depend on any specific AI tool.

The Knowledge Base is a long-term project memory.

Its purpose is to preserve:

1. Understanding of the current system
2. Business knowledge
3. Architecture knowledge
4. Important workflows
5. Important rules and constraints
6. Design decisions and their reasoning
7. Meaningful system changes
8. Known risks and reliability concerns

The Knowledge Base exists outside the Git repository.

It should normally NOT be committed together with the source code.

---

# 2. Core Principles

## 2.1 Source Code Is the Primary Source of Truth

The current source code, configuration, database schema,
tests and authoritative project documentation are the primary
sources of truth.

The Knowledge Base is a derived representation of the project.

When Knowledge conflicts with the current implementation:

1. Trust the current source code.
2. Identify the outdated Knowledge.
3. Update the Knowledge Base.
4. Explain the discrepancy when relevant.

Never blindly trust old Knowledge.

---

## 2.2 Knowledge Is Long-Lived

Knowledge should describe durable information that helps
humans or AI understand the project.

Examples:

- Architecture
- Business domains
- Business rules
- State machines
- Important data relationships
- Important dependencies
- Reliability mechanisms
- Important constraints

Do not use Knowledge as a copy of the source code.

---

## 2.3 Knowledge Is Incremental

Do not repeatedly rewrite the entire Knowledge Base.

When new information is discovered:

1. Search existing Knowledge.
2. Identify the appropriate document.
3. Update only the affected sections.
4. Preserve useful existing explanations.
5. Create a new document only when necessary.

Avoid duplicate sources of truth.

---

# 3. Knowledge Structure

The Knowledge Base normally contains:

    index.md

    protocol.md

    knowledge/
        overview.md
        architecture.md
        modules.md
        data-model.md
        business-flows.md
        business-rules.md
        state-machines.md
        integrations.md
        jobs-and-workers.md
        reliability.md
        risks.md

    modules/
        <module>.md

    decisions/
        <decision>.md

    changes/
        <change>.md

    notes/
        <note>.md

---

# 4. Knowledge Categories

## 4.1 Knowledge

Knowledge describes the current system.

Typical topics:

- Project purpose
- Business domains
- Architecture
- Modules
- Data models
- Business flows
- Business rules
- State machines
- Integrations
- Jobs and workers
- Reliability
- Risks

Knowledge should describe current behavior.

---

## 4.2 Modules

Module documents contain deeper knowledge about
important business or technical modules.

A module document may include:

- Purpose
- Responsibilities
- Entry points
- Important APIs
- Important classes/functions
- Data
- Dependencies
- Business flows
- State transitions
- Business rules
- Error handling
- Reliability
- Risks
- Important source paths

Create module documents only for modules that provide
meaningful independent value.

---

## 4.3 Decisions

Decision documents record important design decisions.

A Decision should answer:

- What problem were we solving?
- What was the context?
- What options were considered?
- Which option was selected?
- Why was it selected?
- What are the consequences?
- What alternatives were rejected?

Do not record trivial implementation details.

---

## 4.4 Changes

Change documents record meaningful changes to the system.

A Change should explain:

- What changed?
- Why did it change?
- What was the previous behavior?
- What is the new behavior?
- What modules are affected?
- What is the business impact?
- What is the technical impact?
- What important risks exist?

Do NOT copy Git commit history.

Only record changes that improve long-term understanding.

---

## 4.5 Notes

Notes contain temporary or exploratory information.

Examples:

- Unconfirmed ideas
- Investigation results
- Open questions
- Temporary debugging observations
- Design alternatives under discussion

Notes are not authoritative.

When a Note becomes verified and durable:

    notes/
        ↓
    Knowledge / Decision / Change

---

# 5. Knowledge Confidence

Important knowledge should indicate its confidence.

Use:

- Confirmed
- Inferred
- Unknown

## Confirmed

Directly verified from:

- Source code
- Configuration
- Tests
- Database schema
- API definitions
- Authoritative documentation

## Inferred

Reasonably inferred from available evidence,
but not directly confirmed.

## Unknown

The information is not currently known.

Never present Inferred or Unknown information
as confirmed fact.

---

# 6. Knowledge Freshness

Confidence and Freshness are different concepts.

Confidence answers:

    "How certain are we that this information is correct?"

Freshness answers:

    "Has this information been verified against the current code?"

Use:

- Fresh
- Potentially Stale
- Stale
- Unknown

Example:

    Status: Confirmed
    Freshness: Stale

means:

The information was previously verified,
but the source code has changed since verification.

---

# 7. Source Version

Important Knowledge documents should record the source-code
version against which they were last verified.

Example:

    Status: Confirmed
    Freshness: Fresh

    Last Verified: 2026-08-27
    Verified Commit: a81f92c

The exact format of the source version may depend on the
version control system.

For Git projects, normally use the commit SHA.

---

# 8. Related Source Paths

Important Knowledge documents may declare the source-code
paths that they describe.

Example:

    Related Paths:

    - internal/promotion/**
    - api/promotion/**
    - tests/promotion/**

Related Paths help identify whether a code change
may affect a Knowledge document.

They are an optimization, not an absolute rule.

A change outside the listed paths may still affect
the Knowledge if there are indirect dependencies.

---

# 9. Knowledge Drift

Knowledge Drift occurs when the source code changes after
Knowledge was last verified.

This is expected.

Knowledge may become outdated because:

- Another developer changed the code.
- A feature was added without updating Knowledge.
- A bug fix changed system behavior.
- A refactoring changed architecture.
- An integration changed.
- A database schema changed.
- A long period passed since Knowledge was verified.

The system must assume that Knowledge can become stale.

---

# 10. Drift Detection

Before answering a complex project question or starting
a significant development task:

1. Locate the Project Knowledge Base.
2. Read `index.md`.
3. Identify relevant Knowledge.
4. Check Knowledge freshness when practical.
5. Compare relevant Knowledge verification versions
   with the current source-code version.
6. Inspect relevant code changes when necessary.
7. Determine whether the Knowledge is still valid.

Do NOT assume that every Git commit invalidates all Knowledge.

Only potentially affected Knowledge should be revalidated.

---

# 11. Knowledge Health

The Project Knowledge Index may summarize the health
of important areas.

Example:

    | Area | Status | Verified Commit |
    |---|---|---|
    | Architecture | Fresh | a81f92c |
    | Promotion | Stale | a81f92c |
    | Translation | Fresh | a81f92c |
    | Payment | Potentially Stale | a81f92c |

Interpretation:

## Fresh

Relevant source code has not changed in a way that
appears to affect this Knowledge.

## Potentially Stale

Relevant source code has changed,
but the impact has not yet been fully verified.

## Stale

The Knowledge has been confirmed to be inconsistent
with the current implementation or is known to be outdated.

## Unknown

Freshness cannot currently be determined.

---

# 12. Drift Detection Workflow

When relevant code changed after the last verification:

    Knowledge
        ↓
    Check Verified Commit
        ↓
    Inspect changes
        ↓
    Determine affected areas
        ↓
    Compare code with Knowledge
        ↓
    ┌──────────────────────┐
    │                      │
    ↓                      ↓
  Still valid          Outdated
    │                      │
    ↓                      ↓
 Mark Fresh            Update Knowledge
                           ↓
                     Update verification
                           ↓
                     Update Index

Do not rewrite unrelated Knowledge.

---

# 13. Knowledge Synchronization

Knowledge Synchronization means bringing the Knowledge Base
up to date with the current implementation.

When synchronizing:

1. Identify the last verified source version.
2. Identify the current source version.
3. Inspect relevant changes.
4. Determine affected Knowledge.
5. Compare current implementation with Knowledge.
6. Update only affected documents.
7. Update `Last Verified`.
8. Update `Verified Commit`.
9. Update `Freshness`.
10. Update `index.md`.
11. Create a Change record when appropriate.
12. Create or update a Decision when an important design
    decision was introduced.

Synchronization should be incremental.

---

# 14. Important Rule About Other Developers

Developers are NOT required to update the local Knowledge Base
every time they modify source code.

The Knowledge Base is local and may belong to an individual developer.

Other developers may:

- Add features
- Modify features
- Refactor code
- Fix bugs
- Change architecture

without knowing about this Knowledge Base.

This is expected.

The Knowledge system must therefore detect Knowledge Drift
when the user returns to the project.

Do not assume:

    "Nobody updated Knowledge, therefore the code did not change."

Instead assume:

    "Knowledge may be behind the source code."

---

# 15. Reading Strategy

Do not read the entire Knowledge Base for every question.

Use targeted reading.

Preferred strategy:

    User Question
        ↓
    index.md
        ↓
    Relevant Knowledge
        ↓
    Relevant Module
        ↓
    Source Code Verification
        ↓
    Answer

For simple questions, existing Knowledge may be sufficient.

For important, ambiguous or potentially stale questions,
verify against source code.

---

# 16. When Source Code Must Be Checked

Source code should normally be checked when:

- Knowledge is Stale
- Knowledge is Potentially Stale
- The question concerns current implementation details
- The question concerns exact behavior
- A new feature is being designed
- A bug is being investigated
- A significant architectural change is involved
- The user asks whether the current code actually does something
- There is a conflict between Knowledge and source code

---

# 17. New Feature Workflow

For non-trivial features, follow:

    Understand
        ↓
    Impact Analysis
        ↓
    Design Discussion
        ↓
    Decision
        ↓
    Implementation
        ↓
    Verification
        ↓
    Knowledge Sync
        ↓
    Change Record

---

## Phase 1 — Understand

Read:

- Relevant Knowledge
- Relevant Module documents
- Existing source code
- Tests
- Related Decisions

Understand the current system before proposing changes.

---

## Phase 2 — Impact Analysis

Identify potentially affected:

- Modules
- APIs
- Data models
- Business flows
- State machines
- Dependencies
- Permissions
- Jobs
- Events
- Tests
- Reliability mechanisms
- Operational behavior

---

## Phase 3 — Design Discussion

For significant changes:

1. Explain current behavior.
2. Identify constraints.
3. Identify affected areas.
4. Propose possible approaches.
5. Compare trade-offs.
6. Recommend an approach.
7. Discuss the design with the user.

Do not silently make important architectural
or business decisions.

---

## Phase 4 — Decision

When an important design decision is accepted:

Create or update a Decision document.

The Decision should capture the reasoning,
not just the final implementation.

---

## Phase 5 — Implementation

Implement according to the accepted design.

Do not modify unrelated functionality without justification.

---

## Phase 6 — Verification

Verify the implementation as appropriate:

- Unit tests
- Integration tests
- Data behavior
- Error handling
- Concurrency
- Idempotency
- Transactions
- Backward compatibility
- Performance
- Operational impact

---

## Phase 7 — Knowledge Sync

After implementation:

1. Compare implementation with existing Knowledge.
2. Update affected Knowledge.
3. Update affected Module documents.
4. Update business flows.
5. Update state machines.
6. Update business rules.
7. Update reliability or risks if necessary.
8. Update verification metadata.
9. Record a meaningful Change.

---

# 18. Bug Fix Workflow

When investigating or fixing a bug:

1. Understand expected behavior.
2. Understand actual behavior.
3. Identify root cause.
4. Inspect relevant Knowledge.
5. Determine whether Knowledge was incomplete,
   incorrect or stale.
6. Fix the problem.
7. Verify the fix.
8. Update Knowledge if the bug reveals durable
   system behavior or an important rule.
9. Record a Change when the fix is significant.

Do not create a Change record for every trivial bug fix.

---

# 19. Architecture Changes

Architecture changes require special attention.

Examples:

- New service
- Service removal
- Database migration
- New message queue
- New external dependency
- New caching strategy
- New persistence strategy
- Major API redesign
- Major state-machine change

When these occur, review at least:

- architecture.md
- modules.md
- data-model.md
- integrations.md
- reliability.md
- risks.md

and any affected Module documents.

Important architecture decisions should normally
have a Decision record.

---

# 20. Business Rule Changes

When a change modifies business behavior,
review:

- business-rules.md
- business-flows.md
- state-machines.md
- affected Module documents

Examples:

- New validation rule
- Changed eligibility rule
- Changed status transition
- Changed permission behavior
- Changed calculation logic
- Changed retry behavior
- Changed release conditions

Business rules should not exist only in code
if they are important for understanding the system.

---

# 21. State Machine Changes

Whenever a state machine changes,
explicitly verify:

1. States
2. Allowed transitions
3. Entry conditions
4. Exit conditions
5. Invalid transitions
6. Side effects
7. Retry behavior
8. Failure behavior
9. Persistence
10. Related APIs/events/jobs

Update:

    knowledge/state-machines.md

and affected Module / Business Flow documents.

---

# 22. Reliability Knowledge

For systems involving money, orders, payments,
distributed processing or asynchronous processing,
pay special attention to:

- Idempotency
- Concurrency
- Transactions
- Retry
- Duplicate messages
- Message ordering
- Eventual consistency
- Data consistency
- Failure recovery
- Timeout behavior
- Race conditions
- Partial failure
- Exactly-once assumptions

Do not assume that an operation is idempotent
without verifying the implementation.

---

# 23. Evidence

Important technical claims should have useful evidence.

Examples:

    Source:
    - internal/promotion/service.go
    - PromotionService.Create()

or:

    Database:
    - promotion
    - promotion_history

or:

    Test:
    - internal/promotion/service_test.go

Evidence does not need to be attached to every trivial statement.

The purpose is to make important conclusions verifiable.

---

# 24. Avoiding Duplicate Knowledge

Before creating a new document:

1. Search existing Knowledge.
2. Determine whether the topic already exists.
3. Update the existing document when appropriate.
4. Create a new document only when it represents
   a distinct topic.

Avoid having multiple documents that describe
the same fact differently.

Prefer one canonical source.

---

# 25. Handling Conflicting Knowledge

If two Knowledge documents conflict:

1. Identify the conflicting claims.
2. Check the current source code.
3. Determine the correct behavior.
4. Update the outdated document.
5. Consolidate duplicated information.
6. Record a Decision if the conflict reveals
   an important architectural or business decision.

Never preserve contradictory Knowledge merely
for historical reasons.

Historical reasoning belongs in Decision or Change records.

---

# 26. Historical Information

Knowledge should normally describe the current system.

Historical information belongs in:

    decisions/
    changes/

Do not leave outdated historical behavior in current
Knowledge unless it is necessary to understand the current system.

Example:

Bad:

    Promotion currently supports draft and published.
    Previously it supported only draft and published,
    but before that it...

Better:

    Promotion currently supports:
    draft → pending_review → approved → published

Historical changes belong in:

    changes/2026-08-27-promotion-review.md

---

# 27. Human Review

AI may:

- Analyze
- Propose
- Compare
- Implement
- Synchronize Knowledge

However, important business and architectural decisions
should be explicitly confirmed by the human.

AI must not silently establish major business rules
or architectural decisions.

---

# 28. Uncertainty

Never invent project behavior.

When information cannot be verified:

    Status: Unknown

When information is based on reasonable inference:

    Status: Inferred

When information has been verified:

    Status: Confirmed

When appropriate, explain what evidence would be needed
to confirm an Inferred or Unknown statement.

---

# 29. Minimal Modification Principle

When updating Knowledge:

- Make the smallest useful change.
- Preserve existing useful explanations.
- Avoid unnecessary rewrites.
- Avoid changing unrelated documents.
- Avoid generating large amounts of repetitive text.

The Knowledge Base should become clearer over time,
not larger for the sake of being larger.

---

# 30. No Automatic Full Re-Summarization

Do not periodically regenerate the entire Knowledge Base
from scratch.

Prefer:

    Existing Knowledge
        +
    Code Changes
        ↓
    Targeted Update

Full re-analysis may be performed when:

- The Knowledge Base is severely outdated.
- Architecture has changed substantially.
- The user explicitly requests a full re-analysis.
- The current Knowledge is unreliable.

---

# 31. Knowledge Maintenance Trigger

Knowledge maintenance should be considered when:

- A new feature is completed.
- A significant bug is fixed.
- Architecture changes.
- Business rules change.
- State machines change.
- Data models change.
- External integrations change.
- Reliability mechanisms change.
- The user explicitly requests synchronization.
- Drift is detected.

---

# 32. Standard Knowledge Metadata

Important Knowledge documents may use:

    Status: Confirmed
    Freshness: Fresh

    Last Verified: 2026-08-27
    Verified Commit: a81f92c

    Related Paths:
    - internal/promotion/**
    - api/promotion/**
    - tests/promotion/**

Metadata may be omitted for very small documents,
but important Knowledge should use it.

---

# 33. Standard Decision Metadata

Example:

    Decision ID: DEC-0001
    Status: Accepted
    Date: 2026-08-27

    Related Knowledge:
    - knowledge/state-machines.md
    - modules/promotion.md

---

# 34. Standard Change Metadata

Example:

    Change ID: CHG-2026-001
    Status: Implemented
    Date: 2026-08-27

    Related Decision:
    - decisions/0001-promotion-state-machine.md

    Related Knowledge:
    - knowledge/state-machines.md
    - knowledge/business-flows.md

---

# 35. Standard Workflow for Any Complex Question

When receiving a complex project question:

    1. Locate Knowledge Base
    2. Read protocol.md
    3. Read index.md
    4. Identify relevant Knowledge
    5. Check freshness
    6. Inspect source code if necessary
    7. Answer
    8. Determine whether durable new knowledge was discovered
    9. Update Knowledge if necessary

Do not update Knowledge merely because
a question was asked.

Update it when durable project knowledge
has been discovered or corrected.

---

# 36. Standard Workflow for Project Entry

When an AI agent first enters the project:

    1. Find the Project Knowledge Base.
    2. Read protocol.md.
    3. Read index.md.
    4. Check the Knowledge health.
    5. Identify relevant Knowledge for the current task.
    6. Verify important information against source code.
    7. Continue with the user's request.

Do not immediately scan the entire repository
unless the task requires it.

---

# 37. Standard Workflow for Knowledge Sync

When the user asks:

    "Sync Project Knowledge"

perform:

    1. Determine current source version.
    2. Determine Knowledge verification versions.
    3. Identify changed source areas.
    4. Map changes to affected Knowledge.
    5. Inspect affected implementation.
    6. Identify outdated or missing Knowledge.
    7. Update affected Knowledge.
    8. Update Module documents.
    9. Update Business Flows / Rules / State Machines.
    10. Update Freshness metadata.
    11. Update Verified Commit.
    12. Update index.md.
    13. Create meaningful Change records.
    14. Create Decision records for important decisions.
    15. Report remaining Knowledge gaps.

Do not modify unrelated Knowledge.

---

# 38. Reporting Knowledge Drift

When significant drift is detected,
report it clearly.

Example:

    Knowledge Drift Detected

    Affected areas:

    - Promotion
    - Release Pipeline

    Changes since last verification:

    - Added pending_review state
    - Added approval API
    - Changed publishing conditions

    Knowledge status:

    - modules/promotion.md → Stale
    - state-machines.md → Stale
    - architecture.md → Fresh

    Recommended action:

    Synchronize affected Knowledge.

Do not silently hide significant discrepancies.

---

# 39. Project Knowledge Is Not a Task Log

Do not turn the Knowledge Base into:

- Chat transcripts
- Daily work logs
- Git commit lists
- TODO lists
- Debugging history
- Personal diary

The Knowledge Base exists to preserve
useful long-term project understanding.

---

# 40. Final Principles

The Project Knowledge Base should provide:

    Human Understanding
          +
    AI Understanding
          +
    Long-Term Continuity

The desired workflow is:

    Source Code
        ↓
    Analysis
        ↓
    Knowledge
        ↓
    Human + AI Understanding
        ↓
    New Development
        ↓
    Source Code Changes
        ↓
    Drift Detection
        ↓
    Knowledge Synchronization
        ↓
    Updated Knowledge

The AI agent is replaceable.

The source code evolves.

The Knowledge Base is the long-term project memory.

The Knowledge Base must remain:

- Tool-agnostic
- Local
- Incremental
- Verifiable
- Human-readable
- AI-readable
- Resistant to knowledge drift
- Focused on durable understanding