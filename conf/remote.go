package conf

import (
	"sync"
	"github.com/astaxie/beego"
	"net/http"
	"io/ioutil"
	"strings"
)

type remoteConfig struct {
	sync.RWMutex
	data map[string]string
}

var _remoteConf remoteConfig

func Remote() remoteConfig {

	if len(_remoteConf.data) == 0 {
		_remoteConf.data = make(map[string]string)
	}
	return _remoteConf
}

func (c remoteConfig) Get(key string) string {
	c.RLock()
	val := c.data[key]
	c.RUnlock()
	return val
}

func (c remoteConfig) Load(authCode string) {

	baseApi := beego.AppConfig.String("config_api")

	keysStr, err := httpGet(baseApi + "/keys")
	if err != nil {
		panic("Load Remote Configuration failed.")
	}

	keysStr = strings.Replace(keysStr, "\"", "", -1)
	keys := strings.Split(keysStr, ",")

	_remoteConf.Lock()
	for _, key := range keys {
		val, _ := httpGet(baseApi + "/get/" + key + "?auth=" + authCode)
		//beego.Debug("Load Remote Config: ", key, ":", val)
		_remoteConf.data[key] = val

	}
	_remoteConf.Unlock()
}

func httpGet(uri string) (html string, err error) {

	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
