# gotrew

[![Build Status](https://travis-ci.org/proshik/gotrew.svg?branch=master)](https://travis-ci.org/proshik/gotrew)
[![Go Report Card](https://goreportcard.com/badge/github.com/proshik/gotrew)](https://goreportcard.com/report/github.com/proshik/gotrew)
[![codecov](https://codecov.io/gh/proshik/jalmew/branch/master/graph/badge.svg)](https://codecov.io/gh/proshik/gotrew)
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/proshik/gotrew/issues)

Console Application for translate words. Support english and russian languages.

![Animation](https://github.com/proshik/gotrew/blob/master/public/intro.gif)

## Installation

```console
$ go get github.com/proshik/gotrew
$ cd $GOPATH/src/github.com/proshik/gotrew
$ go install
```

```console
$ gotrew --help
```

## Usage manual

```console
$ gotrew

NAME:
   gotrew - Application for translate words. Support english and russian languages.

USAGE:
   gotrew [global options] command [command options] [arguments...]

VERSION:
   0.1.5

COMMANDS:
     translate, t  translate words mode
     provider, p   show and select provider for translate
     help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

Command <provider> for work with available providers translated.

```console
$ gotrew.exe provider

NAME:
   gotrew provider - show and select provider for translate

USAGE:
   gotrew provider command [command options] [arguments...]

COMMANDS:
     list    list available providers
     select  select a provider
     config  set config for providers

OPTIONS:
   --help, -h  show help
```

Need select one provider from list, and configure him. (Now support only Yandex Dictionary API and chosen by `default`)

```console
$ gotrew provider list
```

You may print follow command for help

```console
$ gotrew provider config -h

NAME:
   gotrew.exe provider config - set config for providers

USAGE:
   gotrew.exe provider config [command options] [arguments...]

OPTIONS:
   --token value, -t value  provider dictionary token
```

#### `--token`
Specifies token for access to external service of translating.

Example configure follow next 

```console
$ gotrew provider config yandex --token=<yandex_dictionary_token>
```

## TODO

- tests; 
- translate not only words, but and texts;
- integration with few other translate providers.

## Patch 

Welcome!