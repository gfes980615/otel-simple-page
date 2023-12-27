package lib

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
)

var (
	HTTPMethodList     = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	HTTPStatusCodeList = []string{"200", "203", "400", "403", "404", "500", "502", "504"}
	CpidList           = []string{"iam", "fbt", "cam", "ent", "spa", "aak", "ale", "apc", "wbc", "vws", "tmv", "cus", "svp", "uic", "amc", "asp", "mco", "obp", "iem", "asm", "agc", "von", "ops", "anc", "aal", "sps", "acp", "arc", "act", "adg", "scm", "apm", "asr", "sss", "rex", "sfs", "ssp", "sae", "ars", "adm", "awb", "ath", "rca", "ati", "atl", "asc", "ase", "arp", "ara", "clr", "hdl", "eiq", "ddl", "sdl", "dla", "dld", "atc", "psd", "aei", "ves", "aep", "eai", "xms", "ani", "xns", "ams", "ami", "asi", "tpc", "tpi", "seb", "xbc", "cms", "rms", "asg", "sgi", "ssb", "xes", "srp", "sea", "sjg", "irs", "mab", "tig", "sos", "xcs", "vcs", "crp", "acl", "dlt", "ass", "azs", "szn", "art", "sph", "sig", "adi", "dva", "aed", "aod", "sza", "aps", "wae", "air", "asd", "sno", "cta", "sdf", "can", "sao", "pao", "sca", "scf", "sds", "scr", "sim", "pds", "pdd", "sna", "pna", "pei", "pdi", "sem", "stp", "ptp", "scc", "esc", "sas", "sct", "scs", "sof", "csg", "pos", "sws", "swf", "ptn", "stn", "pts", "sts", "ppo", "spo", "pan", "mam", "eps", "xsm", "pva", "riv", "siv", "pns", "aad", "opa", "okt", "olp", "goo", "rsa", "sta", "sdw", "tag", "msp", "sre", "msc", "acv", "tas", "mkt", "ade", "sis", "ers", "cti", "bif", "sma", "piv", "sbc", "tms", "fbs", "kef", "ald", "epp", "axf", "ddr", "hkc", "stu", "snb", "saa", "prm", "eoa", "stg", "idp"}

	highCpidList = CpidList[0:70]
	lowCpidList  = CpidList[70:]

	vipIDs = getCustomerID("vip_ids", 20)
	ids    = getCustomerID("normal_ids", 6980)
)

type TestRawData struct {
	Timestamp      time.Time
	ServiceName    string
	CustomerID     string
	Cpid           string
	HTTPMethod     string
	HTTPStatusCode string
	API            string
	ResponseTime   int
}

func GenRawData(t time.Time) TestRawData {
	return TestRawData{
		Timestamp:      t,
		ServiceName:    "local-test",
		CustomerID:     randomCustomerID(),
		Cpid:           randomCpid(),
		HTTPMethod:     randomHTTPMethod(),
		HTTPStatusCode: randomHTTPStatusCode(),
		API:            randomAPI(randomCpid()),
		ResponseTime:   randomResponseTime(),
	}
}

func randomCustomerID() string {
	if time.Now().UnixNano()%11 >= 2 {
		return vipIDs[rand.Intn(len(vipIDs))]
	}
	return ids[rand.Intn(len(ids))]
}

func randomCpid() string {
	if time.Now().UnixNano()%11 >= 2 {
		return highCpidList[rand.Intn(len(highCpidList))]
	}
	return CpidList[rand.Intn(len(lowCpidList))]
}

func randomHTTPMethod() string {
	return HTTPMethodList[rand.Intn(len(HTTPMethodList))]
}

func randomHTTPStatusCode() string {
	code_2xx := []string{"200", "203", "204", "205", "206", "207", "208", "226"}
	code_4xx := []string{"400", "403", "404", "405"}
	code_5xx := []string{"500", "502", "504", "505", "599"}
	if time.Now().UnixNano()%11 >= 2 {
		return code_2xx[rand.Intn(len(code_2xx))]
	}
	if time.Now().UnixNano()%11 == 1 {
		return code_4xx[rand.Intn(len(code_4xx))]
	}
	if time.Now().UnixNano()%11 == 0 {
		return code_5xx[rand.Intn(len(code_5xx))]
	}
	return HTTPStatusCodeList[rand.Intn(len(HTTPStatusCodeList))]
}

func randomAPI(cpid string) string {
	prefix := []string{"internal", "public", "ui", "external"}
	version := []string{"v1", "v2", "v3"}
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz")
	api := fmt.Sprintf("/%s/%s/%s/%s",
		prefix[rand.Intn(len(prefix))],
		cpid,
		version[rand.Intn(len(version))],
		string(letterRunes[rand.Intn(len(letterRunes))]),
	)
	return api
}

func randomResponseTime() int {
	if time.Now().UnixNano()%11 >= 2 {
		return 10 + rand.Intn(1000)
	}

	return 1000 + rand.Intn(10000)
}

func getCustomerID(idFile string, num int) []string {
	if _, err := os.Stat(idFile); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Create new customer id file: ", idFile)
		f, err := os.Create(idFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		ids := genCustomerIDV2(num)
		for _, id := range ids {
			f.WriteString(id + "\n")
		}
		return ids
	}
	fmt.Println("Read from file: ", idFile)
	ids := []string{}
	f, err := os.Open(idFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for {
		id := ""
		_, err := fmt.Fscanf(f, "%s\n", &id)
		if err != nil {
			break
		}
		ids = append(ids, id)
	}
	return ids
}

func genCustomerIDV2(num int) []string {
	ids := []string{}
	for i := 0; i < num; i++ {
		ids = append(ids, uuid.New().String())
	}
	return ids
}
