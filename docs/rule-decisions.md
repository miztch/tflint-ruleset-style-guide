# Rule Decisions

This document tracks Terraform [Style Guide](https://developer.hashicorp.com/terraform/language/style)
recommendations that are **not** covered by
[tflint-ruleset-terraform](https://github.com/terraform-linters/tflint-ruleset-terraform),
and records the decisions made about them so the same discussions and research
are not repeated later.

This document covers only candidates decided **not** to implement.
Implemented rules are listed in the [README](../README.md), and candidates
decided to implement are tracked as GitHub issues.

## Decided not to implement

| Candidate | Reason |
| --- | --- |
| Local module source path (`./modules/<name>`) | Common environment-split monorepos reference modules as `../../modules/<name>`, which the guide's literal path doesn't match; either the rule is noisy or the check drifts from the guide's wording |
| Ordered provider blocks | "Default provider first" needs cross-file ordering, which is arbitrary in HCL; "alias as first argument" alone is rarely violated and low-value |

## Out of scope

Recommendations that are impractical for TFLint's HCL static analysis:

- `.gitignore` hygiene (state files, `.terraform`, plan files)
- Module repository naming (`terraform-<PROVIDER>-<NAME>`)
- Presence of `README` / `LICENSE` / `examples/` (Standard Module Structure)
- Branching strategy, environment separation, workspace layout, state sharing, secrets management, testing, policy
- "Avoid overuse of variables/locals", "use `count`/`for_each` sparingly" — no objective threshold
