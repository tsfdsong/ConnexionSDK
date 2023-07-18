package web

import (
	"github/Connector-Gamefi/ConnectorGoSDK/config"
	"github/Connector-Gamefi/ConnectorGoSDK/core/common/const_def"
	"github/Connector-Gamefi/ConnectorGoSDK/distmiddleware/skywalking"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/middleware"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/middleware/log_mid"
	"github/Connector-Gamefi/ConnectorGoSDK/web/common"
	"github/Connector-Gamefi/ConnectorGoSDK/web/dashboard/dashcontroller"
	"github/Connector-Gamefi/ConnectorGoSDK/web/marketplace"
	"github/Connector-Gamefi/ConnectorGoSDK/web/rpc"
	"github/Connector-Gamefi/ConnectorGoSDK/web/sdkbackend/sdkcontroller"
	"net/http"
	"os"
	"strings"

	v3 "github.com/SkyAPM/go2sky-plugins/gin/v3"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/unrolled/secure"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		//java code
		// c.Writer.Header().Set("Access-Control-Allow-Origin", c.Writer.Header().Get("Origin"))
		// c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PUT,DELETE")
		// c.Writer.Header().Set("Access-Control-Allow-Headers", c.Writer.Header().Get("Access-Control-Allow-Headers"))
		c.Next()
	}
}

// JWTAuthMiddleware
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": common.AuthCheck,
				"msg":  "the auth of header is empty",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": common.AuthCheck,
				"msg":  "the format of auth is wrong",
			})
			c.Abort()
			return
		}

		mc, err := common.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": common.AuthCheck,
				"msg":  "invalided Token",
			})
			c.Abort()
			return
		}

		c.Set("address", mc.Address)
		c.Next()
	}
}

func TLSHandler(port string) gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     ":" + port,
		})

		err := secureMiddleware.Process(c.Writer, c.Request)
		if err != nil {
			return
		}

		c.Next()
	}
}

func SDKServerRoute() *gin.Engine {
	router := gin.New()

	middlewareLogConfig := config.GetMiddlewareLogConfig()
	recoverFile, err := os.OpenFile(middlewareLogConfig.RecoverLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil || recoverFile == nil {

		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("open recover log file failed")
		}
		if recoverFile == nil {
			logger.Logrus.Error("open recover log file failed:recoverFile is nil")
		}

		return nil
	}

	router.Use(log_mid.Logger(middlewareLogConfig.VisitLogFile, middlewareLogConfig.SkipPath...), gin.RecoveryWithWriter(recoverFile))
	router.Use(CORSMiddleware())

	//TLS config
	// router.Use(TLSHandler(config.GetServerConfig().HttpPort))

	router.GET("/dashboard/auth", GetAuth)

	//skywalking middle
	router.Use(v3.Middleware(router, skywalking.GetSkyTrace()))

	backend := router.Group("/backend")
	{
		//backend
		backend.GET("/user/getAssets", sdkcontroller.GetUserAssets)
		backend.GET("/user/getNFTAssets", sdkcontroller.GetUserNFTAssets)
		backend.POST("/withdraw/setExamine", sdkcontroller.WithdrawExamineSet)
		backend.POST("/deposit/repairOrder", sdkcontroller.DepositRepairOrder)
		backend.POST("/withdraw/repairOrder", sdkcontroller.WithdrawRepairOrder)

		backend.POST("/parseSpecifiedBlockLog", sdkcontroller.ParseSpecifiedBlockLog)
		backend.POST("/parseLogSwitch", sdkcontroller.ParseLogSwitch)
	}

	v1 := router.Group("/v1")
	{
		// v1.Use(JWTAuthMiddleware())
		v1.GET("/dashboard/gameFTAssets", dashcontroller.QueryGameFTAssets)
		v1.GET("/dashboard/gameNFTAssets", dashcontroller.QueryGameNFTAssets)
		v1.POST("/dashboard/gameAssetDetail", dashcontroller.QueryGameAssetDetail)
		v1.POST("/dashboard/chainFTAssets", dashcontroller.ChainFTAssets)
		v1.GET("/nftAssets", dashcontroller.QueryEquipment)

		//ft deposit and withdraw
		v1.POST("/dashboard/ft/deposit", middleware.ReSubmitMiddleware(const_def.FT_DEPOSIT, "address"), dashcontroller.DepositGameERC20Token)
		v1.POST("/dashboard/ft/prewithdraw", middleware.ReSubmitMiddleware(const_def.FT_WITHDRAW, "address"), dashcontroller.WithdrawGameERC20Token)
		v1.POST("/dashboard/ft/withdraw", dashcontroller.ClaimGameERC20Token)

		//nft deposit and withdraw
		v1.POST("/dashboard/nft/prewithdraw", dashcontroller.NFTPreWithdrawController)
		v1.GET("/dashboard/nft/claim", dashcontroller.NFTClaimController)
		v1.POST("/dashboard/nft/deposit", dashcontroller.NFTDepositController)
	}

	marketplaceRouter := router.Group("/marketplace")
	{
		marketplaceRouter.GET("/orderList", marketplace.MarketOrderList)
		marketplaceRouter.GET("/orderDetail", marketplace.OrderDetail)
		marketplaceRouter.GET("/profile/activity", marketplace.Activity)
	}

	rpcRouter := router.Group("/rpc")
	{
		rpcRouter.POST("/ft/newDeposit", rpc.NewFtDeposit)
		rpcRouter.POST("/ft/newWithdraw", rpc.NewFtWithdraw)
		rpcRouter.POST("/ft/confirmedDeposit", rpc.ConfirmedFtDeposit)
		rpcRouter.POST("/ft/confirmedWithdraw", rpc.ConfirmedFtWithdraw)

		rpcRouter.POST("/nft/newMintWithdraw", rpc.NewNftMintWithdraw)
		rpcRouter.POST("/nft/newUpdateWithdraw", rpc.NewNftUpdateWithdraw)
		rpcRouter.POST("/nft/newDeposit", rpc.NewNftDeposit)
		rpcRouter.POST("/nft/confirmedWithdraw", rpc.ConfirmedNftWithdraw)
		rpcRouter.POST("/nft/confirmedDeposit", rpc.ConfirmedNftDeposit)
	}

	return router
}
