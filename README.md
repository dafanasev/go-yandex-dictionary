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
