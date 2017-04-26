Walter
======

<p align="center">
<img src="https://dl.dropboxusercontent.com/u/10177896/walter-logo-readme.png"/>
</p>

[![wercker status](https://app.wercker.com/status/4fcb4b110909fc45775d12641f5cf037/m "wercker status")](https://app.wercker.com/project/bykey/4fcb4b110909fc45775d12641f5cf037)

Walter is a tiny deployment pipeline tool.

----

Blogs
=====

* http://ainoya.io/walter
* http://walter-cd.net


----

Overview
========

Walter is a simple command line tool that automates build, test and deployment of applications or servers.

----

Getting Started
===============


How to install
--------------

Get a binary file from [GitHub Releases](https://github.com/walter-cd/walter/releases) and place it in `$PATH`.


Wrirte your pipeline
--------------------

Write your command pipeline in `pipeline.yml`.


```yaml
build:
  tasks:
    - name: setup build
      command: echo "setting up ..."
    - name: run build
      command: echo "building ..."
  cleanup:
    - name: cleanup build
      command: echo "cleanup build ..."

deploy:
  tasks:
    - name: run deploy
      command: echo "deploying ..."
  cleanup:
    - name: cleanup
      command: echo "cleanup deploy ..."
```


Run walter
----------

```
$ walter -build -deploy
INFO[0000] Build started
INFO[0000] [setup build] Start task
INFO[0000] [setup build] setting up ...
INFO[0000] [setup build] End task
INFO[0000] [run build] Start task
INFO[0000] [run build] building ...
INFO[0000] [run build] End task
INFO[0000] Build succeeded
INFO[0000] Build cleanup started
INFO[0000] [cleanup build] Start task
INFO[0000] [cleanup build] cleanup build ...
INFO[0000] [cleanup build] End task
INFO[0000] Build cleanup succeeded
INFO[0000] Deploy started
INFO[0000] [run deploy] Start task
INFO[0000] [run deploy] deploying ...
INFO[0000] [run deploy] End task
INFO[0000] Deploy succeeded
INFO[0000] [cleanup] Start task
INFO[0000] [cleanup] cleanup deploy ...
INFO[0000] [cleanup] End task
INFO[0000] Deploy cleanup succeeded
```

That's it.

----

Other features
==============

Environment variables
---------------------

You can use environment variables.

```yaml
deploy:
  tasks:
    - name: release files
      command: ghr -token $GITHUB_TOKEN $VERSION pkg/dist/$VERSION
```


Working directory
-----------------

You can specify a working directory of a task.

```yaml
build:
  tasks:
    - name: list files under /tmp
      command: ls
      directory: /tmp
```

Conditions to run tasks
----------

```yaml
build:
  tasks:
    - name: list files under /tmp
      command: ls
      only_if: test -d /tmp
```


Get stdout of a previous task
-----------------------------

Tasks get stdout of a previous task through a pipe.

```yaml
build:
  tasks:
    - name: setup build
      command: echo "setting up"
    - name: run build
      command: cat
```

The second "run build" task outputs "setting up".


Parallel tasks
--------------

You can define parallel tasks.

```yaml
build:
   tasks:
     - name: parallel tasks
       parallel:
           - name: task 1
             command: echo task 1
           - name: task 2
             command: echo task 2
           - name: task 3
             command: echo task 3
```

You can also mix serial tasks in parallel tasks.

```yaml
build:
   tasks:
     - name: parallel tasks
       parallel:
           - name: task 1
             command: echo task 1
           - name: task 2
             serial:
               - name: task 2-1
                 command: echo task 2-1
               - name: task 2-2
                 command: echo task 2-2
           - name: task 3
             command: echo task 3
```

   
Split pipeline definitions and include them
-------------------------------------------

You can split pipeline definitions in other files and include them.

**pipeline.yml**

```yaml
build:
  tasks:
    - include: task1.yml
    - include: task2.yml
```

**task1.yml**

```yaml
- name: task1
  command: echo task1
```

**task2.yaml**

```yaml
- name: task2
  command: echo task2
```

You can also run single definition file.

```
$ walter -build -config task2.yml
```


Wait for some conditions
------------------------

You can make tasks wait for some conditions.

```yaml
build:
  tasks:
    - name: launch solr
      command: bin/solr start
    - name: post data to solr index
      command: bin/post -d ~/tmp/foobar.js
      wait_for:
        host: localhost
        port: 8983
        state: ready
```

Available keys and values are these:


| Key     | Value (value type)  | Description                                         |
|:--------|:--------------------|:----------------------------------------------------|
| delay   | second (float)      | Seconds to wait after the previous stage finish     |
| host    | host (string)       | IP address or host name                             |
| port    | port number (int)   | Port number                                         |
| file    | file name (string)  | File name|
| state   | state of the other key (string) | Two types(present/ready or absent/unready) of states are supported. |



Notification
------------

Walter supports notification of task results to Slack.

```yaml
notify:
  - type: slack
    channel: serverspec
    url: $SLACK_WEBHOOK_URL
    icon_url: http://example.jp/walter.jpg
    username: walter
```

Other services(ex. HipChat) are not supported currently.

----

Changes in v2
=============


Pipeline definition format
--------------------------

Pipeline definition in v1:

```yaml
pipeline:
  - name: setup build
    type: command
    command: echo "setting up ..."
  - name: run build
    type: command
    command: echo "building ..."
cleanup:
  - name: cleanup build
    type: command
    command: echo "cleanup build ..."

```

In v2:

```yaml
build:
  tasks:
    - name: setup build
      command: echo "setting up ..."
    - name: run build
      command: echo "building ..."
  cleanup:
    - name: cleanup build
      command: echo "cleanup build ..."
```

Separate build and deploy phase
-------------------------------

You can separate build and deploy phases in v2:

```yaml
build:
  tasks:
    - name: setup build
      command: echo "setting up ..."
    - name: run build
      command: echo "building ..."
  cleanup:
    - name: cleanup build
      command: echo "cleanup build ..."

deploy:
  tasks:
    - name: run deploy
      command: echo "deploying ..."
  cleanup:
    - name: cleanup
      command: echo "cleanup deploy ..."
```

You can run both phases at once or each phase separately.

```
# Run build and deploy phases
$ walter -build -deploy

# Run build phase only
$ walter -build

# Run deploy phase only
$ walter -deploy
```

Format of wait_for
------------------

You must define parameters for `wait_for` in one line in v1:

```yaml
pipeline:
  - name: launch solr
    command: bin/solr start
  - name: post data to solr index
    command: bin/post -d ~/tmp/foobar.js
    wait_for: host=localhost port=8983 state=ready
```

In v2, you must define parameters for `wait_for` with mapping of yaml.

```
build:
  tasks:
    - name: launch solr
      command: bin/solr start
    - name: post data to solr index
      command: bin/post -d ~/tmp/foobar.js
      wait_for:
        host: localhost
        port: 8983
        state: ready
```


Definition of notification
--------------------------

In v1:

```yaml
messenger:
  type: slack
  channel: serverspec
  url: $SLACK_WEBHOOK_URL
  icon_url: http://example.jp/walter.jpg
  username: walter
```

In v2:

```yaml
notify:
  - type: slack
    channel: serverspec
    url: $SLACK_WEBHOOK_URL
    icon_url: http://example.jp/walter.jpg
    username: walter
```

The key `messenger` was changed to `notify` and you can define multiple notification definitions in v2.


Special variables
-----------------

Special variables like `__OUT`, `__ERR`, `__COMBINED` and `__RESULT` are obsoleted in v2.

Tasks get stdout of a previous task through a pipe.

```yaml
build:
  tasks:
    - name: setup build
      command: echo "setting up"
    - name: run build
      command: cat
```

The second "run build" task outputs "setting up".

I think this is suffient for defining pipelines. Special variables bring complexity for pipelines.


----

Contributing
============

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

