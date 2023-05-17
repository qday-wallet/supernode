package push

import (
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/store/config"
	"github.com/uduncloud/easynode/store/service"
	"log"
	"testing"
	"time"
)

func Init() service.DbMonitorAddressInterface {
	cfg := config.LoadConfig("./../../../cmd/store/config.json")
	return NewChService(&cfg, xlog.NewXLogger())
}

func TestClickhouseDb_AddMonitorAddress(t *testing.T) {
	s := Init()
	log.Println(s.AddMonitorAddress(200, &service.MonitorAddress{Address: "0xe5cB067E90D5Cd1F8052B83562Ae670bA4A211a8", Token: "5fe5f231-7051-4caf-9b52-108db92edbb4", BlockChain: 200, TxType: "0x2", Id: time.Now().UnixMilli()}))
}

func TestClickhouseDb_GetAddressByToken(t *testing.T) {
	s := Init()
	log.Println(s.GetAddressByToken(200, "5fe5f231-7051-4caf-9b52-108db92edbb4"))
}

func TestClickhouseDb_NewToken(t *testing.T) {
	s := Init()
	log.Println(s.NewToken(&service.NodeToken{
		Token: "1231token",
		Email: "123@qq.com",
		Id:    time.Now().UnixMilli(),
	}))
}

func TestClickhouseDb_UpdateToken(t *testing.T) {

	s := Init()
	log.Println(s.UpdateToken("1231token", &service.NodeToken{
		Token: "1231token",
		Email: "125@qq.com",
		Id:    time.Now().UnixMilli(),
	}))
}

func TestClickhouseDb_GetNodeTokenByEmail(t *testing.T) {
	s := Init()
	log.Println(s.GetNodeTokenByEmail("125@qq.com"))
}
