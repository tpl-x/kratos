# kratos-tpl
my go kratos Project Template base on offical kratos layout
> The BSR allows 10 unauthenticated code generation requests per hour, 
 with a burst of up to 10 requests. If you send more than 10 unauthenticated 
 requests per hour using remote plugins, you’ll receive a rate limit error.
 https://buf.build/docs/bsr/rate-limits
> 
## Features
 + use [go-task](https://github.com/go-task/task) rather than make
 + use [buf](https://github.com/bufbuild/buf) for proto build
 + built in zap with lumbjack
 + use [goreleaser](https://github.com/goreleaser/goreleaser) to cross build
## related tools
goreleaser (Optional)
```bash
go install github.com/goreleaser/goreleaser@latest
```
Task (Optional)
```bash
go install github.com/go-task/task/v3/cmd/task@latest
```
Buf
```bash
go install github.com/bufbuild/buf/cmd/buf@latest
```
for more info, please refer [buf document](https://buf.build/docs/installation)

wire
```bash
go install github.com/google/wire/cmd/wire@latest
```
kratos
```bash
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest && kratos upgrade
```
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
> to generate api only using buf ,you can follow the steps below:
> 
> step 1:update buf dep
> ```bash
> buf dep update
>```
> step 2 :generate api from protobuf
> ```bash
> buf generate
   ```
