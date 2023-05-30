package proxy

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"server_template/service/account"
	"strings"
	"time"
)

type CustomProxyConfig struct {
	ProxyServer string `json:"proxy_server"`
	Key         string `json:"key"`
	Timeout     int64  `json:"timeout"`
}

var (
	proxyConfig CustomProxyConfig
	httpClient  *http.Client
)

func InitCustomProxy(config CustomProxyConfig) {
	proxyConfig = config
	httpClient = &http.Client{
		//Transport: tr,
		Timeout: time.Duration(int64(time.Second) * proxyConfig.Timeout), //超时时间
	}
}

//https://xxx.com/proxy?key=&url=
func HttpProxyHandler(c *gin.Context) {

	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Referer, X-SMTPUsername-Agent")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == http.MethodOptions {
		c.Writer.WriteHeader(200)
		return
	}

	logged, _ := account.IsLogged(c)
	if !logged {
		return
	}

	newUrl, _ := url.Parse(proxyConfig.ProxyServer)

	q := newUrl.Query()
	q.Add("url", c.Query("url"))
	q.Add("key", proxyConfig.Key)
	newUrl.RawQuery = q.Encode()

	req, _ := http.NewRequest(c.Request.Method, newUrl.String(), c.Request.Body)
	copyRequestHeader(c, req)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	copyResponseHeader(c, resp)
	buf := make([]byte, 128)
	io.CopyBuffer(c.Writer, resp.Body, buf)
}

func copyResponseHeader(c *gin.Context, resp *http.Response) {
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		cookie.SameSite = http.SameSiteNoneMode
		cookie.Secure = true
		c.Writer.Header().Add("Set-Cookie", cookie.String())
	}
	c.Writer.WriteHeader(resp.StatusCode)
	for k, vList := range resp.Header {
		if k != "Set-Cookie" {
			for _, v := range vList {
				c.Writer.Header().Add(k, v)
			}
		}
	}
}

func copyRequestHeader(c *gin.Context, req *http.Request) {
	for k, vList := range c.Request.Header {
		newKey := k
		if strings.HasPrefix(k, "X-") {
			newKey = k[len("X-"):]
			c.Request.Header.Del(newKey)
			req.Header.Del(newKey)
		}
		for _, v := range vList {
			req.Header.Add(newKey, v)
		}
	}
}
