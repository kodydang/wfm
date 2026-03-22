## ADDED Requirements

### Requirement: alias command exists
The system SHALL provide an `alias` subcommand under `kd-wfm` that prints shell alias definitions to stdout.

#### Scenario: Command appears in help
- **WHEN** user runs `kd-wfm --help`
- **THEN** `alias` SHALL appear in the list of available subcommands with a short description

### Requirement: Alias output format
The system SHALL print one alias definition per line in a format compatible with both bash and zsh.

#### Scenario: Output format
- **WHEN** user runs `kd-wfm alias`
- **THEN** each line SHALL follow the format `alias <name>='kd-wfm <name>'`

#### Scenario: gitcm alias included
- **WHEN** user runs `kd-wfm alias`
- **THEN** the output SHALL include `alias gitcm='kd-wfm gitcm'`

### Requirement: Alias list derived from registered commands
The system SHALL dynamically generate alias output from Cobra's registered subcommands, excluding internal commands.

#### Scenario: Internal commands excluded
- **WHEN** user runs `kd-wfm alias`
- **THEN** `help`, `completion`, and `alias` itself SHALL NOT appear in the output

#### Scenario: New commands automatically included
- **WHEN** a new subcommand is registered in root.go
- **THEN** it SHALL automatically appear in `kd-wfm alias` output without code changes to the alias command

### Requirement: Dotfiles portability
The system SHALL support common dotfiles integration patterns.

#### Scenario: Pipe to shell config
- **WHEN** user runs `kd-wfm alias >> ~/.zshrc`
- **THEN** the alias lines SHALL be appended and usable after the next shell reload

#### Scenario: Eval pattern
- **WHEN** user adds `eval "$(kd-wfm alias)"` to their shell config
- **THEN** all aliases SHALL be active in new shell sessions
