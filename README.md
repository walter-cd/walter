Walter
========

<p align="center">
<img src="https://dl.dropboxusercontent.com/u/10177896/walter-logo-readme.png"/>
</p>

[![wercker status](https://app.wercker.com/status/4fcb4b110909fc45775d12641f5cf037/m "wercker status")](https://app.wercker.com/project/bykey/4fcb4b110909fc45775d12641f5cf037)

Walter is a tiny deployment pipeline template.

Blogs
==========

* http://ainoya.io/walter
* http://walter-cd.net

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

In the above setting, `parallel command 1`, `parallel command 2` and `parallel command 3` are executed in parallel.

## Import predefined stages

Walter provides the feature to import predefined stages in specified files.
To import stages to pipeline configuration file, we use **require** block and
add the list of file names into the block.

For example, the following example import the stages defined in **conf/mystage.yml**

```yaml
require:
    - conf/mystages.yml

pipeline:
  - call: mypackage::hello
  - call: mypackage::goodbye

```

In the above setting, the stages ("mypacakge::hello" and "mypackage::goodbye") which are defined in "mystage.yml" are specified.

The files specified in pipeline configuration file need to have two blocks **namespace** and **stages**.
In **namespace**, we add the package name, the package name is need to avoid collisions of the stage names in multiple required files.
The **stages** block contains the list of stage definitions. We can define the stages same as the stages in pipeline configurations.

For example, the following configuration is the content of conf/mystages.yml imported from the above pipeline configuration file.

```yaml
namespace: mypackage

stages:
  - def:
      name: hello
      command: echo "May I help you majesty!"
  - def:
      name: goodbye
      command: echo "Goobye majesty."
```

As we see that stages **hello** and **goodbye** are defined in the file.

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
  type: hipchat2
  base_url: BASE_URL
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
|   base_url       |  Base API URL (for private servers; hipchat2 only) |

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

|  Element  | description                                                                                                                          |
|:---------:|:-------------------------------------------------------------------------------------------------------------------------------------|
|   type    |  Service type (currently Walter supports github only)                                                                                |
|   token   |  [GitHub token](https://help.github.com/articles/creating-an-access-token-for-command-line-use/)                                     |
|   repo    |  Repository name                                                                                                                     |
|   from    |  Account or organization name (if the repository is own by a organization)                                                           |
|   update  |  Update file which contains the result and time of the last execution  (default **.walter**)                                         |
|   branch  |  Branch name pattern (When this value is filled, Walter only checks the branches who name is matched with the filled regex pattern)  |

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

## Reusing the results from stages

Walter stores the results of preceding stages. The stages can make use of the results of finished stages using the three special
variables (__OUT, __ERR and __RESULT) in Walter configuration files.

* **__OUT** -  output flushed to standard output 
* **__ERR** - output flushed to standard error 
* **__RESULT** - execution result (true or false)

The three variables are maps whose keys are stage names and the value are results of the stages. For example, we want the standard output
result of the stage named "stage1", we write __OUT["stage1"].

The following is a sample configuration with a specical value.

```yaml
pipeline:
  - name: stage_1
    command: echo "hello world"
  - name: stage_2
    command: echo __OUT["stage_1"]
```

Walter with the above configuraiton outputs "hello world" twice, since the second stage (stage_2) flushes the standard output result of first stage (stage_1).

## Wait running stages until the conditions are satisfied

Walter stage starts imidiately after the previous stage finish, but some stages need to wait for some action such as port is ready or file are created.
```wait_for``` feature supports the actions which need to be ready before the stages begin.

wait_for is defined as a property of stage.

```yaml
- pipeline:
  - name: launch solr
    command: bin/solr start
  - name: post data to solr index
    command: bin/post -d ~/tmp/foobar.js
    wait_for: host=localhost port=8983 state=open
```

The **wait_for** property takes the **key** **value** pairs. Key has sevaral variations. The value depends on the key type. The following table shows the supported key value pairs and the description.


| Key     | Value (value type)  | Description                                         |
|:--------|:--------------------|:----------------------------------------------------|
| delay   | second (float)      | The second to wait after the previous stage finish |
| port    | port number (int)   | Port number                                         |
| file    | file name (string)  | File to be created in the previous stages           |
| host    | host (string)       | IP address or host name                             |
| state   | state of the other key (string) | Four types of states are supported. The possible value is dependent to the other Keys. |

State values
--------------

There are seveal **state** values and possible state values are depend on the other key.

| State value       | Description                                         |
|:------------------|:----------------------------------------------------|
| exist / ready     | Specified port is ready or file is created.         |
| delete / absent   | port is not active or file does not exist           |
