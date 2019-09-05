gostudy
=======

Some people learning [Go](https://golang.org).

Setup
-----

Install Go:

```sh
brew install go
```

Make sure the Go bin directory is in your path (you'll probably want to put this in a shell startup script):

```sh
export PATH="$HOME/go/bin:$PATH"
```

Install some useful tools:

```sh
brew install httpie
go get golang.org/x/tools/cmd/goimports
go get golang.org/x/lint/golint
```

Install or enable Go support in your editor/IDE; I use https://github.com/fatih/vim-go

Clone this repo:

```sh
git clone https://github.com/cultureamp/gostudy
```

The `master` branch is empty, ready for you to proceed. You'll find various other branches with example code.
