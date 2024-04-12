package bilibili

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

const ckRaw = `buvid3=B767D3D5-9EBF-29D3-2EF7-29CB3A3311CE68136infoc; b_nut=1712902068; b_lsid=9D8D5A55_18ED0EB5314; _uuid=114E10D14-6E98-59A1-72B8-810C26F5DF2F866968infoc; enable_web_push=DISABLE; buvid_fp=c983ea9e89578d4cacec0fa7685013ab; buvid4=66CC1819-26AF-6887-611A-3491E061709C68718-024041206-2Q6xaryFLHEwY9bCN8w5oA%3D%3D; home_feed_column=4; browser_resolution=1082-1271; FEED_LIVE_VERSION=V_WATCHLATER_PIP_WINDOW2; header_theme_version=CLOSE`

func TestWbiSignQuery(t *testing.T) {
	rURL, err := url.Parse("https://api.bilibili.com/x/space/wbi/acc/info?mid=1850091")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	wbi := NewDefaultWbi()

	q, err := wbi.SignQuery(rURL.Query(), time.Now())
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	rURL.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", rURL.String(), nil)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	//req.Header.Add("Referer", "https://www.bilibili.com/") // 设置会导致失败
	req.Header.Add("Cookie", ckRaw) // 不设置会导致失败
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
		return
	}

	if strings.Contains(string(body), "风控校验失败") {
		t.Fatal("风控校验失败")
		return
	}
}
