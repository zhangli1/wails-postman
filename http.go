package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"myproject/sqliteDb"
	"myproject/tools"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var GlobalFeildMap = make(map[string]string, 10)

func Server() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/saveDir", handleSaveDir)
	http.HandleFunc("/getDirList", handleGetDirList)
	http.HandleFunc("/save", handleSave)
	http.HandleFunc("/delRequest", handleDelRequest)
	http.HandleFunc("/delDir", handleDelDir)
	http.HandleFunc("/saveGlobalFeild", handlesaveGlobalFeild)
	http.HandleFunc("/getGlobalFeild", handlegetGlobalFeild)

	// 启动服务器
	port := "8089"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	log.Printf("Server started on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// 服务首页
func serveHome(w http.ResponseWriter, r *http.Request) {

}

// 添加目录
func handleSaveDir(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	// 仅允许 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body) //读取一次后，无法再次读取，所以要重置一道

	var params struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	err := json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %s", err), http.StatusBadRequest)
		return
	}

	// 重置 r.Body，以便后续可以重新读取
	r.Body = io.NopCloser(bytes.NewReader(body))

	name := params.Name
	if name == "" {
		fmt.Fprintf(w, "param error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	type ApiResponse struct {
		Status  int
		Message string
	}
	var retData ApiResponse
	retData.Status = 500

	if params.Id == 0 {
		collectionsInfo := sqliteDb.Collections{
			ObjectId:   1,
			Name:       name,
			CreateTime: tools.GetCurrentTimeFomat(),
			UpdateTime: tools.GetCurrentTimeFomat(),
		}

		collectionsRet, err := sqliteDb.NewCollectionsModel().Create(collectionsInfo)
		if collectionsRet == 0 || err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(retData)
			return
		}
	} else {
		collectionsUpdate := sqliteDb.Collections{
			Name:       name,
			UpdateTime: tools.GetCurrentTimeFomat(),
		}
		sqliteDb.NewCollectionsModel().Update(int64(params.Id), collectionsUpdate)
	}

	w.WriteHeader(http.StatusOK)

	data := ApiResponse{
		Status:  200,
		Message: "Hello, Go HTTP Server!",
	}

	// 编码并返回 JSON
	json.NewEncoder(w).Encode(data)

	return
}

// 替换全局变量
func ReplaceGlobalFeild(content string) string {
	if len(GlobalFeildMap) > 0 {
		for k, v := range GlobalFeildMap {
			kk := fmt.Sprintf("{{%v}}", k)
			content = strings.Replace(content, kk, v, -1)
		}
	}
	return content
}

func handleSave(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	// 仅允许 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body) //读取一次后，无法再次读取，所以要重置一道

	var params struct {
		FolderId int64  `json:"folderId"`
		Id       int64  `json:"id"`
		Name     string `json:"name"`
		Addr     string `json:"addr"`
		Content  string `json:"content"`
	}

	err := json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %s", err), http.StatusBadRequest)
		return
	}

	// 重置 r.Body，以便后续可以重新读取
	r.Body = io.NopCloser(bytes.NewReader(body))

	addr := params.Addr
	content := params.Content
	folderId := params.FolderId
	name := params.Name
	id := params.Id
	if addr == "" || folderId == 0 || name == "" {
		fmt.Fprintf(w, "param error")
		return
	}

	//randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	//// 生成 0 到 99 之间的随机数
	//randomValue := randGen.Intn(100)

	requestObj := sqliteDb.NewRequestModel()

	var retId int64 = 0
	if id == 0 {
		request := sqliteDb.Request{
			Active:        1,
			Addr:          addr,
			CollectionsId: folderId,
			Name:          name,
			CreateTime:    tools.GetCurrentTimeFomat(),
			UpdateTime:    tools.GetCurrentTimeFomat(),
		}
		if content != "" {
			//替换全局变量
			request.Content = tools.JsonCompression(content)
		}
		retId, _ = requestObj.Create(request)
	} else {
		request := sqliteDb.Request{
			Active:        1,
			Addr:          addr,
			CollectionsId: folderId,
			Name:          name,
			CreateTime:    tools.GetCurrentTimeFomat(),
			UpdateTime:    tools.GetCurrentTimeFomat(),
		}
		if content != "" {
			request.Content = tools.JsonCompression(content)
		}
		requestObj.Update(id, request)

		retId = id
	}

	var updateMap = make(map[string]interface{}, 0)
	updateMap["active"] = 0

	requestObj.UpdateByNotMap(retId, updateMap)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	data := ApiResponse{
		Status:  200,
		Message: "Hello, Go HTTP Server!",
	}

	// 编码并返回 JSON
	json.NewEncoder(w).Encode(data)

	return
}

