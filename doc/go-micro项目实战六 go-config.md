### go-config
你可能已经发现了，我们之前的项目里的所有数据库连接和一些其他需要配置的东西我们都硬编码在代码里了。这并不合理。因此我们这一节来引入go-config，来解决这个问题。

go-config 官方文档说这是一个动态的可插拔的配置库。

### go-config的使用
#### 1.在项目根目录创建config.json文件
```
{
  "mysql" : {
    "host" : "192.168.0.111",
    "port" : "3306",
    "user" : "mytestroot",
    "password" : "mytestroot",
    "database" : "shopping"
  }
}
```

#### 2.修改main.go
引入go-config
```
import(
	...
	"github.com/micro/go-config"
	...
)
```
在main func 里引用
加载配置文件
```
err := config.LoadFile("./config.json")
if err != nil {
	log.Fatalf("Could not load config file: %s", err.Error())
	return
}
conf := config.Map()
```
修改以前的database.go里的数据库连接
```
func CreateConnection(dbconf map[string]interface{}) (*gorm.DB, error) {
	host := dbconf["host"]
	port := dbconf["port"]
	user := dbconf["user"]
	dbName := dbconf["database"]
	password := dbconf["password"]
	return gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, host, port, dbName,
	),
	)
}
```
将配置信息传进来
```
//db
db, err := CreateConnection(conf["mysql"].(map[string]interface{}))
defer db.Close()
```

### 总结
go-config的使用没有什么难点。
golang在表示json数据的时候，通常会用map[string]interface{}来表示。如果json结构是一个多维的话，是无法通过map[index1][index2]来取得值的。会编译不通过，提示interface{} 不是map。可以使用强转来将interface转成map[string]interface{}。或者可以使用config.Get(...path)，来获取最终配置的值。Get方法的实现是也是循环输入进去的路径，然后每一层都将interface自动转成map[string]interface{}来调的。