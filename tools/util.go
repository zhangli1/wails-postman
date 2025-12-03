package tools

import (
	"bytes"
	"fmt"
	"github.com/bytedance/sonic"
	"io"
	"math/rand"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 将结构体转换为 JSON串
func StructToJsonStr(input interface{}) string {
	jsonData, err := sonic.Marshal(&input)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

// 获取当时时间
func GetCurrentTimeFomat() *time.Time {
	// 获取当前时间
	currentTime := time.Now()
	t := time.Unix(currentTime.Unix(), 0)
	return &t
}

// 获取当时时间
func GetCurrentTime(format ...string) string {
	var f string = "2006-01-02 15:04:05"
	if len(format) > 0 {
		f = format[0]
	}
	// 获取当前时间
	currentTime := time.Now()
	t := time.Unix(currentTime.Unix(), 0)
	// 格式化输出：年月日时分秒
	return t.Format(f)
}

// 获取当前的 Unix 时间戳（秒）
func GetTimestamp() int64 {
	// 获取当前时间
	currentTime := time.Now()
	// 获取当前时间的 Unix 时间戳（秒）
	return currentTime.Unix()
}

// 获取当前的 Unix 时间戳（毫秒）
func GetUnixMilli() int64 {
	// 获取当前时间
	currentTime := time.Now()
	// 获取当前时间的 Unix 时间戳（毫秒）
	return currentTime.UnixMilli()
}

// 获取当前的 Unix 时间戳（微秒）
func GetUnixMicro() int64 {
	// 获取当前时间
	currentTime := time.Now()
	// 获取当前时间的 Unix 时间戳（微秒）
	return currentTime.UnixMicro()
}

// 时间戳（秒）转时间格式
func TimestampToDateTime(timestamp int64) *time.Time {
	t := time.Unix(timestamp, 0)
	return &t
}

// 时间戳（秒）转时间格式(输入的格式)
func TimestampToDateTimeByFormat(timestamp int64, format ...string) string {
	var f string = "2006-01-02 15:04:05"
	if len(format) > 0 {
		f = format[0]
	}
	t := time.Unix(timestamp, 0)
	// 格式化输出：年月日时分秒
	return t.Format(f)
}

// string to int64
func StringToInt64(str string) int64 {
	var ret int64 = 0
	if str == "" {
		return ret
	}

	if strings.Contains(str, ".") {
		floatVal, _ := strconv.ParseFloat(str, 64)
		return int64(floatVal)
	}

	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Sprintf("%v %v", err, string(debug.Stack()))
		return ret
	}

	return num
}

// string to int
func StringToInt(str string) int {
	var ret int = 0
	if str == "" {
		return ret
	}
	if strings.Contains(str, ".") {
		floatVal, _ := strconv.ParseFloat(str, 64)
		return int(floatVal)
	}
	num, err := strconv.Atoi(str)
	if err != nil {
		fmt.Sprintf("%v %v", err, string(debug.Stack()))
		return ret
	}
	return num
}

// int to String
func IntToString(num int) string {
	var str string = "0"
	str = strconv.Itoa(num)
	return str
}

// int64 to string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

// string to float64
func StringToFloat64(str string) float64 {
	var ret float64 = 0
	if str == "" {
		return ret
	}
	floatValue, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Sprintf("StringToFloat64 fail %v", string(debug.Stack()))
		return ret
	}
	return floatValue
}

// float64 to string
func Float64ToString(num float64) string {
	return strconv.FormatFloat(num, 'f', -1, 64) // 自动选择精度
}

// 判断sync.Map是否为空
func IsEmpty(m *sync.Map) bool {
	empty := true
	m.Range(func(key, value interface{}) bool {
		empty = false
		return false // 一旦找到元素，就可以停止遍历
	})
	return empty
}

// 概率算法
func Probability(probability string) bool {
	value := StringToInt(probability)

	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 生成 0 到 99 之间的随机数
	randomValue := randGen.Intn(100)

	// 判断事件是否命中（根据概率判断）
	if randomValue < value {
		return true
	}

	return false
}

// 获取唯一ID
func GetUniqId() string {
	return Int64ToString(time.Now().UnixNano())
}

// 数组连接成字符串
func ConcatWithBuilder(strs []string) string {
	var builder strings.Builder
	// 预分配总长度（关键！）
	total := 0
	for _, s := range strs {
		total += len(s)
	}
	builder.Grow(total) // 减少扩容次数

	for _, s := range strs {
		builder.WriteString(s)
	}
	return builder.String()
}

// 获取客户端IP
func GetRealIP(r *http.Request) string {
	// 从标准头获取
	ip := r.Header.Get("X-Real-IP")
	//log.Debug("GetRealIP X-Real-IP: %v", ip)

	// 从代理头获取
	if ip == "" {
		ips := r.Header.Get("X-Forwarded-For")
		//log.Debug("GetRealIP X-Forwarded-For: %v", ip)
		if ips != "" {
			for _, i := range strings.Split(ips, ",") {
				ip = strings.TrimSpace(i)
				if ip != "" {
					break
				}
			}
		}
	}

	// 回退到远程地址
	if ip == "" {
		//log.Debug("RemoteAddr: %v", r.RemoteAddr)
		if strings.Contains(r.RemoteAddr, ":") {
			ip = strings.Split(r.RemoteAddr, ":")[0]
		} else {
			ip = r.RemoteAddr
		}
	}

	if ip == "[" {
		ip = ""
	}

	return ip
}

// http post请求
func HttpPost(url string, postData interface{}, headerMap map[string]string) (bool, string) {
	jsonData, err := sonic.Marshal(&postData)
	if err != nil {
		fmt.Sprintf("HttpPost request params fail, %v %v %v %v", url, postData, err, string(debug.Stack()))
		return false, ""
	}

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Sprintf("HttpPost request initfail, %v %v %v %v", url, postData, err, string(debug.Stack()))
		return false, ""
	}

	req.Header.Set("Content-Type", "application/json")
	if len(headerMap) > 0 {
		for k, v := range headerMap {
			req.Header.Set(k, v)
		}

	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Sprintf("HttpPost request fail, %v %v %v %v", url, postData, err, string(debug.Stack()))
		return false, ""
	}
	defer resp.Body.Close()

	//log.Debug("Status Code: %v", resp.Status)
	//for k, v := range resp.Header {
	//	log.Debug("%s: %v\n", k, v)
	//}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Sprintf("HttpPost request response fail, %v %v %v %v", url, postData, err, string(debug.Stack()))
		return false, ""
	}
	return true, string(body)

}

// json串压缩
func JsonCompression(jsonStr string) string {
	// 解析 JSON 到空接口
	var data interface{}
	if err := sonic.Unmarshal([]byte(jsonStr), &data); err != nil {
		panic(err)
	}

	// 重新编组为紧凑格式
	compactJSON, err := sonic.Marshal(data)
	if err != nil {
		panic(err)
	}

	return string(compactJSON)
}
