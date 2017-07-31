package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/astaxie/beego/orm"
	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	"github.com/googollee/go-socket.io"
	"github.com/satori/go.uuid"
	"github.com/toolkits/file"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// APIConfig ...
type APIConfig struct {
	CreateProxy         string `json:"createProxy"`
	DecodeJWT           string `json:"decodeJWT"`
	EncodeJWT           string `json:"encodeJWT"`
	GetBalance          string `json:"getBalance"`
	GetCoin             string `json:"getCoin"`
	GetIPFS             string `json:"getIPFS"`
	PrepareIPFS         string `json:"prepareIPFS"`
	PrepareProxy        string `json:"prepareProxy"`
	Send                string `json:"send"`
	SignTransaction     string `json:"signTransaction"`
	TransferProxyChange string `json:"transferProxyChange"`
	TransferProxySign   string `json:"transferProxySign"`
	VerifyJWT           string `json:"verifyJWT"`
}

// DBConfig ...
type DBConfig struct {
	Address string `json:"address"`
	Idle    int    `json:"idle"`
	Max     int    `json:"max"`
}

// JWTConfig ...
type JWTConfig struct {
	ServerPrivateKey string `json:"serverPrivateKey"`
	ServerPublicKey  string `json:"serverPublicKey"`
	ServerURL        string `json:"serverURL"`
}

// PathConfig ...
type PathConfig struct {
	Claim string `json:"claim"`
	Login string `json:"login"`
}

// GlobalConfig ...
type GlobalConfig struct {
	API      *APIConfig  `json:"api"`
	Database *DBConfig   `json:"database"`
	Delegate string      `json:"delegate"`
	JWT      *JWTConfig  `json:"jwt"`
	Path     *PathConfig `json:"path"`
	Port     int         `json:"port"`
}

var (
	configFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

// Config ...
func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func setConfig(newConfig *GlobalConfig) {
	configLock.RLock()
	defer configLock.RUnlock()
	config = newConfig
}

func parseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("config file not specified: use -c $filename")
	}
	if !file.IsExist(cfg) {
		log.Fatalln("config file specified not found:", cfg)
	}
	configFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file", cfg, "error:", err.Error())
	}
	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file", cfg, "error:", err.Error())
	}
	setConfig(&c)
}

