defmodule GlowServer.MixProject do
  use Mix.Project

  def project do
    [
      app: :example,
      version: "0.1.0",
      elixir: "~> 1.12",
      start_permanent: Mix.env() == :prod,
      deps: deps(),
      erlc_paths: ["src", "gen"],
      compilers: [:gleam | Mix.compilers()], # Gleam must go first
    ]
  end

  def application do
    [
      extra_applications: [:logger]
    ]
  end

  defp deps do
    [
      {:gleam_stdlib, "~> 0.16.0"},
      {:mix_gleam, "~> 0.1.0"},
    ]
  end
end
