### Requirement: init command exists
The system SHALL provide an `init` subcommand under `kd-wfm` that installs the binary and registers a shell alias.

#### Scenario: Command appears in help
- **WHEN** user runs `kd-wfm --help`
- **THEN** `init` SHALL appear in the list of available subcommands with a short description

### Requirement: Binary installation
The system SHALL copy the running binary to `~/.kd-wfm/kd-wfm`, creating the directory if it does not exist.

#### Scenario: Directory created if missing
- **WHEN** user runs `kd-wfm init` and `~/.kd-wfm/` does not exist
- **THEN** the command SHALL create `~/.kd-wfm/` with permissions `0755`

#### Scenario: Binary copied
- **WHEN** user runs `kd-wfm init`
- **THEN** the running binary SHALL be copied to `~/.kd-wfm/kd-wfm` with permissions `0755`

#### Scenario: Source equals destination
- **WHEN** the running binary is already at `~/.kd-wfm/kd-wfm`
- **THEN** the command SHALL print "binary already installed at ~/.kd-wfm/kd-wfm, skipping" and skip the copy

### Requirement: Shell alias registration
The system SHALL append `alias kd-wfm='~/.kd-wfm/kd-wfm'` to `~/.zshrc` if not already present.

#### Scenario: Alias appended when absent
- **WHEN** user runs `kd-wfm init` and the alias line is not in `~/.zshrc`
- **THEN** the command SHALL append `alias kd-wfm='~/.kd-wfm/kd-wfm'` to `~/.zshrc`

#### Scenario: Alias skipped when present
- **WHEN** user runs `kd-wfm init` and the alias line already exists in `~/.zshrc`
- **THEN** the command SHALL print "alias already in ~/.zshrc, skipping" and not append

#### Scenario: zshrc created if missing
- **WHEN** `~/.zshrc` does not exist
- **THEN** the command SHALL create it and write the alias line

### Requirement: Idempotency
The system SHALL be safe to run multiple times without side effects.

#### Scenario: Re-run produces no duplicates
- **WHEN** user runs `kd-wfm init` a second time after a successful first run
- **THEN** both the binary copy and alias append SHALL be skipped with informational messages
