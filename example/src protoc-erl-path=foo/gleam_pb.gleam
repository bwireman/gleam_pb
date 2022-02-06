import gleam/list
import gleam/dynamic
import gleam/option
import gleam/erlang/atom

// helper funcs and types
pub type Undefined(a) {
  Undefined
  Wrapper(a)
}

pub fn encode(m: dynamic.Dynamic, name: atom.Atom) -> BitString {
  encode_msg(m, name, [True])
}

pub fn option_to_gpb(o: option.Option(a)) -> dynamic.Dynamic {
  case o {
    option.Some(v) ->
      v
      |> dynamic.from
    option.None ->
      Undefined
      |> dynamic.from
  }
}

pub fn wrapper_to_option(w: Undefined(a)) -> option.Option(a) {
  case w {
    Undefined -> option.None
    Wrapper(v) ->
      v
      |> option.Some
  }
}

pub fn force_a_to_b(a: a) -> b {
  a
  |> dynamic.from
  |> dynamic.unsafe_coerce
}

external fn encode_msg(dynamic.Dynamic, atom.Atom, List(Bool)) -> BitString =
  "gleam_gpb" "encode_msg"
