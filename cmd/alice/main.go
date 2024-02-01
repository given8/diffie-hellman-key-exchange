package main

import (
	"bytes"
	"diffie-hellman/pkg"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

const privateNum = 13

func main() {
	http.HandleFunc("/start", func(writer http.ResponseWriter, r *http.Request) {
		sharedKey := math.Mod(math.Pow(pkg.PublicNumA, privateNum), pkg.PublicNumB)

		fmt.Printf("Alice's shared key %v\n", sharedKey)

		body, err := json.Marshal(pkg.SharedKey{
			Key: sharedKey,
		})

		if err != nil {
			fmt.Print(err)
			panic("could not marshal shared key")
		}
		_, err = http.Post("http://localhost:4000/key", "application/json", bytes.NewBuffer(body))
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("could not send shared key to bob"))
			return
		}
		writer.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/key", func(writer http.ResponseWriter, r *http.Request) {
		var bobKey pkg.SharedKey
		err := json.NewDecoder(r.Body).Decode(&bobKey)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("Error converting bob's key to a float"))
			return
		}
		finalKey := math.Mod(math.Pow(bobKey.Key, privateNum), pkg.PublicNumB)

		fmt.Printf("final key %v\n", finalKey)
	})

	http.ListenAndServe(":3000", nil)
	fmt.Println("alice is listening")
}
