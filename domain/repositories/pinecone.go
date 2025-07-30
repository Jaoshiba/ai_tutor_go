package repositories

import (
	"fmt"
	"go-fiber-template/domain/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"google.golang.org/protobuf/types/known/structpb"

	cohere "github.com/cohere-ai/cohere-go/v2"
	cohereClient "github.com/cohere-ai/cohere-go/v2/client"
)

type PineconeRepository struct {
	pineconeIdxConn *pinecone.IndexConnection
}

type IPineconeRepository interface {
	UpsertVector(chapter entities.ChapterDataModel, co *cohereClient.Client, ctx *fiber.Ctx) error
}

func NewPineconeRepository(pineconeIdxConn *pinecone.IndexConnection) IPineconeRepository {
	return &PineconeRepository{
		pineconeIdxConn: pineconeIdxConn,
	}
}

func (pc *PineconeRepository) UpsertVector(chapter entities.ChapterDataModel, co *cohereClient.Client, ctx *fiber.Ctx) error {
	fmt.Println("upsertVector call...")
	// limit := uint32(100)
	// namespaces, err := indexConn.ListNamespaces(ctx.Context(), &pinecone.ListNamespacesParams{
	// 	Limit: &limit,
	// })
	// if err != nil {
	// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err, "message": "Failed to list namespaces."})
	// }

	//get userId to create namespace
	// userIdRaw := ctx.Locals("userId")
	// userId := userIdRaw.(string)
	userId := uuid.NewString()
	namespace := userId + "-namespace"

	fmt.Println("namespace: ", namespace)

	// namespacesList := namespaces.Namespaces
	// for _, namespace := range namespacesList {
	// 	if namespace.Name == userId+"-namespace" {
	// 		return nil
	// 	}
	// }

	//this user have no namespace yet
	// denseValue := make([]float32, pineconeDimension)
	//add initial record to namespace to create namespace

	vectorId := uuid.NewString()
	metadata, err := structpb.NewStruct(map[string]interface{}{
		"chapterid":   chapter.ChapterId,
		"chaptername": chapter.ChapterName,
		"userid":      chapter.UserID,
		"courseid":    chapter.CouseId,
	})
	if err != nil {
		fmt.Println("error create metadata: ", err)
		return err
	}

	//create denseValue from chapter content
	denseValue, err := EmbeddingText(co, chapter.ChapterContent, ctx)
	if err != nil {
		fmt.Println("error create denseValue: ", err)
		return err
	}

	fmt.Println("denseValue: ", denseValue)

	upsertBody := entities.UpsertBodyChapter{
		Vectors: []*pinecone.Vector{
			{
				Id:       vectorId,
				Values:   &denseValue,
				Metadata: metadata,
			},
		},
	}

	namespaceConn := pc.pineconeIdxConn.WithNamespace(namespace)
	_, err = namespaceConn.UpsertVectors(ctx.Context(), upsertBody.Vectors)
	if err != nil {
		fmt.Println("error upsert vectors: ", err)
		return err
	}

	return nil
}

func EmbeddingText(co *cohereClient.Client, text string, ctx *fiber.Ctx) ([]float64, error) {
	resp, err := co.V2.Embed(
		ctx.Context(),
		&cohere.V2EmbedRequest{
			Texts:          []string{text},
			Model:          "embed-v4.0",
			InputType:      cohere.EmbedInputTypeSearchDocument,
			EmbeddingTypes: []cohere.EmbeddingType{cohere.EmbeddingTypeFloat},
		},
	)
	if err != nil {
		return nil, err
	}

	fmt.Println("resp: ", resp)
	result := resp.Embeddings.Float
	return result, nil
}

/*
vectors := []*pinecone.Vector{
            {
            	Id:           "abc-1",
                Values:       &denseValues,
                Metadata:     metadata,
                SparseValues: &sparseValues,
            },
    }
*/