// 允许跨域
func allowCrossOrigin(w http.ResponseWriter, r *http.Request) {
	// 1. 核心：设置所有必要的 CORS 头（必须在 w.Write 之前）
	// 允许所有源（开发环境用 *，生产环境指定具体域名如 "https://xxx.com"）
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 允许的请求方法（覆盖常用方法，按需添加）
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	//允许的请求头（若前端传了自定义头，需在这里添加，如 Token、Content-Type）
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	// 允许前端读取的响应头（可选）
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
	// 预检请求缓存时间（秒，减少 OPTIONS 请求次数）
	w.Header().Set("Access-Control-Max-Age", "86400")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func handlegetGlobalFeild(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	// 2. 处理预检请求 OPTIONS：直接返回 204 No Content（无需响应体）
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// 仅允许 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	globalFieldObj := sqliteDb.NewGlobalFieldModel()

	list, _ := globalFieldObj.GetList()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 编码并返回 JSON
	json.NewEncoder(w).Encode(list)

	return
}

func handlesaveGlobalFeild(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	// 仅允许 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body) //读取一次后，无法再次读取，所以要重置一道

	type paramFeild struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	var params []paramFeild

	err := json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %s", err), http.StatusBadRequest)
		return
	}

	// 重置 r.Body，以便后续可以重新读取
	r.Body = io.NopCloser(bytes.NewReader(body))

	if len(params) == 0 {
		fmt.Fprintf(w, "param error")
		return
	}

	globalFieldObj := sqliteDb.NewGlobalFieldModel()
	globalFieldObj.BatchDelete()

	var feilds []sqliteDb.GlobalField
	for _, v := range params {
		feilds = append(feilds, sqliteDb.GlobalField{
			Field: v.Name,
			Value: v.Value,
		})
	}

	globalFieldObj.Create(feilds)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	data := ApiResponse{
		Status:  200,
		Message: "Hello, Go HTTP Server!",
	}

	// 编码并返回 JSON
	json.NewEncoder(w).Encode(data)

	return
}

type ApiResponse struct {
	Status  int
	Message string
}

const htmlFile string = "index5.html"

// 提供 index3.html 文件
func serveIndex(w http.ResponseWriter, subFS fs.FS) {
	data, err := fs.ReadFile(subFS, htmlFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("无法读取 %s", htmlFile), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

// 设置正确的 Content-Type
func setContentType(w http.ResponseWriter, path string) {
	switch filepath.Ext(path) {
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".json":
		w.Header().Set("Content-Type", "application/json")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	}
}

type FolderRequestStruct struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Url     string `json:"url"`
	Content string `json:"content"`
	Active  bool   `json:"active"`
}

type FolderStruct struct {
	Id       string                `json:"id"`
	Name     string                `json:"name"`
	Expanded bool                  `json:"expanded"`
	Requests []FolderRequestStruct `json:"requests"`
	ChildNum int                   `json:"childNum"`
}

// 获取目录列表
func handleGetDirList(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	// 仅允许 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	type ApiResponse struct {
		Status  int
		Message string
	}
	var retData ApiResponse
	retData.Status = 500

	list, err := sqliteDb.NewCollectionsModel().GetList()
	log.Println("debug: %v %v", list, err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(retData)
		return
	}

	var folders []FolderStruct

	requestObj := sqliteDb.NewRequestModel()
	for _, item := range list {
		requestList, err := requestObj.FindByCollectionsId(item.ID)
		if err != nil {
			continue
		}

		var expanded bool = false
		var req = make([]FolderRequestStruct, 0)
		for _, reqItem := range requestList {
			var active bool = false
			if reqItem.Active > 0 {
				active = true
				expanded = true

			}
			req = append(req, FolderRequestStruct{
				Id:      tools.Int64ToString(reqItem.ID),
				Name:    reqItem.Name,
				Url:     reqItem.Addr,
				Content: reqItem.Content,
				Active:  active,
			})
		}

		folders = append(folders, FolderStruct{
			Id:       tools.Int64ToString(item.ID),
			Name:     item.Name,
			Expanded: expanded,
			Requests: req,
			ChildNum: len(req),
		})

	}

	w.WriteHeader(http.StatusOK)

	// 编码并返回 JSON
	json.NewEncoder(w).Encode(folders)

	return
}

// 删除请求
func handleDelRequest(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	// 仅允许 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body) //读取一次后，无法再次读取，所以要重置一道

	var params struct {
		Id int64 `json:"id"`
	}

	err := json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %s", err), http.StatusBadRequest)
		return
	}

	// 重置 r.Body，以便后续可以重新读取
	r.Body = io.NopCloser(bytes.NewReader(body))

	id := params.Id
	if id == 0 {
		fmt.Fprintf(w, "param error")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	requestObj := sqliteDb.NewRequestModel()
	count, err := requestObj.Delete(id)
	if count == 0 || err != nil {
		fmt.Fprintf(w, "del error")
		return
	}

	w.WriteHeader(http.StatusOK)

	data := ApiResponse{
		Status:  200,
		Message: "Hello, Go HTTP Server!",
	}

	// 编码并返回 JSON
	json.NewEncoder(w).Encode(data)
	return
}

// 删除目录
func handleDelDir(w http.ResponseWriter, r *http.Request) {
	allowCrossOrigin(w, r)
	// 仅允许 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body) //读取一次后，无法再次读取，所以要重置一道

	var params struct {
		Id int64 `json:"id"`
	}

	err := json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %s", err), http.StatusBadRequest)
		return
	}

	// 重置 r.Body，以便后续可以重新读取
	r.Body = io.NopCloser(bytes.NewReader(body))

	id := params.Id
	if id == 0 {
		fmt.Fprintf(w, "param error")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	collectionsObj := sqliteDb.NewCollectionsModel()
	count, err := collectionsObj.Delete(id)
	if count == 0 || err != nil {
		fmt.Fprintf(w, "del error")
		return
	}

	w.WriteHeader(http.StatusOK)

	data := ApiResponse{
		Status:  200,
		Message: "Hello, Go HTTP Server!",
	}

	// 编码并返回 JSON
	json.NewEncoder(w).Encode(data)
	return
}
