# `gleam_pb`

WIP protobuf support for Gleam ✨

## Progress

- [X] Gleam Type generation
  - [ ] custom functions that better handle default values
  - [ ] stop including unnecessary imports
  - [ ] `gleam format` generated files
- [ ] message encoding
- [ ] message decoding
- [ ] grpc

## Example

This message

```protobuf
syntax = "proto3";

package protos;

message Example {
    
    message Response {
        int32 val = 1;
        string user = 2;
    }

    enum OptionType {
        foo = 0;
        bar = 1;
        baz = 2;
    }

    repeated OptionType options = 1;

    oneof ResponseOrError {
        Response response = 2;
        string error = 3;
    }
}
```

Becomes

```rust
import gleam/option
import gleam/list
import gleam/map
import gleam_pb

pub type OptionType {
  OptionTypefoo
  OptionTypebar
  OptionTypebaz
}

pub type ExampleResponseOrError {
  ExampleResponseOrErrorresponse(response: option.Option(Response))
  ExampleResponseOrErrorerror(error: String)
}

pub type Example {
  Example(
    options: list.List(OptionType),
    response_or_error: ExampleResponseOrError,
  )
}

pub type Response {
  Response(val: Int, user: String)
}
```
