# Lazy Deploy

Lazy Deploy is a tool which helps you publish multiple Gitlab branches by building a dependency network. Imagine the problem:

- In order to deploy service "A", you must deploy services "B" and "C".
- Service "C", requires service "D" to be deployed.
- Service "B" has no dependencies.

If you are lazy or forgetful like me, you end up wasting a bunch of time in this process. So, I've written this CLI tool to automate this task.

# Installation

Download a binary for your operating system in the releases page or check the "Building from source" steps.

# Building from source 

- Make sure you have installed Go.
- Clone this project.
- `cd` into the project directory and run `go build`.

# Usage

Please check the help documentation by running the program without any arguments: `lazy-deploy`
