package main

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/adapter/orm"
	"github.com/beego/beego/v2/core/config/yaml"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"v1-go-api/graph/models"
	_ "v1-go-api/routers"
	"v1-go-api/task"
)

var (
	logger = logs.GetBeeLogger()
)

func init() {
	dbIp, err := beego.AppConfig.String("dbIp")
	if err != nil {
		fmt.Println(fmt.Sprintf("get config dbIp err: %s", err.Error()))
		panic(err)
	}
	dbPort, err := beego.AppConfig.String("dbPort")
	if err != nil {
		fmt.Println(fmt.Sprintf("get config dbPort err: %s", err.Error()))
		panic(err)
	}
	dbUser, err := beego.AppConfig.String("dbUser")
	if err != nil {
		fmt.Println(fmt.Sprintf("get config dbUser err: %s", err.Error()))
		panic(err)
	}
	dbPass, err := beego.AppConfig.String("dbPass")
	if err != nil {
		fmt.Println(fmt.Sprintf("get config dbPass err: %s", err.Error()))
		panic(err)
	}
	dbName, err := beego.AppConfig.String("dbName")
	if err != nil {
		fmt.Println(fmt.Sprintf("get config dbName err: %s", err.Error()))
		panic(err)
	}

	// ; sqlconn = swap_user:xr6M54me06hG6YL6@tcp(127.0.0.1:3306)/youswap
	sqlConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?loc=Local", dbUser, dbPass, dbIp, dbPort, dbName)

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", sqlConn)
	fmt.Println(fmt.Sprintf("regist database success"))
}

func main() {

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))

	logger.EnableFuncCallDepth(true)
	logger.Async()

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	err := beego.LoadAppConfig("yaml", "conf/app.yaml")
	conf, err := yaml.ReadYmlReader("conf/app.yaml")
	if err != nil {
		logger.Error("loadConfig err: {}", err)
		panic(err)
	}
	graphNode := conf["GraphNode"].(map[string]interface{})

	for k, v := range graphNode {
		bytes, err := json.Marshal(v)
		if err != nil {
			logger.Error("marshal err {}", err)
		}
		graphData := &models.GraphData{}
		err = json.Unmarshal(bytes, &graphData)
		if err != nil {
			//logger.Error("unmarshal err {}", err)
		}
		if graphData.Igrone {
			continue
		}
		if len(graphData.Name) == 0 {
			graphData.Name = k
		}

		if graphData.BundleId == "" {
			graphData.BundleId = "1"
		}
		task.GraphDataList[k] = graphData
		logger.Info(k)
	}

	go task.Start()
	//task.BuildAirdrop()
	//
	logger.Info("load config success ")
	beego.Run()

	time.Sleep(time.Hour)
}
