# ads-toyrepo

A small Go module, used to try `ads` end to end. Its packages are `cache`, `ledger`, `parse`, and `report`.

It is a copy of `bench/toyrepo` from `sondalab-ai/agentic-dev-system`. The `ads` benchmark uses the same code.

## Try ads here

1. `git clone https://github.com/sondalab-ai/ads-toyrepo.git && cd ads-toyrepo`
2. `ads init`: the onboarding wizard. It sets up the signing key, ssh-agent, the hooks, and the GitHub rulesets.
3. `ads doctor`: checks the setup and changes nothing.

`go test ./...` must pass.
