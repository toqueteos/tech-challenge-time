# Pento Tech Challenge

I'm used to organizing Go APIs following _Mat Ryer's style_. He has iterated this style over the years, the latest version is explained in his talk at [GopherCon Europe 2019](https://www.youtube.com/watch?v=8TLiGHJTlig).

## Requirements

- Docker 19.03+
- docker-compose (which usually comes along Docker in most OSes, but here's a [link][dc] just in case)
- A bash-compatible shell (git bash on Windows, bash/zsh/fish on Linux/OSX).

[dc]: https://github.com/docker/compose

## Quickstart

- Open up a terminal at this project's root

```bash
$ docker-compose up -d
```

- Go to <http://localhost:8080> in your prefered web browser
- That's it!

There's two special endpoints:

- `GET` <http://localhost:8080/_/debug/vars> to see expvar metrics (simplified Prometheus builtin into Go's stdin)
- `POST` <http://localhost:8080/timer/_fake> to generate some fake data to play around with

## How did I work with this?

Although I could just do `go run cmd/backend/main.go` (from the project's root) I've used [modd](https://github.com/cortesi/modd/releases) for live reloading.

Download the binary, put it in your path and run `modd` from the project's root.

## What about the frontend?

Well... There's nothing special required, I just put all frontend code inside `static/main.jsx`.

I decided NOT to use a bundler (rollup/webpack) or a template like `create-react-app` to _keep things simple_.

The end result is almost 400 lines of JavaScript in a single file which is not that simple after all but I aimed for a days work and in the end I couldn't clean it up.
