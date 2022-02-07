package sharksdk

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/shark/minigame-common/conf"
)

type BaseResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Time int             `json:"time"`
	Data json.RawMessage `json:"data"`
}

type GoldQueryResponse struct {
	Gold int64 `json:"gold"`
}

type GoldUseResponse struct {
	Gold int64 `json:"gold"`
}

func signUrl(values url.Values) string {
	valueArr := make([]string, 0)
	for _, vs := range values {
		for _, v := range vs {
			valueArr = append(valueArr, v)
		}
	}
	valueArr = append(valueArr, conf.Ini.SharkSdk.Secret)
	sort.Strings(valueArr)
	args := strings.Join(valueArr, "")
	sign := fmt.Sprintf("%x", md5.Sum([]byte(args)))
	sign = strings.ToUpper(sign)
	return sign
}

func post(method string, values url.Values) (*BaseResponse, error) {
	url := fmt.Sprintf("%s%s", conf.Ini.SharkSdk.Url, method)
	//log.Println(url)
	now := time.Now().Unix()
	values.Add("time", fmt.Sprintf("%d", now))
	sign := signUrl(values)
	values.Add("sign", sign)
	//log.Printf("aaa %+v\n", values)
	res, err := http.PostForm(url, values)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	//log.Printf("[sharksdk] post url=%s body=%s\n", url, string(body))
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	baseResponse := &BaseResponse{}
	json.Unmarshal(body, baseResponse)
	return baseResponse, nil
}

func GoldQuery(token string, openid string) (*GoldQueryResponse, error) {
	values := url.Values{}
	values.Add("token", token)
	values.Add("user_id", openid)
	baseResponse, err := post("/h5game/get_user_gold", values)
	if err != nil {
		return nil, err
	}
	if baseResponse.Code != 1 {
		return nil, errors.New(baseResponse.Msg)
	}
	result := &GoldQueryResponse{}
	if err := json.Unmarshal(baseResponse.Data, result); err != nil {
		return nil, err
	}
	return result, nil
}

func GoldUse(token string, openid string, amount int64) (*GoldUseResponse, error) {
	gameId := conf.Ini.Game.Id
	values := url.Values{}
	values.Add("game_id", fmt.Sprintf("%d", gameId))
	values.Add("user_id", openid)
	values.Add("token", token)
	values.Add("amount", fmt.Sprintf("%d", amount))
	baseResponse, err := post("/h5game/use_gold", values)
	if err != nil {
		return nil, err
	}
	if baseResponse.Code != 1 {
		return nil, errors.New(baseResponse.Msg)
	}
	result := &GoldUseResponse{}
	if err := json.Unmarshal(baseResponse.Data, result); err != nil {
		return nil, err
	}
	return result, nil
}
