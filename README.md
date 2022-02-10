# gleam_pb

Protobuf support for Gleam ✨

---

gleam_pb wraps the excellent [gpb](https://github.com/tomas-abrahamsson/gpb) erlang library and generates idiomatic Gleam types 🤘

## Progress

- [X] Gleam Type generation
  - [X] custom functions that better handle default values
  - [X] gleam format generated files
  - [ ] stop including unnecessary imports
- [X] message encoding
- [X] message decoding
- [ ] improve UX
  - [X] call protoc-erl internally
  - [ ] helper functions
- [ ] grpc

## API

### Generated types

gleam_pb generally follows gpb's type generation, but makes it easier to use from Gleam.

| protobuf | gleam_pb | gpb |
|---|---|---|
| double,float | Float | float() |
| int32, int64, uint32, uint64, sint32, sint64, fixed32, fixed64, sfixed32, sfixed64 | Int | integer() |
| bool | Bool | true \| false |
| enum | Zero Paramater Multi Constructor Type | atom() |
| message | Option(\<CustomType\>) | record \| undefined |
| string | String | unicode string |
| bytes | BitString | binary() |
| oneof | Option(\<CustomType\>) with multiple constructors | {chosen_field, value} |
| map | unordered list of tuples List(#(Key, Value)) | [{key, value}] |

### Functions

gleam_pb generates functions to make using the types easier

- function to generate the message with protobuf's default values named new_\<custom_type\>() -> \<CustomType\>

- functions to encode and decode the messages
  - encode_\<custom_type\>(m: \<CustomType\>) -> BitString
  - encode_\<custom_type\>(b: BitString) -> \<CustomType\>

There are also several other functions intended for usage by gleam_pb

### Example

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
import gleam/pair
import gleam/dynamic
import gleam/erlang/atom
import gleam_pb

/// protos package types generated by gleam_pb
/// DO NOT EDIT
pub type OptionType {
  OptionTypefoo
  OptionTypebar
  OptionTypebaz
}

pub type ExampleResponseOrError {
  ResponseOrErrorresponse(response: option.Option(Response))
  ResponseOrErrorerror(error: String)
}

pub type Example {
  Example(
    options: List(OptionType),
    response_or_error: option.Option(ExampleResponseOrError),
  )
}

pub type Response {
  Response(val: Int, user: String)
}

pub fn new_example() {
  Example(list.new(), option.None)
}

pub fn new_response() {
  Response(0, "")
}

pub fn encode_example(m: Example) -> BitString {
  let name = atom.create_from_string("protos.Example")

  extract_example(name, m)
  |> gleam_pb.encode(name)
}

pub fn decode_example(b: BitString) -> Example {
  let name = atom.create_from_string("protos.Example")
  decode_msg_example(b, name)
  |> reconstruct_example
}

pub fn encode_response(m: Response) -> BitString {
  let name = atom.create_from_string("protos.Example.Response")

  extract_response(name, m)
  |> gleam_pb.encode(name)
}

pub fn decode_response(b: BitString) -> Response {
  let name = atom.create_from_string("protos.Example.Response")
  decode_msg_response(b, name)
  |> reconstruct_response
}

//internal functions continue ...
```

## Usage

gleam_pb and gpb must be used together to generate working Gleam code.

Example Script

```bash
# make sure protoc-gen-gleam is in you're path or add it manually using --plugin
protoc --plugin=protoc-gen-gleam -I . --gleam_out="output_path=./src:./src" protos/*.proto
```

### `gleam_pb` Flags

- 'output_path': (Required) specifies the desired output path
- 'protoc_erl_path': path to gpb's protoc-erl
- 'gpb_header_include': path to prepend to the header include for gpb. See [Issues](#Issues) for more info
  - if you need a variable include here, remember that [erlang header resolution](https://www.erlang.org/doc/reference_manual/macros.html) is quite clever and can use environment variables

```bash
protoc -I . \
  --gleam_out="gpb_header_include=$ENV/include/,output_path=./src,protoc_erl_path=bin/protoc-erl:./src" \
  protos/*.proto
```

### Known Issues

#### Includes aren't working?!

```erlang
% generated in `gleam_gpb.erl`
-include("gpb.hrl"). % -> update to point to the correct header post Gleam compilation
```
