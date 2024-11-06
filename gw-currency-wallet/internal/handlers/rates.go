package handlers

import (
    pb "github.com/zxckontolxd/proto-exchange/exchange"
    log "github.com/sirupsen/logrus"

    "google.golang.org/grpc"
    "net/http"
	"github.com/gin-gonic/gin"
)

func Rates(ctx *gin.Context) {
    // не знаю как лучше, устанавливать ли постоянно соединение или установить одно соединение на всю программу
    // пока оставлю так
	exchangerConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Errorf("Cannot dial with grpc service (gw-exchanger): %v", err)
        //TODO json
        return
    }
    defer exchangerConn.Close()

    exchangerClient := pb.NewExchangeServiceClient(exchangerConn)

    response, err := exchangerClient.GetExchangeRates(ctx, &pb.Empty{})
    if err != nil {
        log.Errorf("Cannot get rates: %v", err)
        //TODO json
        return
    }
    //map<string, float> rates = 1; // ключ: валюта, значение: курс
    ctx.JSON(http.StatusOK, gin.H{
        "rates": gin.H{
            "USD": response.Rates["USD"],
            "RUB": response.Rates["RUB"],
            "EUR": response.Rates["EUR"],
        },
    })
}