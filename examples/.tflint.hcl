plugin "terraform" {
  enabled = false

  version = "0.14.1"
  source  = "github.com/terraform-linters/tflint-ruleset-terraform"
}

plugin "style-guide" {
  enabled = true
}
