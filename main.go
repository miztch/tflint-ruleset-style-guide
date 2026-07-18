package main

import (
	"github.com/miztch/tflint-ruleset-style-guide/project"
	"github.com/miztch/tflint-ruleset-style-guide/rules"

	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "style-guide",
			Version: project.Version,
			Rules: []tflint.Rule{
				rules.NewStyleGuideMetaArgumentsBlankLineRule(),
				rules.NewStyleGuideOrderedOutputArgumentsRule(),
				rules.NewStyleGuideOrderedResourceArgumentsRule(),
				rules.NewStyleGuideOrderedVariableArgumentsRule(),
				rules.NewStyleGuideTypeRepetitionRule(),
				rules.NewStyleGuideTypeVariablesExceptAnyRule(),
			},
		},
	})
}
