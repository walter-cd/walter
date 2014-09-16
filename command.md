---
layout: default
title: Command
---

This page describes the Plumber command usage.
Plumber provides the **plumber** command, which run the build or deployment process definded in the input configuration file.

quick start
-------------

As the first step we build plumber.

    $git clone https://github.com/recruit-tech/plumber.git
    $cd plumber
    $./build

Then, we run plumber with the sample configuration file included in the plumber project source.

    $bin/plumber -c tests/fixtures/pipeline.yml

Options
---------

The plubmer command supports the follwoing options

      -alsologtostderr=false: log to standard error as well as files
      -c="./pipeline.yml": pipeline.yml file
      -f=false: Skip execution of subsequent stage after failing to exec the upstream stage
      -log_backtrace_at=:0: when logging hits line file:N, emit a stack trace
      -log_dir="": If non-empty, write log files in this directory
      -logtostderr=false: log to standard error instead of files
      -stderrthreshold=0: logs at or above this threshold go to stderr
      -v=0: log level for V logs
      -vmodule=: comma-separated list of pattern=N settings for file-filtered logging

