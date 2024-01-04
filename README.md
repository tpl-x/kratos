# kratos-tpl
my go kratos Project Template base on offical kratos layout

## Features
 + use [go-task](https://github.com/go-task/task) rather than make
 + use [buf](https://github.com/bufbuild/buf) for proto build
 + built in zap with lumbjack
 + use [goreleaser](https://github.com/goreleaser/goreleaser) to cross build
## required tools
goreleaser
```bash
go install github.com/goreleaser/goreleaser@latest
```
Task
```bash
go install github.com/go-task/task/v3/cmd/task@latest
```
for Buf installation, please refer [document](https://buf.build/docs/installation)
## usage
run the command:
```
kratos new <your App Name> -r https://github.com/tpl-x/kratos.git
```
or
```
kratos new <your App Name> -r git@github.com:tpl-x/kratos.git
```
to create your first Application
