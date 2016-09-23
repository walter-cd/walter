# Walter v2

## pipeline.yml

```yaml
build:
  tasks:
    - name: build 1
      command: echo "hello w"
    - name: task 2
      command: echo "hello world, build 2"
    - name: parallel tasks
      parallel:
        - name: parallel build 1
          command: parallel command 1
        - name: parallel build 2
          serial:
            - name: serial build 1
              command: serial build 1 in parallel build 2
            - name: serial build 2
              command: serial build 2 in parallel build 2
        - name: parallel build 3
          command: parallel 
  cleanup:
    - name: cleanup
      command: command for cleanup

deploy:
  tasks:
    - name: deploy 1
      command: echo "hello world, deploy 1"
  cleanup:
    - name cleanup
      command: command for cleanup
```

