# Agent Guardrails & Guidelines

This document outlines the operational protocols for AI agents working on the current project.

## 1. Environment & Tooling
*   **Mise First:** Always prioritize `mise` for environment setup, tool versioning, and task execution.
    *   Check `mise.toml` for available tools and versions.
    *   Use `mise run <task>` for project commands defined in `mise.toml` (or `Makefile` if wrapped).
*   **Tool Documentation:** If you are unsure about a tool's usage or capabilities, **always** ask for or fetch its documentation before proceeding. Do not guess flags or parameters.

## 2. Planning & Implementation
*   **Implementation Plans:** Before writing code for a new feature, check if an implementation plan exists.
    *   If **NO**: Create a detailed `IMPLEMENTATION_PLAN.md` first. This file is used to plan specific features or phases before implementation. **Do not commit this file.**
    *   If **YES**: Follow the existing plan.
*   **Global Plan:** `PLAN.md` tracks the high-level project goals and architecture. Keep it updated if architectural decisions change.
*   **Plan Updates:** If a decision is made during the conversation that alters the agreed-upon approach, **immediately update** the relevant plan (`IMPLEMENTATION_PLAN.md` for details, `PLAN.md` for high-level) to reflect this change.

## 3. Git Workflow & Commits
When asked to commit changes, follow this strict two-step process:

### Step 1: Analysis & Planning
1.  **Analyze** the current `git status` and `git diff`.
2.  **Group** changes into logical, atomic units. Avoid monolithic "catch-all" commits.
3.  **Plan** the commit messages using the **Conventional Commits** standard (e.g., `feat:`, `fix:`, `chore:`, `docs:`).
4.  **Present** this plan to the user for approval.

### Step 2: Execution (After Approval)
1.  **Generate** the specific `git add` and `git commit` commands to execute the approved plan.
2.  Execute the commands only after explicit user confirmation.

## 4. General Behavior
*   Be concise and direct.
*   Make surgical changes; do not modify unrelated files.
*   Validate changes (run tests/linters) before committing.

