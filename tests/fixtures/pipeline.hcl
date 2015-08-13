/*
  An example of the pipeline yaml test file written in HCL
*/

//Global parameters
global {
  dummyOpt = true
}

//The pipeline
pipeline {
  //Why not give it a name?
  name = "my-hcl-pipeline"

  stage {
    name = "command_stage_1"
    type = "command"
    command = "echo 'hello, world'"

    //stages can contain parallel stages
    parallel {
      stage {
        name = "command_stage_2_group_1"
        type = "command"
        command = "echo 'hello, world, command_stage_2_group_1'"
      }
      stage {
        name = "command_stage_3_group_1"
        type = "command"
        command = "echo 'hello, world, command_stage_3_group_1'"

        //let's go deeper!
        parallel {
          stage {
            name = "command_stage_3_group_2"
            type = "command"
            command = "echo 'hello, world, command_stage_3_group_2'"
          }
        }
      }
    }
  }
  stage {
    name = "command_stage_4"
    type = "command"
    command = "echo 'hello, world, stage4' && sleep 1"
  }
  stage {
    name = "command_stage_5"
    type = "command"
    command = "echo 'hello, world, stage5'"
  }
}

cleanup {
  stage {
    name = "cleanup_stage_1"
    type = "command"
    command = "echo 'hello cleanup'"
  }
}
