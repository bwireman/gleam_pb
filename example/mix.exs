defmodule Example.MixProject do
  use Mix.Project

  @app :example

  def project do
    [
      app: @app,
      archives: [mix_gleam: "~> 0.4.0"],
      aliases: MixGleam.add_aliases(aliases()),
      erlc_paths: ["build/dev/erlang/#{@app}/build"],
      erlc_include_path: "build/dev/erlang/#{@app}/include",
      version: "0.1.0",
      elixir: "~> 1.12",
      start_permanent: Mix.env() == :prod,
      deps: deps()
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      extra_applications: [:logger]
    ]
  end

  def aliases do
    []
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      {:gleam_stdlib, "~> 0.19"},
      {:gpb, "~> 4.19"},
      {:gleam_erlang, "~> 0.7.0"},
      {:gleeunit, "~> 0.6", only: [:dev, :test], runtime: false},
    ]
  end
end
