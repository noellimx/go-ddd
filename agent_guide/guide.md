# Agent

Using `copilot`.

## Prerequisites

- An active GitHub Copilot subscription
- `brew install copilot-cli`

## Run

- check bin in path `which copilot`
- run cli `copilot` 
  - check authenticated `/user list`

### Prompts
!! only for vscode
Example file: `.github/prompts/explain-code.prompt.md`

### Custom Instructions

1. Repo-wide instructions — create `.github/copilot-instructions.md`
2. Scoped instructions — example: `.github/instructions/comment.instructions.md`

### Agents
roles copilot can take on. customized agent profiles are in `.github/agents`

Example: `.github/agents/red.agent.md`

### References

https://docs.github.com/en/copilot/how-tos/copilot-cli/use-copilot-cli-agents/overview
https://docs.github.com/en/copilot/how-tos/copilot-cli/set-up-copilot-cli/install-copilot-cli
https://docs.github.com/en/copilot/tutorials/customization-library/prompt-files/your-first-prompt-file
https://docs.github.com/en/copilot/how-tos/copilot-cli/customize-copilot/add-custom-instructions
https://docs.github.com/en/copilot/how-tos/use-copilot-agents/coding-agent/create-custom-agents

### Samples
https://github.com/github/awesome-copilot/tree/main/agents