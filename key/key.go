package main

import (
	"github.com/gorilla/securecookie"
	"fmt"
	"encoding/hex"
)


func main() {
	key := securecookie.GenerateRandomKey(32)
	fmt.Println(hex.EncodeToString(key))
}
