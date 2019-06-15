package getter

import (	
	"github.com/go-clog/clog"
	"github.com/Aiicy/htmlquery"
	"github.com/henson/proxypool/pkg/models"
	"strings"
	"strconv"
	"time"
	"sync"
)

func StartKDL(ipChan chan<- *models.IP, wg *sync.WaitGroup){
	for i:=1; i<=2901; i++{
		temp := kDL(i)
		//log.Println("[run] get into loop")
		for _, v := range temp {
			//log.Println("[run] len of ipChan %v",v)
			ipChan <- v
		}
		time.Sleep(2 * time.Second)
	}
	wg.Done()
}
var count int = 0
// KDL get ip from kuaidaili.com
func kDL(i int) (result []*models.IP) {
	clog.Info("[kuaidaili] start")
	pollURL := "http://www.kuaidaili.com/free/inha/"
	pollURL = pollURL + strconv.Itoa(i) + "/"
	doc,_ := htmlquery.LoadURL(pollURL)
	trNode, err := htmlquery.Find(doc, "//tbody//tr")
	if err != nil {
		clog.Warn("KDL find error: %v", err)
	}
	for i := 0; i < len(trNode); i++ {
		IP := models.NewIP()
		tdNode, _ := htmlquery.Find(trNode[i],"//td")
		IP.Data = htmlquery.InnerText(tdNode[0])
		IP.Port = htmlquery.InnerText(tdNode[1])
		IP.Type = htmlquery.InnerText(tdNode[2])
		IP.Protocol = strings.ToLower(htmlquery.InnerText(tdNode[3]))
		IP.Position = htmlquery.InnerText(tdNode[4])
		IP.Speed = extractSpeed(htmlquery.InnerText(tdNode[5]))
		result = append(result, IP)
		count++
	}

	clog.Info("[kuaidaili] done, %d", count)
	return
}
