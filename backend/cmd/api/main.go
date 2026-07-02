package main

import(
	"fmt"
	"net/http"
)

func main(){
	http.HandleFunc("/api/health", func(w http.ResponseWriter,r *http.Request){
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "El Psy Kongroo")
	})

	fmt.Println("Serwer działa na porcie 8081")
	if err := http.ListenAndServe(":8081", nil); err !=nil{
		panic(err)
	}
}