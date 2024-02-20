# Lazy Deploy

Lazy Deploy is a tool which helps you publish multiple Gitlab branches by building a dependency network. Imagine the problem:

- In order to deploy service "A", you must deploy services "B" and "C".
- Service "C", requires service "D" to be deployed.
- Service "B" has no dependencies.

If you are lazy or forgetful like me, you end up wasting a bunch of time in this process. So, I've written this CLI tool to automate this task.

# Installation / Build

- Clone this project.
- Build or run with Go.
    - Place it somewhere in FS, I'd be honored if you include us in your `$PATH`.

# Usage

Please check the help documentation by running the program without any arguments: `Lazy-Deploy`
