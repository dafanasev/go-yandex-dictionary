go-yandex-dictionary
====================

Yandex Dictionary Go (golang) API wrapper

Usage:

```
package main

import (
  "fmt"
  "github.com/icrowley/go-yandex-dictionary"
)

func main() {
  dict := yandex_dictionary.New("YOUR_API_KEY")

  langs, err := dict.GetLangs()

  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println(langs)
  }

  definition, err := dict.Lookup(&yandex_dictionary.Params{Lang: "en-ru", Text: "Dog"})

  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println(definition.Def[0].Text)
  }
}
```
