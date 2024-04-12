package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", BuscaCepHandler)

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

}

func BuscaCepHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Erro de request!"))
		return
	}

	cep, error := BuscaCep(cepParam)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cep)
}

func BuscaCep(cep string) (*ViaCEP, error) {
	resp, error := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if error != nil {
		return nil, error
	}
	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}

	var c ViaCEP
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}
	return &c, nil
}
