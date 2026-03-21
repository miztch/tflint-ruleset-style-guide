# tflint-ruleset-style-guide

[![Build Status](https://github.com/miztch/tflint-ruleset-style-guide/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/miztch/tflint-ruleset-style-guide/actions)

A TFLint ruleset based on the [Terraform Style Guide](https://developer.hashicorp.com/terraform/language/style).

## Installation

Add the following to your `.tflint.hcl`:

```hcl
plugin "style-guide" {
  enabled = true

  version = "0.1.0"
  source  = "github.com/miztch/tflint-ruleset-style-guide"

  signing_key = <<-KEY
  -----BEGIN PGP PUBLIC KEY BLOCK-----
  (YOUR PUBLIC KEY HERE)
  -----END PGP PUBLIC KEY BLOCK-----
  KEY
}
```

Then run:

```bash
tflint --init
```

## Rules

| Name | Description | Severity | Enabled |
| --- | --- | --- | --- |
| [style_guide_typed_variables_except_any](docs/rules/style_guide_typed_variables_except_any.md) | Disallow `any` as variable type | WARNING | ✔ |
| [style_guide_type_repetition](docs/rules/style_guide_type_repetition.md) | Disallow repeating the resource type in the resource name | WARNING | ✔ |
| [style_guide_meta_arguments_blank_line](docs/rules/style_guide_meta_arguments_blank_line.md) | Require blank lines around meta-arguments | WARNING | ✔ |

---

## Building the plugin

```bash
make
```

Install the built plugin locally:

```bash
make install
```