# GoKvSqlite
* simply, it just key-value Sqlite based store package, so simply and easy-to-use

* Example:
```golang
package main

import (
	"log"

	store "github.com/tm-sah/GoKvSqlite"
)

func main() {

	db, err := store.Client("db.sql")
	if err != nil {
		log.Fatal(err)
	}
	db.Set(
		"name",
		"Mahdi",
	)
	db.Set(
		"age:1",
		16,
	)

	keys, err := db.Keys("*") // or just "" it will work, also.
	if err != nil {
		log.Fatal(err)
	}
	for index, key := range keys {
		println(index, key)
	}
}

```
# Available methods:
``` Keys, Get, Set, Exists, Delete ```
  
* for simple porpuses.
* it accept all type of storing, int, int64, string, map, slices, etc.
  
