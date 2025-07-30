package datasources

import (
	"os"

	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

func NewPincecone() (*pinecone.IndexConnection, error) {

	pinceconeApiKey := os.Getenv("PINECONE_API_KEY")
	indexHostName := os.Getenv("PINECONE_INDEX_HOSTNAME")

	var indexConn *pinecone.IndexConnection

	pc, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey: pinceconeApiKey,
	})
	if err != nil {
		return indexConn, err
	}

	indexConn, err = pc.Index(pinecone.NewIndexConnParams{Host: indexHostName})
	if err != nil {
		return indexConn, err
	}

	return indexConn, nil

}
