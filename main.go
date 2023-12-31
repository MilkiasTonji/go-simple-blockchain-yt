package main

// import packages
import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

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

type Blockchain struct {
	blocks []*Block
}

// blockchain variable to store blocks like a database
var BlockChain *Blockchain

// generate hash and add to block
func (b *Block) generateHash() {
	bytes, _ := json.Marshal(b.Data)

	data := string(b.Pos) + b.TimesTamp + string(bytes) + b.PrevHash

	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))

}

// create a block
func CreateBlock(prevBlock *Block, checkoutItem BookCheckout) *Block {
	block := &Block{}
	block.Pos = prevBlock.Pos + 1
	block.TimesTamp = time.Now().String()
	block.PrevHash = prevBlock.Hash
	block.generateHash()
	return block
}

// struct method is different from a regular function in go.
// struct method
func (bc *Blockchain) AddBlock(data BookCheckout) {
	prevBlock := bc.blocks[len(bc.blocks)-1] // take current block

	block := CreateBlock(prevBlock, data)
	if validBlock(block, prevBlock) {
		bc.blocks = append(bc.blocks, block)
	}
}

func validBlock(block, prevBlock *Block) bool {

	if prevBlock.Hash != block.Hash {
		return false
	}
	if !block.validateHash(block.Hash) {
		return false
	}

	if prevBlock.Pos+1 != block.Pos {
		return false
	}
	return true
}

func (b *Block) validateHash(hash string) bool {
	b.generateHash()
	if b.Hash != hash {
		return false
	}
	return true
}

func writeBlock(w http.ResponseWriter, r *http.Request) {
	var checkoutItem BookCheckout
	err := json.NewDecoder(r.Body).Decode(&checkoutItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not create a block %v: ", err)
		w.Write([]byte("could not create a block "))
		return
	}
	// call struct method
	BlockChain.AddBlock(checkoutItem)

}

// request and response because we are using mux
func newBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	// encode json data to struct
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not crate: %v", err)
		w.Write([]byte("could not create new book"))
		return
	}

	// create new Id for book
	h := md5.New()
	io.WriteString(h, book.ISBN+book.PublishedDate)
	book.ID = fmt.Sprintf("%x", h.Sum(nil))

	resp, err := json.MarshalIndent(book, "", " ")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not marshal payload: %v", err)
		w.Write([]byte("could not save book data"))
		return
	}

	// set status and create a new book
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GenesisBlock() *Block {
	return CreateBlock(&Block{}, BookCheckout{IsGenesis: true})
}

func NewBlockChain() *Blockchain {
	return &Blockchain{[]*Block{GenesisBlock()}}
}

// get blockchains

func getBlockChain(w http.ResponseWriter, r *http.Request) {
	jbytes, err := json.MarshalIndent(BlockChain.blocks, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	io.WriteString(w, string(jbytes))

}

// define main func
func main() {

	BlockChain = NewBlockChain()

	r := mux.NewRouter()
	r.HandleFunc("/", getBlockChain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newBook).Methods("POST")

	// function executes as soon as the program starts
	go func() {
		for _, block := range BlockChain.blocks {
			fmt.Printf("Prev hash: %x\n", block.PrevHash)
			bytes, _ := json.MarshalIndent(block.Data, "", " ")
			fmt.Printf("Data: %v\n", string(bytes))
			fmt.Printf("Hash: %x\n", block.Hash)
			fmt.Println()
		}
	}()

	log.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", r))
}
