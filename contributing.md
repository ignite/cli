# Contributing Guidelines

Before submitting a PR to the Ignite CLI repository, please review and follow these guidelines to ensure consistency and smooth collaboration across the project.

If you have suggestions or want to propose changes to these guidelines, start a new [Discussion topic](https://github.com/ignite/cli/discussions/new) to gather feedback.

For setup instructions, see [Set Up Your Ignite CLI Development Environment](dev-env-setup.md).

To contribute to docs and tutorials, see [Contributing to Ignite CLI Docs](https://docs.ignite.com/contributing).

We appreciate your contribution!

- [Contributing Guidelines](#contributing-guidelines)
  - [Providing Feedback](#providing-feedback)
  - [Opening Pull Requests (PRs)](#opening-pull-requests-prs)
    - [Choosing a Good PR Title](#choosing-a-good-pr-title)
    - [Reviewing Your Own Code](#reviewing-your-own-code)
    - [Commit Guidelines \& Title Conventions](#commit-guidelines--title-conventions)
    - [Do Not Rebase After Opening a PR](#do-not-rebase-after-opening-a-pr)
  - [Contributing to Documentation](#contributing-to-documentation)
    - [Ask for Help](#ask-for-help)
  - [Prioritizing Issues with Milestones](#prioritizing-issues-with-milestones)
  - [Issue Title Conventions and Labeling](#issue-title-conventions-and-labeling)

## Providing Feedback

- Before opening an issue, search for [existing open and closed issues](https://github.com/ignite/cli/issues) to check if your question has already been addressed. If a relevant issue exists, consider commenting on it instead of opening a duplicate issue.
- For feedback, questions, or suggestions, open a [Discussion topic](https://github.com/ignite/cli/discussions/new) to share your thoughts. Providing detailed information, such as use cases and links, will make the discussion more productive and actionable.

- For quick questions or informal feedback, join the **#🛠️ build-chains** channel in the official [Ignite Discord](https://discord.gg/ignite).

## Opening Pull Requests (PRs)

Please review relevant issues and discussions before opening a PR to ensure alignment with ongoing work.

### Choosing a Good PR Title

- Keep PR titles concise (fewer than 60 characters).
- Follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) guidelines for structuring your titles. For example: `feat(services/chain)`, `fix(scaffolding)`, `docs(migration)`.
- Your PR title should reflect the purpose of the changes and follow a consistent format.

### Reviewing Your Own Code

- Manually test your changes before submitting a PR or adding new commits.
- Ensure all CI checks pass before requesting a review. Your PR should show **All checks have passed** with a green checkmark.

### Commit Guidelines & Title Conventions

- **Standardized Issue Prefixes:**  
  Issue titles should begin with one of the following standardized prefixes, depending on the type of action being taken:

  - **`FIX:`** for resolving bugs or problems within existing features.
  - **`INIT:`** for creating new components, features, or initiatives.
  - **`UPDATE:`** for making improvements or modifications to existing functionality.
  - **`META:`** for larger, multi-step initiatives that consist of multiple tasks (e.g., epics).

  **Examples:**

  - `FIX: Resolve crash during chain initialization`
  - `INIT: Add staking module to example chain`
  - `UPDATE: Improve performance of block synchronization`
  - `META: Overhaul user permissions system`

- **Why Standardized Prefixes?**  
  The use of standardized prefixes ensures that the focus is on what needs to be done, making the task clear and actionable. This approach avoids redundancy with Conventional Commits, which are used for PR titles and commit messages to capture the purpose of the change. By separating the action (described by the prefix) from the nature of the issue (captured by labels), we reduce duplication and improve clarity. For example, if the issue is labeled `type:bug`, there’s no need to state "bug" in the title—the `FIX:` prefix already implies that the task involves resolving a bug.

- **Labels for Characteristics:**  
  Labels are used to classify the characteristics, elements, and descriptors of the issue or initiative. Labels help clarify the type of issue, the component involved, and its priority or status, without cluttering the title. Here are some examples:

  - **Type:** Describes the nature of the issue.

    - `type:bug` – Something isn't working.
    - `type:feat` – A new feature to be implemented.
    - `type:refactor` – Refactoring code without adding features.

  - **Component:** Specifies the part of the system the issue is related to.

    - `component:scaffold` – Related to scaffolding configuration or logic.
    - `component:frontend` – Related to frontend components.
    - `component:network` – Related to networking features or configurations.

  - **Status:** Indicates the current status of the issue or PR.
    - `status:needs-triage` – Needs to be reviewed and prioritized.
    - `status:blocked` – Cannot proceed until the blocking matter is resolved.
    - `status:help wanted` – Additional input or attention is needed.

### Do Not Rebase After Opening a PR

- Avoid rebasing commits once a PR is open for review. Instead, add additional commits as needed.
- Force pushes are acceptable only when the PR is in draft mode and hasn't been reviewed yet.

PRs will be squashed into a single commit when merged, so don't worry about having too many commits during the review process. The final PR title will be used as the commit message.

## Contributing to Documentation

Changes to the Ignite CLI codebase often require updates to the corresponding documentation. Please ensure that you update relevant documentation when making code changes.

- For changes to the [Developer Guide](https://docs.ignite.com/guide) and tutorials, update content in the `/docs/docs/02-guide` folder.
- For changes to the [Ignite CLI Reference](https://docs.ignite.com/references/cli), update the `./ignite/cmd` package where the command is defined. Do not edit auto-generated docs under `docs/docs/08-references/01-cli.md`.

### Ask for Help

If you start a PR but cannot complete it for any reason, don’t hesitate to ask for help. Another contributor can take over and finish the work.

## Prioritizing Issues with Milestones

We use Git Flow as our branch strategy, with each MAJOR release linked to a milestone. Core maintainers manage the prioritization of issues on the project board to ensure that the most critical work is addressed first.

- **Priority Labels (P0-P3):**  
  Issues are classified based on their urgency and impact, which helps guide the team’s focus during each release cycle:

  - **P0:** Urgent ("drop everything"); requires immediate attention and resolution. These issues take precedence over all other work.
  - **P1:** High priority ("important matter"); important and should be addressed promptly, though not as immediately critical as P0 issues.
  - **P2:** Medium priority ("sometime soon"); should be addressed but can be scheduled after P0 and P1 issues are resolved.
  - **P3:** Low priority ("nice to have"); nice to have but can be deferred or addressed as time permits.

- **Milestones and Workflow:**  
  Each milestone represents a MAJOR release. Issues are assigned to milestones based on their priority and relevance to the release goals. The project board is used to track and manage the progress of these issues. This structured workflow ensures that urgent tasks (P0) are addressed immediately, while lower-priority tasks (P3) are handled as resources allow.

- **Next Milestone:**  
  The **Next** milestone is used for issues or features that are not tied to a specific release but are still relevant to the project’s roadmap. These issues will be addressed when higher-priority work has been completed, or as part of future planning.

Check the [project board](https://github.com/ignite/cli/projects/7) to see what we're working on and what’s planned.

## Issue Title Conventions and Labeling

To maintain consistency across issues and PRs, follow these guidelines for issue titles:

- **Standardized Prefixes:** Begin with one of the standardized prefixes:

  - `FIX:` for resolving bugs.
  - `INIT:` for new components or projects.
  - `UPDATE:` for improving or modifying existing features.
  - `META:` for meta tasks involving multiple sub-tasks or actions.

- **Labels for Characteristics:** Use labels to classify the nature of the issue, such as its type, component, or status. Labels help describe the various elements of the issue or task, making it easier to manage and prioritize.

By combining standardized prefixes with well-organized labels, we maintain clarity, avoid redundancy, and ensure that issues and PRs are properly categorized and actionable.
