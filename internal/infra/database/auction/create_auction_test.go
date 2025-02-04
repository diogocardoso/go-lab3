package auction

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/diogocardoso/go_lab3/internal/entity/auction_entity"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDatabase(t *testing.T) (*mongo.Database, func()) {
	// Configura conexão com MongoDB de teste
	ctx := context.Background()
	mongoURI := "mongodb://localhost:27017"

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Usar um banco de dados específico para testes
	dbName := "test_auction_db"
	database := client.Database(dbName)

	// Retornar função de cleanup
	cleanup := func() {
		// Limpar a coleção após o teste
		database.Collection("auctions").Drop(ctx)
		client.Disconnect(ctx)
	}

	return database, cleanup
}

func TestAutomaticAuctionClosure(t *testing.T) {
	// ambiente de teste
	os.Setenv("AUCTION_DURATION", "1s")

	// Inicializa banco de dados de teste
	database, cleanup := setupTestDatabase(t)
	defer cleanup()

	// Criar repositório
	repo := NewAuctionRepository(database)
	defer repo.Stop()

	// Criar um leilão de teste
	auction := &auction_entity.Auction{
		Id:          "test-auction-id",
		ProductName: "Test Product",
		Category:    "Test Category",
		Description: "Test Description",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now().Add(-2 * time.Second), // Já expirado
	}

	// Criar o leilão
	ctx := context.Background()
	err := repo.CreateAuction(ctx, auction)
	assert.NoError(t, err, "Não deveria haver erro ao criar o leilão")

	// Esperar tempo suficiente para o processo de fechamento executar
	time.Sleep(3 * time.Second)

	// Verificar se o leilão foi fechado
	var result AuctionEntityMongo
	mongoErr := repo.Collection.FindOne(ctx, bson.M{"_id": auction.Id}).Decode(&result)
	assert.NoError(t, mongoErr, "Não deveria haver erro ao buscar o leilão")

	// Verificar se o status foi atualizado para fechado
	assert.Equal(t, auction_entity.Completed, result.Status, "O leilão deveria estar fechado")
}

// Teste adicional para verificar criação de leilão
func TestCreateAuction(t *testing.T) {
	// Inicializar banco de dados de teste
	database, cleanup := setupTestDatabase(t)
	defer cleanup()

	// Criar repositório
	repo := NewAuctionRepository(database)
	defer repo.Stop()

	// Criar um leilão de teste
	auction := &auction_entity.Auction{
		Id:          "test-create-auction",
		ProductName: "Test Product",
		Category:    "Test Category",
		Description: "Test Description",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	}

	// Criar o leilão
	ctx := context.Background()
	err := repo.CreateAuction(ctx, auction)
	assert.NoError(t, err, "Não deveria haver erro ao criar o leilão")

	// Verificar se o leilão foi criado corretamente
	var result AuctionEntityMongo
	mongoErr := repo.Collection.FindOne(ctx, bson.M{"_id": auction.Id}).Decode(&result)
	assert.NoError(t, mongoErr, "Não deveria haver erro ao buscar o leilão")
	assert.Equal(t, auction.Id, result.Id)
	assert.Equal(t, auction.ProductName, result.ProductName)
}
