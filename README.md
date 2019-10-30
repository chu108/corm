# corm

corm是一个伪ORM，目的为了简化操作，不损失性能，不使用反射，操作起来更简便

## 支持的操作

* Select、SelectRaw
* Where、WhereRaw、WhereIn、WhereNotIn、WhereLike、WhereNotLike、WhereBetween
* Join、LeftJoin、RightJoin
* OrderBy、GroupBy、Limit、Having
* First、Get
* Sum、Max、Min、Count、Exists
* Insert、Update
* PrintSql

## 示例模型
```
type project struct {
    id int64
    name string
}
```

## 示例数据连接

```go
//MasterDB 初始化好的数据库连接
MasterDB
```

### 获取一条数据
```
project := &project{}
err = corm.GetDb(MasterDB)
        .Tab("project")
        .Select("id","name")
        .Where("id", ">", 100)
        .OrderBy("id", "desc")
        .First(&project.id, &project.name)
if err != nil {
    panic(err)
}
//读取结果
fmt.Println(project)
```


### 获取多条数据
```

//返回Slicp类型数据
project := make([]*project, 0, 10)
err = corm.GetDb(MasterDB).Tab("project")
        .Select("id","name")
        .Where("id", ">", 30287016)
        .OrderBy("id", "desc")
        .Get(func(rows *sql.Rows) {
            obj := new(project)
            _ = rows.Scan(&obj.id, &obj.name)
            project = append(project, obj)
        })
if err != nil {
    panic(err)
}
//读取结果
for k, v := range project {
    fmt.Println(k, v)
}

//返回Map类型数据
project := make(map[int64]*project)
err = corm.GetDb(MasterDB).Tab("project")
        .Select("id","name")
        .Where("id", ">", 30287016)
        .OrderBy("id", "desc")
        .Get(func(rows *sql.Rows) {
            obj := new(project)
            _ = rows.Scan(&obj.id, &obj.name)
            project[obj.id] = obj
        })
if err != nil {
    panic(err)
}
//读取结果
for k, v := range project {
    fmt.Println(k, v)
}

```

### 条件组合
```
//多条件
err = corm.GetDb(MasterDB).Tab("project")
    .Select("id", "name")
    .Where("id", ">", 30287016)
    .Where("name", "=", "张三")
    .OrderBy("id", "desc")

//多条件
err = corm.GetDb(MasterDB).Tab("project")
    .SelectRaw("id, name")
    .WhereRaw("id > 30287016")
    .Where("name", "=", "张三")
    .OrderBy("id", "desc")

//join
err = corm.GetDb(MasterDB).Tab("users").join("group","user.group_id = group.id")
    .Select("users.id","users.name","group.name")
    .Where("users.id", ">", 30287016)
    .Where("users.name", "=", "张三")
    .OrderBy("users.id", "desc")

//join as
err = corm.GetDb(MasterDB).Tab("users as u").join("group as g","u.group_id = g.id")
    .Select("u.id","u.name","g.name")
    .Where("u.id", ">", 30287016)
    .Where("u.name", "=", "张三")
    .OrderBy("u.id", "desc")
    
```

## 插入记录
```go
insertId, err := corm.GetDb(MasterDB).Tab("users")
    .Insert(map[string]interface{}{
        "nickname":"aaaaaaa",
        "name":"aaaaaaa",
        "phone":1231231234,
    })
if err != nil {
    panic(err)
}
//返回插入ID
fmt.Println(insertId)
```

### 修改数据
```go
num, err := corm.GetDb(MasterDB).Tab("users")
    .Where("name", "=", "aaaaaaa")
    .WhereIn("id", 212, 213)
    .Update(map[string]interface{}{
        "nickname":   "eeeeeeeeeee",
        "name":       "eeeeeeeeeee",
        "phone":      "eeeeeeeeeee",
    })
if err != nil {
    panic(err)
}
//更新行数
fmt.Println(num)
```

### 统计查询
```go
//Count
countNum, err := corm.GetDb(MasterDB).Tab("users")
    .Where("name", "=", "aaaaaaa")
    .WhereIn("id", 212, 213)
    .Count()
if err != nil {
    panic(err)
}
fmt.Println(countNum)

//Sum
sumNum, err := corm.GetDb(MasterDB).Tab("users")
    .Where("name", "=", "aaaaaaa")
    .WhereIn("id", 212, 213)
    .Sum("money")
if err != nil {
    panic(err)
}
fmt.Println(sumNum)

//Max
maxNum, err := corm.GetDb(MasterDB).Tab("users")
    .Where("name", "=", "aaaaaaa")
    .WhereIn("id", 212, 213)
    .Max("money")
if err != nil {
    panic(err)
}
fmt.Println(maxNum)

//Min
minNum, err := corm.GetDb(MasterDB).Tab("users")
    .Where("name", "=", "aaaaaaa")
    .WhereIn("id", 212, 213)
    .Min("money")
if err != nil {
    panic(err)
}
fmt.Println(minNum)
```