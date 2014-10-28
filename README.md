Walter
========

[![wercker status](https://app.wercker.com/status/9b663fe6a4c4eace9a0a3be1fe71757e/m/master "wercker status")](https://app.wercker.com/project/bykey/9b663fe6a4c4eace9a0a3be1fe71757e)

Walter is a tiny deployment pipeline template.

Overview
==========

Walter automates the deployment process of appllications or servers from configuration to sofware deployment.

Getting Started
===============

Requirements
-------------
- Go 1.3 or greater
- Mercurial 2.9 or greater

How to Build
-------------

You can build Walter with the follwoing commands.

    $git clone https://github.com/recruit-tech/walter.git
    $cd walter
    $./build

How to contribute
====================

We welcome any contributions through Github pull requests.
When you make changes such as new features and the add the tests, please run test before throw the pull reqest.
You can run test with the test.sh script.

    $sh test.sh

Configuration setting
======================

Walter has one configuration file, which specifies a set of tasks needed to build or deploy target applicaion or service. 
More specifically, users specify the order of the tasks for the deployment. Each task is called as **stage**, and the flow is called **pipeline** in Walter.

## Pipeline setting

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

### Stage setting

Stage in pipeline has three elements, **name** **stage_type** and the **configurations** if needed. The configurations depend on the stage_type.
For example command_stage type has **command** configuration, which specify the shell command run in the stage.
The following is the table on the stage_type and the parameters.

#### Command stage

Command stage executes one command. Users specify Command stage adding **command** in stage_type.

The follwoing is the parameter of Command stage.

|  Configuration | Optional   | meaning                                |
|:--------------:|:----------:|:--------------------------------------:|
|   command      | false      | shell command run in the stage         |

#### Shell script stage
Command stage executes specified shell scrpit file. Users specify Command stage adding **shell_script** in stage_type.

The follwoing is the parameter of Command stage.

|  Configuration   | Optional   | meaning                                |
|:----------------:|:----------:|:--------------------------------------:|
|   shell _script  | false      | shell script file run in the stage     |

## Report setting
Walter supports to submits the messages to messeging service (currently only HipChat is supported). To submit a messege,
users need to add a **messenger** block into the configuraiton file. Following is a sample of the yaml block.

    messenger:
        type: hipchat
        room_id: ROOM_ID
	    token: TOKEN
		from: USER_NAME

To report the result to the messenger service added with the above setting,
users add **messege** attribute with **true** into the stage they want to know the results.

     pipeline:
          command_stage_1:
             stage_type: command
             command: echo "hello, world"
			 message: true
         command_stage_2:
             stage_type: command
             command: echo "hello, world, command_stage_2"
			 message: true
