package auction_controller

import (
	"context"
	"log"
	"net/http"

	"github.com/diogocardoso/go_lab3/configuration/logger"
	"github.com/diogocardoso/go_lab3/configuration/rest_err"
	"github.com/diogocardoso/go_lab3/internal/entity/auction_entity"
	"github.com/diogocardoso/go_lab3/internal/infra/api/web/validation"
	"github.com/diogocardoso/go_lab3/internal/usecase/auction_usecase"
	"github.com/gin-gonic/gin"
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

type AuctionController struct {
	auctionUseCase auction_usecase.AuctionUseCaseInterface
}

func NewAuctionController(auctionUseCase auction_usecase.AuctionUseCaseInterface) *AuctionController {
	return &AuctionController{
		auctionUseCase: auctionUseCase,
	}
}

func (u *AuctionController) CreateAuction(c *gin.Context) {
	var auctionInputDTO auction_usecase.AuctionInputDTO

	if err := c.ShouldBindJSON(&auctionInputDTO); err != nil {
		log.Println("Erro ao ler body da requisição", err)
		restErr := validation.ValidateErr(err)

		c.JSON(restErr.Code, restErr)
		return
	}

	err := u.auctionUseCase.CreateAuction(context.Background(), auctionInputDTO)
	if err != nil {
		restErr := rest_err.ConvertError(err)
		logger.Error("Erro ao criar leilão", err)
		c.JSON(restErr.Code, restErr)
		return
	}

	c.Status(http.StatusCreated)
}
