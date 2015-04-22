Walter
========

<p align="center">
<img src="https://dl.dropboxusercontent.com/u/10177896/walter-logo-readme.png"/>
</p>

[![wercker status](https://app.wercker.com/status/4fcb4b110909fc45775d12641f5cf037/m "wercker status")](https://app.wercker.com/project/bykey/4fcb4b110909fc45775d12641f5cf037)

Walter is a tiny deployment pipeline template.

Overview
==========

Walter automates the deployment process of applications or servers from configuration to software deployment.

Getting Started
===============

Requirements
-------------
- Go 1.3 or greater
- Mercurial 2.9 or greater

How to Build
-------------

You can build Walter with the following commands.

```
$ git clone git@github.com:walter-cd/walter.git
$ cd walter
$ ./build
```

How to contribute
====================

We welcome any contributions through Github pull requests.
When you make changes such as new features and the add the tests, please run test before throw the pull request.
You can run test with the test.sh script.

```
$ sh test.sh
```

Configuration setting
======================

Walter has one configuration file, which specifies a set of tasks needed to build or deploy target application or service.
More specifically, users specify the order of the tasks for the deployment. Each task is called as **stage**, and the flow is called **pipeline** in Walter.

## Pipeline setting

The configuration format of Walter is Yaml. The yaml configuration file need to have one pipeline block, which has more
than one stage element.

The following is a sample configuration of Walter.

```yaml
pipeline:
  - name: command_stage_1
    type: command
    command: echo "hello, world"
  - name: command_stage_2
    type: command
    command: echo "hello, world, command_stage_2"
  - name: command_stage_3
    type: command
    command: echo "hello, world, command_stage_3"
```

As we see, the pipeline block has three stages and the stage type is command, each of which run **echo** command and has the stage name
(such as **command_stage_1**). User can name arbitrary name of each stage. The commands are executed with the same order as the pipeline configuration.

### Stage setting

Stage in pipeline has three elements, **name**, **type** and **configurations**. configuration elements are optional. The elements of configurations depend on the type.
For example command_stage type has **command** configuration, which specify the shell command run in the stage.
The following is the table on the type and the parameters.

#### Command stage

Command stage executes one command. Users specify Command stage adding **command** in type.

The following is the parameter of Command stage.

|  Configuration | Optional   | meaning                                                                     |
|:--------------:|:----------:|:----------------------------------------------------------------------------|
|   command      | false      | shell command run in the stage                                              |
|   only_if      | true       | run specified command on when the condition written in only_if is satisfied |
|   directory    | true       | the directory where walter runs the specified command                       |

#### Shell script stage
Shell script stage executes specified shell script file. Users specify Shell script stage adding **shell** in type.

The following is the parameter of Shell script stage.

|  Configuration   | Optional   | meaning                                |
|:----------------:|:----------:|:--------------------------------------:|
|   file           | false      | shell script file run in the stage     |

## Parallel stages

 You can set child stages and run these stages in parallel like this.
 
```yaml
pipeline:
  - name: parallel stages
    parallel:
      - name: parallel command 1
        type: command
        command: parallel command 1
      - name: parallel command 2
        type: command
        command: parallel command 2
      - name: parallel command 3
        type: command
        command: parallel command 3
```

`parallel command 1`, `parallel command 2` and `parallel command 3` are executed in parallel.

## Cleanup pipeline

Walter configuraiton can have one **cleanup** block; cleanup is another pipeline which needs to be executed after a pipeline has either failed or passed.
In the cleanup block, we can add command or shell script stages. The below example create a log file in pipeline and then cleanup the log file in the cleaup steps.


```yaml
pipeline:
  - name: start pipeline
    command: echo “pipeline” > log/log.txt
cleanup:
  - name: cleanup
    command:  rm log/*
```


## Reporting function
Walter supports to submits the messages to messaging services.

### Report configuration
To submit a message, users need to add a **messenger** block into the configuration file. The following is a sample of the yaml block with HipChat.

```yaml
messenger:
  type: hipchat
  room_id: ROOM_ID
  token: TOKEN
  from: USER_NAME
```

To report the full output of stage execution to the specified messenger service added with the above setting,
users add **report_full_output** attribute with **true** into the stage they want to know the command outputs.

```yaml
pipeline:
  - name: command_stage_1
    type: command
    command: echo "hello, world"
    report_full_output: true # If you put this option, messenger sends the command output "hello, world"
  - name: command_stage_2
    type: command
    command: echo "hello, world, command_stage_2"
    # By default, report_full_output is false
```

### Report types
Walter supports HipChat API v1 and v2 as the messenger type.

|  Messenger Type  | meaning                                                                                   |
|:----------------:|:-----------------------------------------------------------------------------------------:|
|   hipchat        |  [HipChat (API v1)](https://www.hipchat.com/docs/api)                                     |
|   hipchat2       |  [HipChat (API v2)](https://www.hipchat.com/docs/apiv2)                                   |
|   slack          |  [Slack Incoming Webhook integration](https://my.slack.com/services/new/incoming-webhook) |

### Report configuration
To activate the report function, we need to specify the properties for messenger type. The needed properties are different in each messenger type.

#### hipchat and hipchat2

|  Property Name   | meaning                                                                                   |
|:----------------:|:----------------------------------|
|   room_id        |  Room name                        |
|   token          |  HipChat token                    |
|   from           |  Account name                     |

#### slack

|  Property Name   | meaning                                                                                   |
|:----------------:|:----------------------------------|
|   channel        |  Channel name                     |
|   username       |  User name                        |

## Service coordination

Walter provides a coordination function to a project hosting service, GitHub. Specifically the service function provides two roles

- Check if a new commit or pull requests from the repository
- Run pipeline to the latest commit and pull requests when there are newer ones than the last update

### Service configuration
To activate service coordination function, we add a "service" block to the Walter configuration file. "service" block contains several elements (type, token, repo, from, update).

```yaml
service:
  type: github
  token: ADD_YOUR_KEY
  repo: YOUR_REPOSITORY_NAME
  from: YOUR_ACCOUNT_OR_GROUP_NAME
  update: UPDATE_FILE_NAME
```

The following shows the description of each element.

|  Element  | description                                                                                          |
|:---------:|:----------------------------------------------------------------------------------------------------|
|   type    |  Service type (currently Walter supports github only)                                                |
|   token   |  [GitHub token](https://help.github.com/articles/creating-an-access-token-for-command-line-use/)  |
|   repo    |  Repository name                                                                                     |
|   from    |  Account or organization name (if the repository is own by a organization)                           |
|   update  |  Update file which contains the result and time of the last execution                                |

## Embedding Environment Variables

Users add environment variables in Walter configuration files. The names of environment variables are expanted into the the values of environment variables.
Environment variables in the configuration files are valuable when we need to write the sensitive information such as
tokens of messenger service or passwords of a external systems.

The following is the format of embedding of the environment variables.

```
$ENV_NAME
```

We write the envrionment variable into **ENV_NAME**. The following configuration file specify the GitHub Token by embedding the environment variable, GITHUB_TOKEN.

```yaml
service:
  type: github
  token: $GITHUB_TOKEN
  repo: my-service-repository
  from: service-group
  update: .walter
```
