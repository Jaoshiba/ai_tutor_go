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
	UpsertVector(chapter entities.RefDataModel, co *cohereClient.Client, ctx *fiber.Ctx) error
}

func NewPineconeRepository(pineconeIdxConn *pinecone.IndexConnection) IPineconeRepository {
	return &PineconeRepository{
		pineconeIdxConn: pineconeIdxConn,
	}
}

func (pc *PineconeRepository) UpsertVector(ref entities.RefDataModel, co *cohereClient.Client, ctx *fiber.Ctx) error {
	fmt.Println("upsertVector call...")
	// limit := uint32(100)
	// namespaces, err := indexConn.ListNamespaces(ctx.Context(), &pinecone.ListNamespacesParams{
	// 	Limit: &limit,
	// })
	// if err != nil {
	// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err, "message": "Failed to list namespaces."})
	// }

	//get userId to create namespace
	userIdRaw := ctx.Locals("userID")
	if userIdRaw == nil {
		fmt.Println("Error: User ID not found in context locals for DocSearchService.")
		return fiber.NewError(fiber.StatusUnauthorized, "User ID not found in context")
	}
	userIdStr, ok := userIdRaw.(string)
	if !ok || userIdStr == "" {
		fmt.Println("Error: Invalid or missing user ID format in context locals for DocSearchService.")
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or missing user ID")
	}

	userId := userIdStr
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

	courseId := ctx.Locals("courseID")

	vectorId := uuid.NewString()
	metadata, err := structpb.NewStruct(map[string]interface{}{
		"refId":    ref.RefId,
		"moduleId": ref.ModuleId,
		"courseId": courseId,
		"Title":    ref.Title,
		"Link":     ref.Link,
		"searchat": ref.SearchAt,
	})
	if err != nil {
		fmt.Println("error create metadata: ", err)
		return err
	}

	denseValue, err := EmbeddingText(co, ref.Content, ctx)
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

func EmbeddingText(co *cohereClient.Client, text string, ctx *fiber.Ctx) ([]float32, error) {
	resp, err := co.V2.Embed(
		ctx.Context(),
		&cohere.V2EmbedRequest{
			Texts:          []string{text},
			Model:          "embed-multilingual-v3.0",
			InputType:      cohere.EmbedInputTypeSearchDocument,
			EmbeddingTypes: []cohere.EmbeddingType{cohere.EmbeddingTypeFloat},
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]float32, len(resp.Embeddings.Float[0]))
	for i, v := range resp.Embeddings.Float[0] {
		result[i] = float32(v)
	}
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
