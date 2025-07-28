package datasources

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

func NewPincecone(ctx *fiber.Ctx) (*pinecone.Client, error) {

	pinceconeApiKey := os.Getenv("PINECONE_API_KEY")

	pc, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey: pinceconeApiKey,
	})
	if err != nil {
		return pc, ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err, "message": "Failed to create Pinecone client."})
	}

	return pc, nil

}
