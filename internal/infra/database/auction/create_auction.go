package auction

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/diogocardoso/go_lab3/configuration/logger"
	"github.com/diogocardoso/go_lab3/internal/entity/auction_entity"
	"github.com/diogocardoso/go_lab3/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
	closeChan  chan bool
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	repo := &AuctionRepository{
		Collection: database.Collection("auctions"),
		closeChan:  make(chan bool),
	}

	// Iniciando a goroutine de verificação
	go repo.checkExpiredAuctions()

	return repo
}

// Função para calcular o tempo do leilão baseado em variáveis de ambiente
func (ar *AuctionRepository) calculateAuctionDuration() time.Duration {
	// Obtendo duração das variáveis de ambiente (exemplo: "24h")
	durationStr := os.Getenv("AUCTION_DURATION")
	if durationStr == "" {
		durationStr = "24h" // Valor padrão
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		logger.Error("Error parsing auction duration", err)
		return time.Hour * 24 // Valor padrão em caso de erro
	}

	return duration
}

// Goroutine para verificar leilões expirados
func (ar *AuctionRepository) checkExpiredAuctions() {
	checkInterval := time.Minute // Intervalo de verificação
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ar.closeChan:
			return
		case <-ticker.C:
			ctx := context.Background()

			// Busca leilões ativos que já expiraram
			currentTime := time.Now().Unix()
			filter := bson.M{
				"status":    auction_entity.Active,
				"timestamp": bson.M{"$lt": currentTime - int64(ar.calculateAuctionDuration().Seconds())},
			}

			update := bson.M{
				"$set": bson.M{"status": auction_entity.Completed},
			}

			// Atualiza todos os leilões expirados
			result, err := ar.Collection.UpdateMany(ctx, filter, update)
			if err != nil {
				logger.Error("Error updating expired auctions", err)
				continue
			}

			if result.ModifiedCount > 0 {
				logger.Info(fmt.Sprintf("Closed %d expired auctions", result.ModifiedCount))
			}
		}
	}
}

// Método para parar a goroutine (útil para testes)
func (ar *AuctionRepository) Stop() {
	close(ar.closeChan)
}

func (ar *AuctionRepository) CreateAuction(ctx context.Context, auction *auction_entity.Auction) *internal_error.InternalError {
	// Definindo o timestamp de criação
	auction.Timestamp = time.Now()

	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auction.Id,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   auction.Condition,
		Status:      auction.Status,
		Timestamp:   auction.Timestamp.Unix(),
	}

	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to create auction", err)
		return internal_error.NewInternalServerError("Error trying to create auction")
	}

	return nil
}
