## ADDED Requirements

### Requirement: gitcm command exists
The system SHALL provide a `gitcm` subcommand under `kd-wfm` that generates and applies a conventional git commit message using the `claude` CLI.

#### Scenario: Command appears in help
- **WHEN** user runs `kd-wfm --help`
- **THEN** `gitcm` SHALL appear in the list of available subcommands with a short description

### Requirement: Staged diff check
The system SHALL abort early when there are no staged changes.

#### Scenario: No staged changes
- **WHEN** user runs `kd-wfm gitcm` and `git diff --cached` returns empty output
- **THEN** the command SHALL print "no staged changes to commit" and exit with a non-zero status

#### Scenario: Staged changes present
- **WHEN** user runs `kd-wfm gitcm` and there are staged changes
- **THEN** the command SHALL proceed to generate a commit message

### Requirement: LLM commit message generation
The system SHALL invoke the `claude` CLI subprocess with an embedded system prompt and the staged diff to generate a commit message.

#### Scenario: claude CLI invoked with diff
- **WHEN** staged changes are present
- **THEN** the system SHALL run `claude` with the system prompt from `prompts/gitcm-system.md` (embedded) and pipe the staged diff as input

#### Scenario: claude CLI not found
- **WHEN** the `claude` binary is not found in PATH
- **THEN** the command SHALL print "claude CLI not found — install Claude Code to use this command" and exit with a non-zero status

#### Scenario: claude CLI returns output
- **WHEN** `claude` exits successfully
- **THEN** the system SHALL use its stdout as the proposed commit message

### Requirement: Commit confirmation
The system SHALL show the generated commit message and require explicit user confirmation before committing.

#### Scenario: Confirmation prompt shown
- **WHEN** a commit message has been generated
- **THEN** the system SHALL print the message and prompt "Commit with this message? [y/N]"

#### Scenario: User confirms
- **WHEN** user enters `y` or `Y` at the confirmation prompt
- **THEN** the system SHALL run `git commit -m "<message>"` and print the git output

#### Scenario: User declines
- **WHEN** user enters anything other than `y`/`Y` at the confirmation prompt
- **THEN** the system SHALL print "Aborted. Generated message:" followed by the message and exit with status 0

### Requirement: System prompt embedded in binary
The system SHALL embed the LLM system prompt file at compile time so the binary is self-contained.

#### Scenario: Prompt file embedded
- **WHEN** the binary is compiled
- **THEN** the content of `prompts/gitcm-system.md` SHALL be embedded via `//go:embed` and available at runtime without filesystem access
