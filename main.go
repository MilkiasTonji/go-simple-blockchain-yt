package main

// import packages
import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// create stucks

type Block struct {
	Pos       int
	Data      BookCheckout
	TimesTamp string
	Hash      string
	PrevHash  string
}

type BookCheckout struct {
	BookID       string `json:"book_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"`
}

type Book struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedDate string `json:"published_date"`
	ISBN          string `json:"isbn"`
}

type BlockChain struct {
	blocks []*Block
}

var BlockChain *BlockChain

// define main func
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", getBlockChain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newBook).Methods("POST")

	log.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", r))
}
