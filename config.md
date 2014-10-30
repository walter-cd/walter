---
layout: default
title: Configuration
---

This page describes the settings of Walter.

# Configuration setting

Walter has one configuration file, which specifies a set of tasks needed to build or deploy target applicaion or service. 
More specifically, users specify the order of the tasks for the deployment. Each task in the flow is called as **stage**,
and the flow is called **pipeline** in Walter.

The configuration format of Walter is Yaml. The yaml configuration file need to have one pipeline block, which has more
than one stage element.

The following is a sample configuration of Walter.

     pipeline:
          command_stage_1:
             stage_type: command
             command: echo "hello, world"
         command_stage_2:
             stage_type: command
             command: echo "hello, world, command_stage_2"
         command_stage_3:
             stage_type: command
             command: echo "hello, world, command_stage_3"

As we see, the pipeline block has three stages and the stage type is command, each of which run **echo** command and has the stage name
(such as **command_stage_1**). User can name arbitary name of each stage. The commands are excuted with the same order as the pipeline configuration.

## Stage setting

Stage in pipeline has three elements, **name** **stage_type** and the **configurations** if needed. The configurations depend on the stage_type.
For example command_stage type has **command** configuration, which specify the shell command run in the stage.
The following is the table on the stage_type and the parameters.

### Command stage

Command stage executes one command. Users specify Command stage adding **command** in stage_type.

The follwoing is the parameter of Command stage.

|  Configuration | Optional   | meaning                                |
|:--------------:|:----------:|:--------------------------------------:|
|   command      | false      | shell command run in the stage         |

### Shell script stage
Command stage executes specified shell scrpit file. Users specify Command stage adding **shell_script** in stage_type.

The follwoing is the parameter of Command stage.

|  Configuration   | Optional   | meaning                                |
|:----------------:|:----------:|:--------------------------------------:|
|   shell _script  | false      | shell script file run in the stage     |

