# appleopensource

An [opensource.apple.com][opensource.apple.com] resource management package and command line tool written in Go.

| **CI (darwin)**                             | **codecov.io**                          | **godoc.org**                      | **Analytics**                |
|:-------------------------------------------:|:---------------------------------------:|:----------------------------------:|:----------------------------:|
| [![circleci.com][circleci-badge]][circleci] | [![codecov.io][codecov-badge]][codecov] | [![godoc.org][godoc-badge]][godoc] | [![Analytics][ga-badge]][ga] |

===

## Install

Installing `appleopensource` command:

```sh
go get -u -v github.com/zchee/appleopensource/cmd/appleopensource
```


## Usage

```
```


## API

See [godoc.org/github.com/zchee/appleopensource](https://godoc.org/github.com/zchee/appleopensource).


## Background


## Contribute

Not yet.  
~~See [CONTRIBUTING.md](CONTRIBUTING.md)~~.


## Acknowledgement

- [Apple Open Source][opensource.apple.com]


## License

appleopensource is released under the [BSD 3-Clause License](https://opensource.org/licenses/BSD-3-Clause).  
[Apple Open Source][opensource.apple.com] project resources is under the [Apple Public Source License Version 2.0][apsl].



[opensource.apple.com]: https://opensource.apple.com
[apsl]: http://www.opensource.apple.com/apsl/

[circleci]: https://circleci.com/gh/zchee/appleopensource
[godoc]: https://godoc.org/github.com/zchee/appleopensource
[codecov]: https://codecov.io/gh/zchee/appleopensource
[release]: https://github.com/zchee/nvim-go/releases
[ga]: https://github.com/zchee/appleopensource

[circleci-badge]: https://img.shields.io/circleci/project/github/zchee/appleopensource.svg?style=flat-square&label=%20%20CircleCI&logoWidth=16&logo=data%3Aimage%2Fsvg%2Bxml%2C%3Csvg+xmlns%3D%27http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%27+width%3D%2740%27+viewBox%3D%270+0+200+200%27%3E%3Cpath+fill%3D%27%23DDD%27+d%3D%27M74.7,100c0-13.2,10.7-23.8,23.8-23.8c13.1,0,23.8,10.7,23.8,23.8c0,13.1-10.7,23.8-23.8,23.8C85.4,123.8,74.7,113.1,74.7,100zM98.5,0C51.8,0,12.7,32,1.6,75.2c-0.1,0.3-0.1,0.6-0.1,1c0,2.6,2.1,4.8,4.8,4.8h40.3c1.9,0,3.6-1.1,4.3-2.8l0,0c8.3-18,26.5-30.6,47.6-30.6c28.9,0,52.4,23.5,52.4,52.4c0,28.9-23.5,52.4-52.4,52.4c-21.1,0-39.3-12.5-47.6-30.6l0,0c-0.8-1.6-2.4-2.8-4.3-2.8H6.3c-2.6,0-4.8,2.1-4.8,4.8c0,0.3,0.1,0.6,0.1,1C12.6,168,51.8,200,98.5,200c55.2,0,100-44.8,100-100C198.5,44.8,153.7,0,98.5,0z%27%2F%3E%3C%2Fsvg%3E
[godoc-badge]: https://img.shields.io/badge/godoc-reference-4F73B3.svg?style=flat-square&label=%20godoc.org
[codecov-badge]: https://img.shields.io/codecov/c/github/zchee/appleopensource.svg?style=flat-square&label=%20%20Codecov%2Eio&logo=data%3Aimage%2Fsvg%2Bxml%3Bbase64%2CPHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI0MCIgaGVpZ2h0PSI0MCIgdmlld0JveD0iMCAwIDI1NiAyODEiPjxwYXRoIGZpbGw9IiNkZGQiIHN0cm9rZT0iI2RkZCIgZD0iTTIxOC41NTEgMzcuNDE5QzE5NC40MTYgMTMuMjg5IDE2Mi4zMyAwIDEyOC4wOTcgMCA1Ny41MzcuMDQ3LjA5MSA1Ny41MjcuMDQgMTI4LjEyMUwwIDE0OS44MTNsMTYuODU5LTExLjQ5YzExLjQ2OC03LjgxNCAyNC43NS0xMS45NDQgMzguNDE3LTExLjk0NCA0LjA3OSAwIDguMTk4LjM3MyAxMi4yNCAxLjExIDEyLjc0MiAyLjMyIDI0LjE2NSA4LjA4OSAzMy40MTQgMTYuNzU4IDIuMTItNC42NyA0LjYxNC05LjIwOSA3LjU2LTEzLjUzNmE4OC4wODEgODguMDgxIDAgMCAxIDMuODA1LTUuMTVjLTExLjY1Mi05Ljg0LTI1LjY0OS0xNi40NjMtNDAuOTI2LTE5LjI0NWE5MC4zNSA5MC4zNSAwIDAgMC0xNi4xMi0xLjQ1OSA4OC4zNzcgODguMzc3IDAgMCAwLTMyLjI5IDYuMDdjOC4zNi01MS4yMjIgNTIuODUtODkuMzcgMTA1LjIzLTg5LjQwOCAyOC4zOTIgMCA1NS4wNzggMTEuMDUzIDc1LjE0OSAzMS4xMTcgMTYuMDExIDE2LjAxIDI2LjI1NCAzNi4wMzMgMjkuNzg4IDU4LjExNy0xMC4zMjktNC4wMzUtMjEuMjEyLTYuMS0zMi40MDMtNi4xNDRsLTEuNTY4LS4wMDdhOTAuOTU3IDkwLjk1NyAwIDAgMC0zLjQwMS4xMTFjLTEuOTU1LjEtMy44OTguMjc3LTUuODIxLjUtLjU3NC4wNjMtMS4xMzkuMTUzLTEuNzA3LjIzMS0xLjM3OC4xODYtMi43NS4zOTUtNC4xMDkuNjM5LS42MDMuMTEtMS4yMDMuMjMxLTEuOC4zNTFhOTAuNTE3IDkwLjUxNyAwIDAgMC00LjExNC45MzdjLS40OTIuMTI2LS45ODMuMjQzLTEuNDcuMzc0YTkwLjE4MyA5MC4xODMgMCAwIDAtNS4wOSAxLjUzOGMtLjEuMDM1LS4yMDQuMDYzLS4zMDQuMDk2YTg3LjUzMiA4Ny41MzIgMCAwIDAtMTEuMDU3IDQuNjQ5Yy0uMDk3LjA1LS4xOTMuMTAxLS4yOTMuMTUxYTg2LjcgODYuNyAwIDAgMC00LjkxMiAyLjcwMWwtLjM5OC4yMzhhODYuMDkgODYuMDkgMCAwIDAtMjIuMzAyIDE5LjI1M2MtLjI2Mi4zMTgtLjUyNC42MzUtLjc4NC45NTgtMS4zNzYgMS43MjUtMi43MTggMy40OS0zLjk3NiA1LjMzNmE5MS40MTIgOTEuNDEyIDAgMCAwLTMuNjcyIDUuOTEzIDkwLjIzNSA5MC4yMzUgMCAwIDAtMi40OTYgNC42MzhjLS4wNDQuMDktLjA4OS4xNzUtLjEzMy4yNjVhODguNzg2IDg4Ljc4NiAwIDAgMC00LjYzNyAxMS4yNzJsLS4wMDIuMDA5di4wMDRhODguMDA2IDg4LjAwNiAwIDAgMC00LjUwOSAyOS4zMTNjLjAwNS4zOTcuMDA1Ljc5NC4wMTkgMS4xOTIuMDIxLjc3Ny4wNiAxLjU1Ny4xMDQgMi4zMzhhOTguNjYgOTguNjYgMCAwIDAgLjI4OSAzLjgzNGMuMDc4LjgwNC4xNzQgMS42MDYuMjc1IDIuNDEuMDYzLjUxMi4xMTkgMS4wMjYuMTk1IDEuNTM0YTkwLjExIDkwLjExIDAgMCAwIC42NTggNC4wMWM0LjMzOSAyMi45MzggMTcuMjYxIDQyLjkzNyAzNi4zOSA1Ni4zMTZsMi40NDYgMS41NjQuMDItLjA0OGE4OC41NzIgODguNTcyIDAgMCAwIDM2LjIzMiAxMy40NWwxLjc0Ni4yMzYgMTIuOTc0LTIwLjgyMi00LjY2NC0uMTI3Yy0zNS44OTgtLjk4NS02NS4xLTMxLjAwMy02NS4xLTY2LjkxNyAwLTM1LjM0OCAyNy42MjQtNjQuNzAyIDYyLjg3Ni02Ni44MjlsMi4yMy0uMDg1YzE0LjI5Mi0uMzYyIDI4LjM3MiAzLjg1OSA0MC4zMjUgMTEuOTk3bDE2Ljc4MSAxMS40MjEuMDM2LTIxLjU4Yy4wMjctMzQuMjE5LTEzLjI3Mi02Ni4zNzktMzcuNDQ5LTkwLjU1NCIvPjwvc3ZnPg==
[release-badge]: https://img.shields.io/github/release/zchee/appleopensource.svg?style=flat-square
[ga-badge]: https://ga-beacon.appspot.com/UA-89201129-1/gist-go-template?flat&useReferer
