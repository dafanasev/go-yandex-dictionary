[![Build Status](https://travis-ci.org/dafanasev/go-yandex-dictionary.svg?branch=master)](https://travis-ci.org/dafanasev/go-yandex-dictionary)
[![GoDoc](https://godoc.org/github.com/dafanasev/go-yandex-dictionary?status.svg)](https://godoc.org/github.com/dafanasev/go-yandex-dictionary)
[![Go Report Card](https://goreportcard.com/badge/github.com/dafanasev/go-yandex-dictionary)](https://goreportcard.com/report/github.com/dafanasev/go-yandex-dictionary)
[![Coverage Status](https://coveralls.io/repos/github/dafanasev/go-yandex-dictionary/badge.svg)](https://coveralls.io/github/dafanasev/go-yandex-dictionary)

go-yandex-dictionary
====================

Go Yandex Dictionary API wrapper

Usage:

```
package main

import (
  "fmt"
  "github.com/dafanasev/go-yandex-dictionary"
)

func main() {
  dict := dictionary.New("YOUR_API_KEY")

  langs, err := dict.GetLangs()

  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println(langs)
  }

  definition, err := dict.Lookup(&dictionary.Params{Lang: "en-ru", Text: "Dog"})

  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println(definition.Def[0].Text)
  }
}
```
