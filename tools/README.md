# 生成工具的使用方法

```bash
# go 用这个
go run gen_struct.go

# python3用这个
python3 gen_struct.py

# python2用这个
python2 gen_struct2.py
```

会得到这样的提示：

```console
请输入Markdown表格，在最后一行之后输入ok表示结束：（退出请输入exit）
```

把以下Markdown表格复制粘贴输入进去并回车：

```md
| 字段         | 类型  | 内容           | 备注                                                         |
| ------------ | ----- | -------------- | ------------------------------------------------------------ |
| id           | num   | 专栏cvid       |                                                              |
| title        | str   | 文章标题       |                                                              |
| state        | num   | 0              | 作用尚不明确                                                 |
| publish_time | num   | 发布时间       | 秒时间戳                                                     |
| words        | num   | 文章字数       |                                                              |
| image_urls   | array | 文章封面       |                                                              |
| category     | obj   | 文章标签       |                                                              |
| categories   | array | 文章标签列表   |                                                              |
| summary      | str   | 文章摘要       |                                                              |
| stats        | obj   | 文章状态数信息 |                                                              |
| like_state   | num   | 是否点赞       | 0：未点赞<br />1：已点赞<br />需要登录(Cookie) <br />未登录为0 |
```

然后输入：

```console
ok
```

并回车，你就会得到：

```go
type T struct {
    Id int `json:"id"` // 专栏cvid
    Title string `json:"title"` // 文章标题
    State int `json:"state"` // 0。作用尚不明确
    PublishTime int `json:"publish_time"` // 发布时间。秒时间戳
    Words int `json:"words"` // 文章字数
    ImageUrls []ImageUrl `json:"image_urls"` // 文章封面
    Category Category `json:"category"` // 文章标签
    Categories []Categorie `json:"categories"` // 文章标签列表
    Summary string `json:"summary"` // 文章摘要
    Stats Stats `json:"stats"` // 文章状态数信息
    LikeState int `json:"like_state"` // 是否点赞。0：未点赞。1：已点赞。需要登录(Cookie) 。未登录为0
}
```

自行复制到go代码中去，把`T`改个名即可。
