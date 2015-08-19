require = ["tests/fixtures/s1_stages.hcl", "tests/fixtures/s2_stages.hcl"]

pipeline {
  stage {
    name = "command_stage_1"
    type = "command"
    command = "echo \"hello, world\""
  }
  stage {
    call = "s1::foo"
  }
  stage {
    call = "s2::foo"
  }
  stage {
    name = "parallel run"
    parallel {
      stage {
        call = "s2::bar"
      }
      stage {
        call = "s2::baz"
      }
    }
  }
}
