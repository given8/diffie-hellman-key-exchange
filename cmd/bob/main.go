package main

import (
	"bytes"
	"diffie-hellman/pkg"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

func main() {

	http.HandleFunc("/key", func(w http.ResponseWriter, r *http.Request) {
		var aliceSharedKey pkg.SharedKey
		err := json.NewDecoder(r.Body).Decode(&aliceSharedKey)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Printf("received alices key %v\n", aliceSharedKey)

		const privateNum = 5
		sharedKey := pkg.SharedKey{
			Key: math.Mod(math.Pow(pkg.PublicNumA, privateNum), pkg.PublicNumB),
		}

		fmt.Printf("Bobs sharedKey %v\n", sharedKey)

		body, err := json.Marshal(sharedKey)

		if err != nil {
			fmt.Print(err)
			panic("could not marshal shared bobs sharedKey")
		}
		_, err = http.Post("http://localhost:3000/key", "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		finalKey := math.Mod(math.Pow(aliceSharedKey.Key, privateNum), pkg.PublicNumB)

		fmt.Printf("final key %v", finalKey)

	})

	http.ListenAndServe(":4000", nil)
}
