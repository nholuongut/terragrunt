engine {
  source  = "github.com/gruntwork-io/terragrunt-engine-opentofu"
  version = "v0.0.5"
  type    = "rpc"
}

inputs = {
  value = "test_input_value_from_terragrunt"
}
