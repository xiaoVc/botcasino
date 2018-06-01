package main

import (
	"fmt"
	"strconv"
	"syscall"

	"github.com/vrecan/death"
	"github.com/zhangpanyi/basebot/logger"
	"github.com/zhangpanyi/basebot/telegram/updater"
	"github.com/zhangpanyi/botcasino/config"
	"github.com/zhangpanyi/botcasino/logic"
	context "github.com/zhangpanyi/botcasino/logic/context"
	"github.com/zhangpanyi/botcasino/logic/notice"
	"github.com/zhangpanyi/botcasino/logic/syncfee"
	"github.com/zhangpanyi/botcasino/logic/timer"
	"github.com/zhangpanyi/botcasino/models"
	"github.com/zhangpanyi/botcasino/pusher"
	"github.com/zhangpanyi/botcasino/remote"
	"github.com/zhangpanyi/botcasino/service"
	"github.com/zhangpanyi/botcasino/storage"
	"github.com/zhangpanyi/botcasino/webrpc"
	"github.com/zhangpanyi/botcasino/withdraw"
	"upper.io/db.v3/sqlite"
)

func main() {
	// 加载配置文件
	config.LoadConfig("master.yml")

	// 初始化日志库
	serveCfg := config.GetServe()
	logger.CreateLoggerOnce(logger.DebugLevel, logger.InfoLevel)

	// 连接到数据库
	err := storage.Connect(serveCfg.BolTDBPath)
	if err != nil {
		logger.Panic(err)
	}
	dbcfg := serveCfg.SQlite
	settings := sqlite.ConnectionURL{
		Database: dbcfg.Database,
		Options:  dbcfg.Options,
	}
	err = models.Connect(settings)
	if err != nil {
		logger.Panic(err)
	}

	orderID, err := models.InsertWithdraw(1000, "byte01", "btc", 100, 1)
	fmt.Println(orderID, err)
	return

	// 创建更新器
	botUpdater, err := updater.NewUpdater(serveCfg.Port, serveCfg.Domain, serveCfg.APIWebsite)
	if err != nil {
		logger.Panic(err)
	}
	webrpc.InitRoute(botUpdater.GetRouter())

	// 连接钱包服务
	remote.NewWalletServerForOnce(serveCfg.WalletService.Address,
		serveCfg.WalletService.Port)

	// 同步转账手续费
	syncfee.CheckFeeStatusAsync()

	// 运行转账服务
	withdraw.RunWithdrawServiceForOnce(6)

	// 启动红包机器人
	context.CreateManagerForOnce(serveCfg.BucketNum)
	bot, err := botUpdater.AddHandler(serveCfg.Token, logic.NewUpdate)
	if err != nil {
		logger.Panic(err)
	}
	logger.Infof("Lucky money bot id is: %d", bot.ID)
	pool := updater.NewPool(2048)
	timer.StartTimerForOnce(bot, pool)

	// 初始化推送配置
	notice.InitBotForOnce(bot)

	// 创建消息推送器
	pusher.CreatePusherForOnce(pool)

	// 启动RPC服务
	address := serveCfg.GRPCBindAddress + ":" + strconv.Itoa(serveCfg.GRPCPort)
	go service.RunService(address)

	// 启动更新服务器
	logger.Infof("Casino server started, grpc listen: %s", address)
	go func() {
		err = botUpdater.ListenAndServe()
		if err != nil {
			logger.Panicf("Casino server failed to listen: %v", err)
		}
	}()

	// 捕捉退出信号
	d := death.NewDeath(syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL,
		syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGALRM)
	d.WaitForDeathWithFunc(func() {
		storage.Close()
		logger.Info("Casino server stoped")
	})
}
