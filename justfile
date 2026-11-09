set shell := ["bash", "-cu"]

default: help

help:
	@just --list

bin:
	go build -o jikascript ./cmd/jikascript

install:
	go install ./cmd/jikascript

[arg("file", long)]
check file="examples/blink.jika": bin
	./jikascript check {{file}}

[arg("file", long)]
[arg("dest", long)]
transpile file="examples/blink.jika" dest="out/blink.c": bin
	mkdir -p "$(dirname "{{dest}}")"
	./jikascript transpile {{file}} --output={{dest}}

[arg("source", long)]
[arg("board", long)]
[arg("optimize", long)]
build source="examples/blink.jika" board="ch32v003" optimize="size": bin
	./jikascript build {{source}} --board={{board}} --optimize={{optimize}}

[arg("source", long)]
[arg("board", long)]
[arg("optimize", long)]
build-c source="out/blink.c" board="ch32v003" optimize="size": bin
	./jikascript build {{source}} --board={{board}} --optimize={{optimize}}

[arg("source", long)]
[arg("board", long)]
[arg("port", long)]
flash source="examples/blink.jika" board="ch32v003" port="": bin
	./jikascript flash {{source}} --board={{board}} {{ if port != "" { "--port=" + port } else { "" } }}