func setResponse(rw http.ResponseWriter, resp map[string]interface{}) {
	if _, ok := resp["auth"]; ok {
		delete(resp, "auth")
	}
	if _, ok := resp["method"]; ok {
		delete(resp, "method")
	}
	if _, ok := resp["params"]; ok {
		delete(resp, "params")
	}
	if _, ok := resp["result"]; ok {
		result := resp["result"].(map[string]interface{})
		if val, ok := result["error"]; ok {
			errors := val.([]string)
			if len(errors) > 0 {
				delete(resp, "result")
				resp["error"] = errors
			} else {
				delete(resp["result"].(map[string]interface{}), "error")
				if val, ok = result["items"]; ok {
					resp["result"] = val
				}
				if val, ok = result["count"]; ok {
					resp["count"] = val
				}
				if val, ok = result["anomalies"]; ok {
					resp["anomalies"] = val
				}
			}
		}
	}
	resp["time"] = getNow()
	renderJSON(rw, resp)
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

func getNow() string {
	t := time.Now()
	now := t.Format("2006-01-02 15:04:05")
	return now
}

func setError(error string, result map[string]interface{}) {
	log.Println("Error =", error)
	result["error"] = append(result["error"].([]string), error)
}

func postByJSON(req *http.Request, destination string, params map[string]interface{}, result map[string]interface{}) map[string]interface{} {
	log.Println("func postByJSON() destination =", destination)
	log.Println("func postByJSON() params =", params)
	JSONString, _ := json.Marshal(params)
	reqPost, err := http.NewRequest("POST", destination, bytes.NewBuffer(JSONString))
	if err != nil {
		setError(err.Error(), result)
	}
	reqPost.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqPost)
	if err != nil {
		setError(err.Error(), result)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	json, err := simplejson.NewJson(buf.Bytes())
	if err != nil {
		setError(err.Error(), result)
	}
	response := map[string]interface{}{}
	response, _ = json.Map()
	return response
}

func postByForm(req *http.Request, destination string, params map[string]string, result map[string]interface{}) map[string]interface{} {
	log.Println("func postByForm() destination =", destination)
	log.Println("func postByForm() params =", params)
	form := url.Values{}
	for key, value := range params {
		form.Add(key, value)
	}
	client := &http.Client{}
	resp, err := client.PostForm(destination, form)
	if err != nil {
		setError(err.Error(), result)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	json, err := simplejson.NewJson(buf.Bytes())
	if err != nil {
		setError(err.Error(), result)
	}
	response := map[string]interface{}{}
	response, _ = json.Map()
	return response
}

func createUser(rw http.ResponseWriter, r *http.Request) {
	errors := []string{}
	result := map[string]interface{}{}
	result["error"] = errors

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	sjson, err := simplejson.NewJson(buf.Bytes())
	if err != nil {
		setError(err.Error(), result)
	}
	user := map[string]interface{}{}
	user, _ = sjson.Map()
	address := user["address"].(string)
	privateKey := user["privateKey"].(string)

	needCoins := true
	URL := Config().API.GetBalance + "/" + address
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		setError(err.Error(), result)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		setError(err.Error(), result)
	}
	defer resp.Body.Close()

	buf = new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	sjson, err = simplejson.NewJson(buf.Bytes())
	if err != nil {
		setError(err.Error(), result)
	}
	response := map[string]interface{}{}
	response, _ = sjson.Map()

	if value, ok := response["1"]; ok {
		balance, err := strconv.ParseFloat(value.(string), 64)
		if (err == nil) && (balance > 10) {
			needCoins = false
		}
	}

	if needCoins {
		URL = Config().API.GetCoin
		params := map[string]interface{}{
			"addr":   address,
			"amount": "10",
		}
		response = postByJSON(r, URL, params, result)
	}

	URL = Config().API.PrepareProxy
	params := map[string]string{
		"delegates":     Config().Delegate,
		"senderAddress": address,
		"userKey":       address,
	}
	response = postByForm(r, URL, params, result)
	rawTranscation := response["rawTx"].(string)

	URL = Config().API.SignTransaction
	paramsJSON := map[string]interface{}{
		"pri_key": privateKey,
		"raw_tx":  rawTranscation,
	}
	response = postByJSON(r, URL, paramsJSON, result)
	signedTranscation := response["result"].(string)

	URL = Config().API.CreateProxy
	params = map[string]string{
		"rawTxSigned":   signedTranscation,
		"senderAddress": address,
		"userKey":       address,
	}
	response = postByForm(r, URL, params, result)
	contract := response["contract"].(map[string]interface{})
	proxy := contract["proxy"].(string)
	user["proxy"] = proxy
	user["controller"] = contract["controller"].(string)
	user["recovery"] = contract["recovery"].(string)

	account := map[string]string{}
	o := orm.NewOrm()
	o.Using("vchain")
	rows := []orm.Params{}
	sql := "SELECT id FROM `vchain`.`users` WHERE publickey = ? LIMIT 1"
	num, err := o.Raw(sql, user["publicKey"]).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num == 0 {
		now := getNow()
		sql = "INSERT INTO `vchain`.`users`(`name`, `phone`,"
		sql += "`privatekey`, `publickey`, `address`, `proxy`,"
		sql += "`controller`, `recovery`, `created`, `updated`) VALUES("
		sql += "?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		_, err := o.Raw(sql, user["name"], user["phone"],
			user["privateKey"], user["publicKey"],
			user["address"], user["proxy"], user["controller"],
			user["recovery"], now, now).Exec()
		if err != nil {
			setError(err.Error(), result)
		} else {
			account["name"] = user["name"].(string)
			account["phone"] = user["phone"].(string)
			account["privateKey"] = user["privateKey"].(string)
			account["publicKey"] = user["publicKey"].(string)
			account["address"] = user["address"].(string)
			account["proxy"] = user["proxy"].(string)
			account["controller"] = user["controller"].(string)
			account["recovery"] = user["recovery"].(string)
		}
	}

	nodes := map[string]interface{}{}
	result["user"] = user
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func generateLoginToken(rw http.ResponseWriter, r *http.Request) {
	errors := []string{}
	result := map[string]interface{}{}
	result["error"] = errors

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	sjson, err := simplejson.NewJson(buf.Bytes())
	if err != nil {
		setError(err.Error(), result)
	}
	payload := map[string]interface{}{}
	payload, _ = sjson.Map()
	clientJWT := payload["clientJWT"].(string)
	URL := Config().API.DecodeJWT
	params := map[string]interface{}{
		"token": clientJWT,
	}
	response := postByJSON(r, URL, params, result)
	body := response["payload"]
	context := body.(map[string]interface{})["context"].(map[string]interface{})
	publicKey := context["clientPublicKey"].(string)

	URL = Config().API.VerifyJWT
	params = map[string]interface{}{
		"pubkey": publicKey,
		"token":  clientJWT,
	}
	response = postByJSON(r, URL, params, result)

	serverJWT := ""
	token := ""
	if verification, ok := response["result"]; ok {
		if verification == "True" {
			token = strings.Replace(uuid.NewV4().String(), "-", "", -1)
			serverContext := map[string]string{
				"clientName":      context["clientName"].(string),
				"scope":           context["scope"].(string),
				"serverPublicKey": Config().JWT.ServerPublicKey,
				"token":           token,
			}
			iat := time.Now().UTC().Unix()
			exp := iat + 300
			serverJSON := map[string]interface{}{
				"iss":     "vport.chancheng.server",
				"aud":     "vport.chancheng.user",
				"iat":     iat,
				"exp":     exp,
				"sub":     "login token",
				"context": serverContext,
			}

			URL = Config().API.EncodeJWT
			params = map[string]interface{}{
				"payload":     serverJSON,
				"private_key": Config().JWT.ServerPrivateKey,
			}
			response = postByJSON(r, URL, params, result)
			serverJWT = response["token"].(string)
		}
	}

	nodes := map[string]interface{}{}
	result["JWT"] = serverJWT
	result["token"] = token
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func validateUsersLoginJWT(rw http.ResponseWriter, r *http.Request) {
	errors := []string{}
	result := map[string]interface{}{}
	result["error"] = errors

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	sjson, err := simplejson.NewJson(buf.Bytes())
	if err != nil {
		setError(err.Error(), result)
	}
	payload, _ := sjson.Map()
	userJWT := payload["userJWT"].(string)
	URL := Config().API.DecodeJWT
	params := map[string]interface{}{
		"token": userJWT,
	}
	response := postByJSON(r, URL, params, result)
	body := response["payload"]
	valid := false
	expired := true
	exp, err := body.(map[string]interface{})["exp"].(json.Number).Int64()
	if err == nil {
		diff := exp - time.Now().UTC().Unix()
		if diff > 0 {
			expired = false
		}
	}
	context := body.(map[string]interface{})["context"].(map[string]interface{})
	proxy := context["userProxy"].(string)
	publicKey := context["userPublicKey"].(string)
	scope := context["scope"].(string)
	token := context["token"].(string)

	URL = Config().API.VerifyJWT
	params = map[string]interface{}{
		"pubkey": publicKey,
		"token":  userJWT,
	}
	response = postByJSON(r, URL, params, result)

	if (response["result"] == "True") && !expired {
		valid = true
	}

	now := getNow()
	o := orm.NewOrm()
	o.Using("vchain")
	rows := []orm.Params{}
	sql := "SELECT token FROM `vchain`.`tokens` WHERE token = ? LIMIT 1"
	num, err := o.Raw(sql, token).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num == 0 {
		sql = "INSERT INTO `vchain`.`tokens`(`token`, `valid`,"
		sql += "`proxy`, `scope`, `created`) VALUES("
		sql += "?, ?, ?, ?, ?)"
		_, err = o.Raw(sql, token, valid, proxy, scope, now).Exec()
		if err != nil {
			setError(err.Error(), result)
		}
	} else if num > 0 {
		sql := "UPDATE `vchain`.`tokens`"
		sql += " SET `valid` = ?, `proxy` = ?, `scope` = ?,"
		sql += " `created` = ? WHERE token = ?"
		_, err = o.Raw(sql, valid, proxy, scope, now, token).Exec()
		if err != nil {
			setError(err.Error(), result)
		}
	}

	nodes := map[string]interface{}{}
	result["valid"] = valid
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func getAttestationData(proxy string, attestationType string, result map[string]interface{}) map[string]string {
	item := map[string]string{}
	claimID := ""
	o := orm.NewOrm()
	o.Using("vchain")
	rows := []orm.Params{}
	sql := "SELECT id FROM `vchain`.`claims` WHERE proxy = ?"
	sql += " AND status = ? AND type = ? LIMIT 1"
	num, err := o.Raw(sql, proxy, "APPROVED", attestationType).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		row := rows[0]
		claimID = row["id"].(string)
	}
	if len(claimID) > 0 {
		sql = "SELECT attestant, attestation, created"
		sql += " FROM `vchain`.`attestations` WHERE claimid = ? LIMIT 1"
		num, err := o.Raw(sql, claimID).Values(&rows)
		if err != nil {
			setError(err.Error(), result)
		} else if num > 0 {
			row := rows[0]
			item["attestant"] = row["attestant"].(string)
			item["attestation"] = row["attestation"].(string)
			item["created"] = row["created"].(string)
		}
	}
	return item
}

func getUserData(proxy string, scope string, result map[string]interface{}) map[string]string {
	user := map[string]string{}
	o := orm.NewOrm()
	o.Using("vchain")
	rows := []orm.Params{}
	sql := "SELECT name, idnumber, phone, email, privatekey, publickey, address"
	sql += ", proxy, controller, recovery, ipfs, description, created"
	sql += " FROM `vchain`.`users` WHERE proxy = ? LIMIT 1"
	num, err := o.Raw(sql, proxy).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		row := rows[0]
		if strings.Contains(scope, "address") {
			user["address"] = row["address"].(string)
		}
		if strings.Contains(scope, "controller") {
			user["controller"] = row["controller"].(string)
		}
		if strings.Contains(scope, "created") {
			user["created"] = row["created"].(string)
		}
		if strings.Contains(scope, "description") {
			user["description"] = row["description"].(string)
		}
		if strings.Contains(scope, "email") {
			user["email"] = row["email"].(string)
		}
		if strings.Contains(scope, "ID") {
			user["ID"] = row["idnumber"].(string)
		}
		if strings.Contains(scope, "ipfs") {
			user["ipfs"] = row["ipfs"].(string)
		}
		if strings.Contains(scope, "name") {
			user["name"] = row["name"].(string)
		}
		if strings.Contains(scope, "phone") {
			user["phone"] = row["phone"].(string)
		}
		if strings.Contains(scope, "proxy") {
			user["proxy"] = row["proxy"].(string)
		}
		if strings.Contains(scope, "publicKey") {
			user["publicKey"] = row["publickey"].(string)
		}
		if strings.Contains(scope, "recovery") {
			user["recovery"] = row["recovery"].(string)
		}
		if strings.Contains(scope, "updated") {
			user["updated"] = row["updated"].(string)
		}
	}
	return user
}

func getResultForWebsocket(token string) map[string]interface{} {
	errors := []string{}
	result := map[string]interface{}{}
	result["error"] = errors

	valid := false
	proxy := ""
	scope := ""
	o := orm.NewOrm()
	o.Using("vchain")
	rows := []orm.Params{}
	sql := "SELECT token, valid, proxy, scope, created"
	sql += " FROM `vchain`.`tokens` WHERE token = ? LIMIT 1"
	row := orm.Params{}
	for i := 1; i <= 300; i++ {
		num, err := o.Raw(sql, token).Values(&rows)
		if err != nil {
			setError(err.Error(), result)
			break
		} else if num > 0 {
			row = rows[0]
			if row["valid"].(string) == "1" {
				valid = true
			}
			proxy = row["proxy"].(string)
			scope = row["scope"].(string)
			break
		}
		if len(scope) == 0 {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	if valid {
		if scope == "ID" {
			result["ID"] = getAttestationData(proxy, scope, result)
		} else {
			result["user"] = getUserData(proxy, scope, result)
			result["login"] = valid
		}
	}

	nodes := map[string]interface{}{}
	result["scope"] = scope
	nodes["result"] = result
	return nodes
}

func createClaim(rw http.ResponseWriter, r *http.Request) {
	errors := []string{}
	result := map[string]interface{}{}
	result["error"] = errors

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	sjson, err := simplejson.NewJson(buf.Bytes())
	if err != nil {
		setError(err.Error(), result)
	}
	payload, _ := sjson.Map()
	claimJWT := payload["claimJWT"].(string)

	URL := Config().API.DecodeJWT
	params := map[string]interface{}{
		"token": claimJWT,
	}
	response := postByJSON(r, URL, params, result)
	body := response["payload"]
	valid := false
	expired := true
	exp, err := body.(map[string]interface{})["exp"].(json.Number).Int64()
	if err == nil {
		diff := exp - time.Now().UTC().Unix()
		if diff > 0 {
			expired = false
		}
	}
	subject := body.(map[string]interface{})["sub"].(string)
	claimType := strings.Replace(subject, "claim for ", "", -1)
	claimType = strings.ToUpper(claimType)
	context := body.(map[string]interface{})["context"].(map[string]interface{})
	proxy := context["userProxy"].(string)
	publicKey := context["userPublicKey"].(string)

	URL = Config().API.VerifyJWT
	params = map[string]interface{}{
		"pubkey": publicKey,
		"token":  claimJWT,
	}
	response = postByJSON(r, URL, params, result)

	if (response["result"] == "True") && !expired {
		valid = true
	}
	claim := map[string]string{}
	if valid {
		now := getNow()
		o := orm.NewOrm()
		o.Using("vchain")
		rows := []orm.Params{}
		sql := "SELECT id, proxy, type, status FROM `vchain`.`claims`"
		sql += " WHERE proxy = ? AND type = ? AND claim = ? LIMIT 1"
		num, err := o.Raw(sql, proxy, claimType, claimJWT).Values(&rows)
		if err != nil {
			setError(err.Error(), result)
		} else if num == 0 {
			sql = "INSERT INTO `vchain`.`claims`(`proxy`, `type`,"
			sql += "`status`, `claim`, `created`, `updated`) VALUES("
			sql += "?, ?, ?, ?, ?, ?)"
			_, err = o.Raw(sql, proxy, claimType, "PENDING", claimJWT, now, now).Exec()
			if err != nil {
				setError(err.Error(), result)
			} else {
				claim["proxy"] = proxy
				claim["type"] = claimType
				claim["status"] = "PENDING"
				claim["content"] = claimJWT
			}
		} else if num > 0 {
			setError("Claim existed.", result)
		}
	}
	nodes := map[string]interface{}{}
	result["valid"] = valid
	result["claim"] = claim
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func generateClaimToken(rw http.ResponseWriter, r *http.Request) {
	errors := []string{}
	result := map[string]interface{}{}
	result["error"] = errors

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	sjson, err := simplejson.NewJson(buf.Bytes())
	if err != nil {
		setError(err.Error(), result)
	}
	payload := map[string]interface{}{}
	payload, _ = sjson.Map()
	clientJWT := payload["clientJWT"].(string)
	URL := Config().API.DecodeJWT
	params := map[string]interface{}{
		"token": clientJWT,
	}
	response := postByJSON(r, URL, params, result)
	body := response["payload"]
	context := body.(map[string]interface{})["context"].(map[string]interface{})
	publicKey := context["clientPublicKey"].(string)

	URL = Config().API.VerifyJWT
	params = map[string]interface{}{
		"pubkey": publicKey,
		"token":  clientJWT,
	}
	response = postByJSON(r, URL, params, result)

	serverJWT := ""
	token := ""
	if verification, ok := response["result"]; ok {
		if verification == "True" {
			token = strings.Replace(uuid.NewV4().String(), "-", "", -1)
			serverContext := map[string]string{
				"clientName":      context["clientName"].(string),
				"serverPublicKey": Config().JWT.ServerPublicKey,
				"token":           token,
			}
			iat := time.Now().UTC().Unix()
			exp := iat + 300
			serverJSON := map[string]interface{}{
				"iss":     "vport.chancheng.server",
				"aud":     "vport.chancheng.user",
				"iat":     iat,
				"exp":     exp,
				"sub":     "claim token",
				"context": serverContext,
			}

			URL = Config().API.EncodeJWT
			params = map[string]interface{}{
				"payload":     serverJSON,
				"private_key": Config().JWT.ServerPrivateKey,
			}
			response = postByJSON(r, URL, params, result)
			serverJWT = response["token"].(string)
		}
	}

	nodes := map[string]interface{}{}
	result["JWT"] = serverJWT
	result["token"] = token
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func getClaims(rw http.ResponseWriter, req *http.Request) {
	errors := []string{}
	result := map[string]interface{}{}
	result["error"] = errors
	claims := []map[string]string{}
	page := 1
	if slice, ok := req.URL.Query()["page"]; ok {
		if len(slice) > 0 {
			value := slice[0]
			valueInt, err := strconv.Atoi(value)
			if err != nil {
				setError(err.Error(), result)
			} else {
				page = valueInt
			}
		}
	}

	limit := 5
	pages := 0
	total := 0
	o := orm.NewOrm()
	o.Using("vchain")
	sql := "SELECT COUNT(*) FROM `vchain`.`claims` WHERE status = ? AND type = ?"
	var rows []orm.Params
	num, err := o.Raw(sql, "PENDING", "ID").Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		row := rows[0]
		countInt, err := strconv.Atoi(row["COUNT(*)"].(string))
		if err == nil {
			pages = int(math.Ceil(float64(countInt) / float64(limit)))
			total = countInt
		}
	}
	if page > pages {
		page = pages
	}
	offset := (page - 1) * limit

	sqlcmd := "SELECT id, status, claim, created "
	sqlcmd += "FROM `vchain`.`claims` "
	sqlcmd += "WHERE status = ? AND type = ? "
	sqlcmd += "ORDER BY created ASC LIMIT ? OFFSET ?"
	num, err = o.Raw(sqlcmd, "PENDING", "ID", limit, offset).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num > 0 {
		for _, row := range rows {
			claim := map[string]string{
				"claimID": row["id"].(string),
				"claim":   row["claim"].(string),
				"status":  row["status"].(string),
				"created": row["created"].(string),
			}
			claims = append(claims, claim)
		}
	}
	nodes := map[string]interface{}{}
	nodes["count"] = len(claims)
	nodes["currentPage"] = page
	nodes["pages"] = pages
	nodes["total"] = total
	result["items"] = claims
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func createAttestation(rw http.ResponseWriter, r *http.Request) {
	errors := []string{}
	result := map[string]interface{}{}
	result["error"] = errors
	attested := false

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	sjson, err := simplejson.NewJson(buf.Bytes())
	if err != nil {
		setError(err.Error(), result)
	}
	payload, _ := sjson.Map()
	attestant := payload["attestant"].(string)
	attestation := payload["attestation"].(string)
	claimID := payload["claimID"].(string)
	claimType := payload["claimType"].(string)
	proxy := payload["proxy"].(string)
	status := payload["status"].(string)

	URL := Config().API.DecodeJWT
	params := map[string]interface{}{
		"token": attestation,
	}
	response := postByJSON(r, URL, params, result)
	body := response["payload"]
	valid := false
	expired := true
	exp, err := body.(map[string]interface{})["exp"].(json.Number).Int64()
	if err == nil {
		diff := exp - time.Now().UTC().Unix()
		if diff > 0 {
			expired = false
		}
	}
	context := body.(map[string]interface{})["context"].(map[string]interface{})
	publicKey := context["attestantPublicKey"].(string)

	URL = Config().API.VerifyJWT
	params = map[string]interface{}{
		"pubkey": publicKey,
		"token":  attestation,
	}
	response = postByJSON(r, URL, params, result)

	if (response["result"] == "True") && !expired {
		valid = true
	}
	item := map[string]string{}
	if valid {
		o := orm.NewOrm()
		o.Using("vchain")
		rows := []orm.Params{}
		sql := "SELECT id FROM `vchain`.`claims`"
		sql += " WHERE id = ? AND proxy = ? AND type = ? LIMIT 1"
		num, err := o.Raw(sql, claimID, proxy, claimType).Values(&rows)
		if err != nil {
			valid = false
			setError(err.Error(), result)
		} else if num > 0 {
			now := getNow()
			sql = "UPDATE `vchain`.`claims`"
			sql += " SET `status` = ?, `updated` = ?"
			sql += " WHERE id = ?"
			_, err = o.Raw(sql, status, now, claimID).Exec()
			if err != nil {
				valid = false
				setError(err.Error(), result)
			} else {
				item["status"] = status
				item["updated"] = now
				attested = true
			}
			sql = "SELECT id FROM `vchain`.`attestations`"
			sql += " WHERE claimid = ? AND attestant = ? LIMIT 1"
			num, err := o.Raw(sql, claimID, attestant).Values(&rows)
			if err != nil {
				valid = false
				setError(err.Error(), result)
			} else if (num == 0) && (status == "APPROVED") {
				sql = "INSERT INTO `vchain`.`attestations`(`claimid`, `attestant`,"
				sql += "`attestation`, `status`, `created`, `updated`) VALUES("
				sql += "?, ?, ?, ?, ?, ?)"
				_, err = o.Raw(sql, claimID, attestant, attestation, status, now, now).Exec()
				if err != nil {
					valid = false
					setError(err.Error(), result)
				} else {
					item["attestant"] = attestant
					item["attestation"] = attestation
				}
			} else if num > 0 {
				sql = "UPDATE `vchain`.`attestations`"
				sql += " SET `attestation` = ?, `status` = ?, `updated` = ?"
				sql += " WHERE claimid = ?"
				_, err = o.Raw(sql, attestation, status, now, claimID).Exec()
				if err != nil {
					valid = false
					setError(err.Error(), result)
				} else {
					item["attestant"] = attestant
					item["attestation"] = attestation
				}
			}
		}
	}
	nodes := map[string]interface{}{}
	nodes["attested"] = attested
	result["items"] = item
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func getAttestation(rw http.ResponseWriter, req *http.Request) {
	errors := []string{}
	result := map[string]interface{}{}
	result["error"] = errors

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	sjson, err := simplejson.NewJson(buf.Bytes())
	if err != nil {
		setError(err.Error(), result)
	}
	payload, _ := sjson.Map()
	attestationJWT := payload["attestationJWT"].(string)

	URL := Config().API.DecodeJWT
	params := map[string]interface{}{
		"token": attestationJWT,
	}
	response := postByJSON(req, URL, params, result)
	body := response["payload"]
	valid := false
	expired := true
	exp, err := body.(map[string]interface{})["exp"].(json.Number).Int64()
	if err == nil {
		diff := exp - time.Now().UTC().Unix()
		if diff > 0 {
			expired = false
		}
	}
	subject := body.(map[string]interface{})["sub"].(string)
	claimType := strings.Replace(subject, "attestation retrieval for ", "", -1)
	claimType = strings.ToUpper(claimType)
	context := body.(map[string]interface{})["context"].(map[string]interface{})
	proxy := context["userProxy"].(string)
	publicKey := context["userPublicKey"].(string)

	URL = Config().API.VerifyJWT
	params = map[string]interface{}{
		"pubkey": publicKey,
		"token":  attestationJWT,
	}
	response = postByJSON(req, URL, params, result)

	if (response["result"] == "True") && !expired {
		valid = true
	}
	item := map[string]string{}
	status := "ERROR"
	attestation := ""
	if valid {
		o := orm.NewOrm()
		o.Using("vchain")
		sql := "SELECT id, status FROM `vchain`.`claims`"
		sql += " WHERE proxy = ? AND type = ? ORDER BY created DESC LIMIT 1"
		var rows []orm.Params
		num, err := o.Raw(sql, proxy, claimType).Values(&rows)
		if err != nil {
			setError(err.Error(), result)
		} else if num > 0 {
			row := rows[0]
			claimID := row["id"].(string)
			status = row["status"].(string)
			if status == "APPROVED" {
				sql = "SELECT attestation FROM `vchain`.`attestations`"
				sql += " WHERE claimid = ? AND status = ? ORDER BY created DESC LIMIT 1"
				num, err := o.Raw(sql, claimID, "ACTIVE").Values(&rows)
				if err != nil {
					setError(err.Error(), result)
				} else if num > 0 {
					row = rows[0]
					attestation = row["attestation"].(string)
				}
			}
		}
	}
	nodes := map[string]interface{}{}
	item["status"] = status
	if len(attestation) > 0 {
		item["attestation"] = attestation
	}
	result["items"] = item
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func generateAuthorizationToken(rw http.ResponseWriter, r *http.Request) {
	errors := []string{}
	result := map[string]interface{}{}
	result["error"] = errors

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	sjson, err := simplejson.NewJson(buf.Bytes())
	if err != nil {
		setError(err.Error(), result)
	}
	payload := map[string]interface{}{}
	payload, _ = sjson.Map()
	clientJWT := payload["clientJWT"].(string)
	URL := Config().API.DecodeJWT
	params := map[string]interface{}{
		"token": clientJWT,
	}
	response := postByJSON(r, URL, params, result)
	body := response["payload"]
	context := body.(map[string]interface{})["context"].(map[string]interface{})
	publicKey := context["clientPublicKey"].(string)

	URL = Config().API.VerifyJWT
	params = map[string]interface{}{
		"pubkey": publicKey,
		"token":  clientJWT,
	}
	response = postByJSON(r, URL, params, result)

	serverJWT := ""
	token := ""
	if verification, ok := response["result"]; ok {
		if verification == "True" {
			token = strings.Replace(uuid.NewV4().String(), "-", "", -1)
			serverContext := map[string]string{
				"requesterName":   context["requesterName"].(string),
				"scope":           context["scope"].(string),
				"serverPublicKey": Config().JWT.ServerPublicKey,
				"token":           token,
			}
			iat := time.Now().UTC().Unix()
			exp := iat + 300
			serverJSON := map[string]interface{}{
				"iss":     "vport.chancheng.server",
				"aud":     "vport.chancheng.user",
				"iat":     iat,
				"exp":     exp,
				"sub":     "authorization request",
				"context": serverContext,
			}

			URL = Config().API.EncodeJWT
			params = map[string]interface{}{
				"payload":     serverJSON,
				"private_key": Config().JWT.ServerPrivateKey,
			}
			response = postByJSON(r, URL, params, result)
			serverJWT = response["token"].(string)
		}
	}

	nodes := map[string]interface{}{}
	result["JWT"] = serverJWT
	result["token"] = token
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func validateUserAauthorizationJWT(rw http.ResponseWriter, r *http.Request) {
	errors := []string{}
	result := map[string]interface{}{}
	result["error"] = errors

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	sjson, err := simplejson.NewJson(buf.Bytes())
	if err != nil {
		setError(err.Error(), result)
	}
	payload, _ := sjson.Map()
	authorizationJWT := payload["authorizationJWT"].(string)
	URL := Config().API.DecodeJWT
	params := map[string]interface{}{
		"token": authorizationJWT,
	}
	response := postByJSON(r, URL, params, result)
	body := response["payload"]
	valid := false
	expired := true
	exp, err := body.(map[string]interface{})["exp"].(json.Number).Int64()
	if err == nil {
		diff := exp - time.Now().UTC().Unix()
		if diff > 0 {
			expired = false
		}
	}
	context := body.(map[string]interface{})["context"].(map[string]interface{})
	proxy := context["userProxy"].(string)
	publicKey := context["userPublicKey"].(string)
	scope := context["scope"].(string)
	token := context["token"].(string)

	URL = Config().API.VerifyJWT
	params = map[string]interface{}{
		"pubkey": publicKey,
		"token":  authorizationJWT,
	}
	response = postByJSON(r, URL, params, result)

	if (response["result"] == "True") && !expired {
		valid = true
	}

	now := getNow()
	o := orm.NewOrm()
	o.Using("vchain")
	rows := []orm.Params{}
	sql := "SELECT token FROM `vchain`.`tokens` WHERE token = ? LIMIT 1"
	num, err := o.Raw(sql, token).Values(&rows)
	if err != nil {
		setError(err.Error(), result)
	} else if num == 0 {
		sql = "INSERT INTO `vchain`.`tokens`(`token`, `valid`,"
		sql += "`proxy`, `scope`, `created`) VALUES("
		sql += "?, ?, ?, ?, ?)"
		_, err = o.Raw(sql, token, valid, proxy, scope, now).Exec()
		if err != nil {
			setError(err.Error(), result)
		}
	} else if num > 0 {
		sql := "UPDATE `vchain`.`tokens`"
		sql += " SET `valid` = ?, `proxy` = ?, `scope` = ?,"
		sql += " `created` = ? WHERE token = ?"
		_, err = o.Raw(sql, valid, proxy, scope, now, token).Exec()
		if err != nil {
			setError(err.Error(), result)
		}
	}

	nodes := map[string]interface{}{}
	result["valid"] = valid
	nodes["result"] = result
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	setResponse(rw, nodes)
}

func main() {
	cfg := flag.String("c", "cfg.json", "specify config file")
	flag.Parse()
	parseConfig(*cfg)
	db := Config().Database
	orm.RegisterDataBase("default", "mysql", db.Address, db.Idle, db.Max)

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		so.On("connection", func(message string) {
			log.Println("connection message =", message)
			sjson, err := simplejson.NewJson([]byte(message))
			if err != nil {
				so.Emit("error", "Invalid websocket request")
			}
			payload := map[string]interface{}{}
			payload, _ = sjson.Map()
			token := payload["token"].(string)
			log.Println("SOCKET token =", token)
			result := getResultForWebsocket(token)
			log.Println("SOCKET result =", result)
			so.Emit(token, result)
			so.Disconnect()
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
		so.Emit("error", err)
	})

	http.HandleFunc("/api/v1/attestations", getAttestation)
	http.HandleFunc("/api/v1/attestations/add", createAttestation)
	http.HandleFunc("/api/v1/authorizations/jwt", validateUserAauthorizationJWT)
	http.HandleFunc("/api/v1/authorizations/token", generateAuthorizationToken)
	http.HandleFunc("/api/v1/claims", getClaims)
	http.HandleFunc("/api/v1/claims/add", createClaim)
	http.HandleFunc("/api/v1/claims/token", generateClaimToken)
	http.HandleFunc("/api/v1/login/jwt", validateUsersLoginJWT)
	http.HandleFunc("/api/v1/login/token", generateLoginToken)
	http.HandleFunc("/api/v1/users/add", createUser)
	http.Handle("/api/v1/socket/", server)

	port := Config().Port
	addr := "0.0.0.0:" + strconv.Itoa(port)
	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}
	log.Println("http.Start ok, listening on", addr)
	log.Fatalln(s.ListenAndServe())
}
