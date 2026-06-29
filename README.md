# kratos-tpl
my [go-kratos](https://github.com/go-kratos/kratos) Project Template base on offical kratos layout
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

AIP Go

This template tracks [aip-go](https://github.com/einride/aip-go) as a Go tool dependency and uses
`go tool protoc-gen-go-aip` from `buf.gen.yaml`, so the generator version is kept in `go.mod`.
The generated `*_aip.go` resource-name helpers are produced when your protobuf APIs use
standard [Google AIP](https://aip.dev/) `google.api.resource` annotations. Runtime helpers
for pagination, filtering, ordering, and field masks are available from `go.einride.tech/aip`.

`aip-go` also provides helpers for common AIP patterns:

| Package | AIP | Purpose |
| --- | --- | --- |
| `pagination` | AIP-132 | Parse and generate offset-based pagination tokens. |
| `filtering` | AIP-160 | Parse filter expressions such as `filter=name="foo"`. |
| `ordering` | AIP-132 | Parse `order_by` strings. |
| `fieldmask` | AIP-134 / AIP-161 | Validate and apply field masks for update APIs. |
| `fieldbehavior` | AIP-203 | Handle field behavior annotations such as `REQUIRED` and `OUTPUT_ONLY`. |
| `resourcename` | AIP-122 | Validate and parse resource names such as `shelves/123/books/456`. |
| `resourceid` | AIP-122 | Generate and validate resource IDs. |
| `validation` | General | Validate common request parameters. |

How to use these helpers in this template:

1. Define the public API shape in `proto/`.
   - Use `google.api.resource` on resource messages when you want generated resource-name helpers.
   - Add standard request fields such as `page_size`, `page_token`, `filter`, `order_by`, and `update_mask` when your API needs List or Update behavior.
   - Keep HTTP routes in `google.api.http`; the Kratos HTTP plugin will keep generating the route handlers.
2. Run `buf generate`.
   - Protobuf, gRPC, HTTP, validation, OpenAPI, and AIP helper files are generated into `api/` and `docs/`.
   - `*_aip.go` files come from `protoc-gen-go-aip` and are generated from resource annotations.
3. Use the helpers in `internal/service` or `internal/biz`.
   - Resource names: use generated helpers such as `v1.ShelfResourceName.UnmarshalString(request.GetName())`.
   - Pagination: use `pagination.ParsePageToken(request)` and set `response.NextPageToken = pageToken.Next(request).String()` when there is another page.
   - Filtering: build declarations with `filtering.NewDeclarations(...)`, then call `filtering.ParseFilter(request, declarations)`.
   - Ordering: call `ordering.ParseOrderBy(request)`, then validate allowed paths with `orderBy.ValidateForPaths(...)` or `orderBy.ValidateForMessage(...)`.
   - Update masks: validate with `fieldmask.Validate(request.GetUpdateMask(), request.GetResource())`, then apply with `fieldmask.Update(...)`.
   - Field behavior: use `fieldbehavior.ClearFields(...)`, `ValidateRequiredFields(...)`, or `ValidateImmutableFieldsNotChanged(...)` around create and update requests.
   - Resource IDs: use `resourceid.ValidateUserSettable(id)` for user-provided IDs or `resourceid.NewSystemGenerated()` for server-generated IDs.

The AIP helpers parse, validate, and normalize API-level input. They do not build SQL or storage queries for you; convert the parsed filter, order, page token, or field mask into repository query parameters in your own `internal/biz` or `internal/data` code.

Example:

`proto/helloworld/v1/greeter.proto` defines a `Shelf` resource with the pattern
`shelves/{shelf}` and a `GetShelf` API. After `buf generate`, `api/helloworld/v1/greeter_aip.go`
contains the generated `ShelfResourceName` helper. Run the service and request
`GET /v1/shelves/default` to see the example response.

kratos
```bash
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest && kratos upgrade
```
if you want to build all platform simply,you may need install [goreleaser](https://github.com/goreleaser/goreleaser),and build with
```bash
goreleaser build --snapshot --clean
```
also you can create a release in your github,the compile is auto completed
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
>  if you want to clean up unused dep. use `buf dep prune `
> 
> step 2 :generate api from protobuf
> ```bash
> buf generate
>   ```
