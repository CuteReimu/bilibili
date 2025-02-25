# 命名规范和编码风格

> [!CAUTION]
> 因为本仓库的用户数有一定的人数，因此**尽可能的不要出现 breaking change** ，除非以下情况：
> - B站自己的接口发生了改动，我们需要修改
> - 实在无法避免 breaking change

## 关于命名规范

- 函数命名参考中文翻译，例如“获取全部直播间分区列表”翻译成`GetLiveAreaList`，每个词一一对应，个别词语在不影响理解的情况下可以省略，例如这里“全部”被省略了。
  - 可以参考一下 **bilibili-API-collect** 库的翻译，例如专栏文集叫`Articles`，专栏文章叫`Article`。
- 函数传入参数的结构体名定义为`GetLiveAreaListParam`，也就是后面加“Param”，然后返回值定义为`GetLiveAreaListResult`，也就是后面加“Result”。
  - 如果完全没有传入参数，没必要弄个空的结构体，留空即可。
  - 对于“Get”类的函数，返回值可以简写成`LiveAreaList`，也就是把“Get”和“Result”省略掉。
  - 传入参数中的`csrf`特殊处理一下，不用调用者手动填写，可以参考一下`live.go`。
- 子类型的名字（例如下面例子中的`Category`）可以自己随意修改，只要易读即可，但请注意不要和别的文件中的子类型出现冲突（可以考虑加个前缀）。
- 子类型如果能够复用尽量复用，例如下面例子中的`Category`和`Categories`用了同一个结构体。会在多个文件中复用的类型考虑写到`type.go`中去。
- 可能为`null`值的结构体加`*`，防止`json.Unmarshal`失败。

例子：

```go
type Article struct {
    Id          int        `json:"id"`           // 专栏cvid
    Title       string     `json:"title"`        // 文章标题
    State       int        `json:"state"`        // 0。作用尚不明确
    PublishTime int        `json:"publish_time"` // 发布时间。秒时间戳
    Words       int        `json:"words"`        // 文章字数
    ImageUrls   []ImageUrl `json:"image_urls"`   // 文章封面
    Category    Category   `json:"category"`     // 文章标签
    Categories  []Category `json:"categories"`   // 文章标签列表
    Summary     string     `json:"summary"`      // 文章摘要
    Stats       Stats      `json:"stats"`        // 文章状态数信息
    LikeState   int        `json:"like_state"`   // 是否点赞。0：未点赞。1：已点赞。需要登录(Cookie) 。未登录为0
}
```

> [!TIP]
> [tools目录下](../tools)包含了方便将Markdown表格转化为Go struct的工具，强烈建议使用。

## 关于go文件的命名

对于**bilibili-API-collect**库的接口，将对应方法放在对应的`.go`文件里。例如：

- **bilibili-API-collect**中的`docs/user`文件夹下面的所有接口，对应方法放在`user.go`文件里
- **bilibili-API-collect**中的`docs/video`文件夹下面的所有接口，对应方法放在`video.go`文件里

对于并非是对应接口的方法（例如一些纯工具类方法），可自行找一个合适的`.go`文件。

## 关于编码风格

本项目不做太多的编码风格限制，就提这样一个建议吧：

- 如果是纯粹的接口调用函数，请参考现有的函数写法。
- 如果是其它函数，不限制编码风格，在提交pull request时会带有golangci-lint检测，请确保通过即可。
