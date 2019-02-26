package logic

import (
	"encoding/json"
	"github.com/Jeffail/gabs"
	"github.com/gosuri/uiprogress/util/strutil"
	"github.com/prometheus/common/log"

	"gopkg.in/cheggaaa/pb.v1"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const sizeBarText = 60

// An Executor is structure
type Executor struct {
	threads  int
	baseURL  string
	username string
	token    string
}

type channelmsg struct {
	job    string
	err    error
	status int
}

var jmutex = &sync.Mutex{}

// NewExecutor returns a new logic that helps to run jenkins jobs.
func NewExecutor(threads int, baseURL, username, token string) Executor {
	return Executor{threads: threads, baseURL: baseURL, username: username, token: token}
}

// Execute runs jenkins jobs asynchronous
func (e *Executor) Execute(orders ...Order) {
	log.Debugln("execute ", e, "    ", orders)
	var cmsg channelmsg
	ch := make(chan channelmsg)
	//fmt.Println("Transfering....")
	wg := new(sync.WaitGroup)
	batch := len(orders)/e.threads + 1
	i2 := batch
	pbList := make([]*pb.ProgressBar, e.threads)
	i := 0
	for i = 0; i < e.threads; i++ {

		pb := newProgressBar(wg)
		pbList[i] = pb
		log.Debug("new worker  %n", i)
		go e.executeJOBAsync(&orders[i], pb, ch, wg)
		if i2 == len(orders) {
			break
		}
		i2 = i2 + batch
		if i2 > len(orders)-1 {
			i2 = len(orders)
		}
	}
	log.Debugln("pbList[0:i] = ", len(pbList[0:i+1]), " i=", i)
	poolTransfer, err := pb.StartPool(pbList[0 : i+1]...)
	if err != nil {
		panic(err)
	}

	defer poolTransfer.Stop()

	//time.Sleep(time.Second*1000)
	for t1 := 0; t1 <= i; t1++ {
		cmsg = <-ch
		if cmsg.status != 0 {
			log.Errorf("failed to build job %s, details: %s", cmsg.job, cmsg.err)
			log.Fatalf("Failed to transfer to FRONT region...")
		}
	}
}

func (e *Executor) executeJOBAsync(order *Order, pb *pb.ProgressBar, ch chan channelmsg, wg *sync.WaitGroup) {
	defer wg.Done()
	defer pb.Finish()
	_, err := e.executeJOB(order, pb)
	if err == nil {
		ch <- channelmsg{status: 0, job: order.Job}
	} else {
		ch <- channelmsg{err: err, status: 1, job: order.Job}
	}
	return
}

func (e *Executor) executeJOB(order *Order, pb *pb.ProgressBar) (string, error) {
	code, lastSuccessfulJobInfo, err := e.jreq(order.Job, "/lastSuccessfulBuild/api/json", map[string]interface{}{})
	if err != nil {
		return "", err
	}

	durationSt := lastSuccessfulJobInfo.Path("duration").Data()
	duration := 100
	_, ok := durationSt.(float64)
	if ok {
		duration = int(durationSt.(float64))
	}

	jmutex.Lock()
	code, lastJobInfo, err := e.jreq(order.Job, "/api/json", map[string]interface{}{})
	if err != nil {
		return "", err
	}
	log.Debug("start1 deploy ")
	code, _, err2 := e.jreq(order.Job, "/build", order.Query)
	log.Debug("start2 deploy " + strconv.Itoa(code) + "\t" + "path:=" + order.Job)
	if err != nil {
		return "", err
	}
	//time.Sleep(time.Second*5)
	jmutex.Unlock()

	if err2 != nil {
		return "", err2
	}
	if code != 201 {
		log.Fatalf("http error code2 " + strconv.Itoa(code) + "\t" + "path:=" + order.Job)
	}
	nextBuildNumber := lastJobInfo.Path("nextBuildNumber").String()
	log.Debugf("nextBuildNumber %s", nextBuildNumber)
	//fmt.Print("build...")
	status := "building"

	ticsize := duration / 600
	pb.SetTotal(ticsize)
	code = 1
	t := 0
	for t > -1 {
		//if true==true{ for gif building
		if t%5 == 0 && t > 1 {
			code, lastBuildInfo, err := e.jreq(order.Job, "/"+string(nextBuildNumber)+"/api/json", map[string]interface{}{})
			if err != nil {
				return "", err
			}
			if t == ticsize {

			}
			code = 200
			if code == 200 {

				if lastBuildInfo.Path("building").Data().(bool) == false {
					//if t == ticsize  { for gif building
					if lastBuildInfo.Path("result").Data().(string) == "SUCCESS" {
						//if true == true{ for gif building
						status = "finished"
						pb.Prefix(strutil.PadRight(order.Job+"("+order.getVals()+")"+": "+status, sizeBarText, []byte(" ")[0]))
						for k := t; k < ticsize; k++ {
							if pb != nil {
								pb.Prefix(strutil.PadRight(order.Job+"("+order.getVals()+")"+": "+status, sizeBarText, []byte(" ")[0]))
								pb.Increment()
							}
							time.Sleep(10 * time.Millisecond)
						}
						time.Sleep(1 * time.Second)
						return "", err
					}

					if pb != nil {
						pb.Prefix(strutil.PadRight(order.Job+"("+order.getVals()+")"+": failed", sizeBarText, []byte(" ")[0]))
						pb.Increment()
					}
					time.Sleep(1 * time.Second)
					id := lastBuildInfo.Path("id").Data().(string)
					path := "/job/" + order.Job + "/" + id + "/console"
					//path := "/job/" + order.Job + "/" + "10" + "/console" for gif building
					log.Fatalf("FAILURE, link " + order.Job + path)

				}
			} else {
				if t > 10 {
					if pb != nil {
						pb.Prefix(strutil.PadRight(order.Job+"("+order.getVals()+")"+": failed", sizeBarText, []byte(" ")[0]))
						pb.Increment()
					}
					log.Fatalf("response error " + strconv.Itoa(code) + "\t" + "path:=" + order.Job)
				}

			}
			if t > 30000 {
				if pb != nil {
					pb.Prefix(strutil.PadRight(order.Job+"("+order.getVals()+")"+": failed", sizeBarText, []byte(" ")[0]))
					pb.Increment()
				}
				log.Fatalf("response error1 " + strconv.Itoa(code) + "\t" + "path:=" + order.Job)
			}
		}

		if pb != nil {
			pb.Increment()
			pb.Prefix(strutil.PadRight(order.Job+"("+order.getVals()+")"+": "+status, sizeBarText, []byte(" ")[0]))
		}
		//time.Sleep(50 * time.Millisecond);for gif building
		time.Sleep(time.Second)
		t++
	}
	log.Fatalf("timeout of deployment process")
	return "", nil
}

func (e *Executor) jreq(job string, api string, query map[string]interface{}) (int, *gabs.Container, error) {
	fn := func(request *http.Request) {
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		request.SetBasicAuth(e.username, e.token)
	}
	fullURL := e.baseURL + "/job/" + job + api
	dataArr := []map[string]string{}
	for key, val := range query {
		dataArr = append(dataArr, map[string]string{"name": key, "value": val.(string)})
	}
	queryJSON := map[string][]map[string]string{
		"parameter": dataArr,
	}
	queryBytes, _ := json.Marshal(queryJSON)
	urlquery := url.Values{}
	urlquery.Add("json", string(queryBytes))
	return req("POST", fullURL, []byte(urlquery.Encode()), fn)
}

func newProgressBar(wg *sync.WaitGroup) *pb.ProgressBar {
	pb := pb.New(3).Prefix(strutil.PadRight("", sizeBarText, []byte(" ")[0]))
	pb.SetWidth(120)
	wg.Add(1)
	return pb
}
