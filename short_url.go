package bilibili

import (
	"github.com/pkg/errors"
	"regexp"
	"strconv"
)

var (
	regBv   = regexp.MustCompile(`(?i)bv([\dA-Za-z]{10})`)
	regLive = regexp.MustCompile(`^https://live.bilibili.com/(\d+)`)
)

// UnwrapShortUrl 解析短链接，传入一个完整的短链接。
//
// 第一个返回值如果是"bvid"，则第二个返回值是视频的bvid (string)。
// 第一个返回值如果是"live"，则第二个返回值是直播间id (int)。
func (c *Client) UnwrapShortUrl(shortUrl string) (string, any, error) {
	resp, err := c.resty.R().Get(shortUrl)
	if resp == nil {
		return "", nil, errors.WithStack(err)
	}
	if resp.StatusCode() != 302 {
		return "", nil, errors.Errorf("解析短链接失败，status code: %d", resp.StatusCode())
	}
	url := resp.Header().Get("Location")
	{
		ret := regBv.FindString(url)
		if len(ret) > 0 {
			return "bvid", ret, nil
		}
	}
	{
		ret := regLive.FindStringSubmatch(url)
		if len(ret) > 0 {
			rid, err := strconv.Atoi(ret[1])
			if err != nil {
				return "", nil, errors.WithStack(err)
			}
			return "live", rid, nil
		}
	}
	return "", nil, errors.New("无法解析链接：" + url)
}
