all: test

test:
	go run main.go demo/flowchart.mmd

.PHONY: resources
resources:
	pkger