package main


import (
	"fmt"
	"google.golang.org/grpc"
	 "github.com/dgraph-io/dgo/v2"
	 "github.com/dgraph-io/dgo/v2/protos/api"
	"context"
	"encoding/json"
	)

	

	type Book struct {
		Uid string `json:"uid", omitempty`
		Tile string `json:"title", omitempty`
		Language string `json:"language", omitempty`
	}

	

	type LibraryPlugins interface {
		Add(*Book) (error)
		Delete(*Book)(error)
		Profile()
	}


	type Librarian struct {}

	var(
		PublishBooks []Book = []Book{}
	)

	type Root struct {
		Books []Book `json:"books"`
	}


func main() {
	
	connection , err := grpc.Dial(":9080", grpc.WithInsecure())	
	
	if err != nil {
		fmt.Println("Timeout", err)
		return 
	}

	defer connection.Close()
	
	fmt.Println("Grpc Connection:", connection)
	
	client := dgo.NewDgraphClient(api.NewDgraphClient(connection))
	
	fmt.Println("Client :", client)
	
	num := 0
	// InitBooks(num)
	PublishBooks = make([]Book, num)
	reads := Book{}
	reads.Uid = "_:Zero"
	reads.Tile =  "Zero to One"
	reads.Language = "English" 
	rowBooks := make([]Book, num)
	rowBooks = append(rowBooks , NewBook(reads))
	PublishBooks = append(PublishBooks, rowBooks[num])
	fmt.Println("Books:", PublishBooks)
	err = DBoperations(client); if err != nil{
		fmt.Println("Timeout operation", err)
		return
	}
	mutate := &api.Mutation{
		CommitNow : true,
	}

	data , err := json.Marshal(PublishBooks[num])
	if err != nil{
		fmt.Println("Timeout marshal", err)
		return
	}
	fmt.Println("Data:", string(data))
	mutate.SetJson = data

	resp , err := client.NewTxn().Mutate(context.Background(), mutate)
	if err != nil {
		fmt.Println("Timeout mutate", err)
		return		
	}

	fmt.Println("Response:", resp)
	vari := map[string]string{"$id1": resp.Uids["Zero"]}
	dgraphquery := `query Book($id1: string){
		book(func: uid($id1)){
			title
			langauge
		}
	}`
	
	queryRes , err := client.NewTxn().QueryWithVars(context.Background(), dgraphquery, vari)
	if err != nil {
		fmt.Println("Timeout query result", err)
		return
	}
	fmt.Println("Query Result:", queryRes)
	schema := Root{} 
	err = json.Unmarshal(queryRes.Json, &schema)
	if err != nil {
		fmt.Println("Timeout unmarshal", err)
		return
	}
	fmt.Println(string(queryRes.Json))
	queryresp , err := client.NewTxn().Query(context.Background(), `schema(pred:[title, language]) {type}`)
	if err != nil {
		fmt.Println("Timeout unmarshal", err)
		return
	}
	fmt.Println("Query Execute:", string(queryresp.Json))
}


func DBoperations (client *dgo.Dgraph) (error){
	err := client.Alter(context.Background(), &api.Operation{
		Schema: `
			Id : string .
			Title : string @index(fulltext) .
			Language : string @index(fulltext,exact) .
		` ,
	})
	return err
}


func NewBook(book Book) Book {
	return Book{
		Uid : book.Uid,
		Tile : book.Tile,
		Language : book.Language,
	}
}