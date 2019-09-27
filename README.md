# fake_orm
伪ORM，目的为了简化操作，不损失性能，不使用反射，操作起来更简便

# 使用方法

### 示例模型
```
    type project struct {
		id int64
		name string
	}
```

### 获取一条数据
```
    
	project := &project{}
	//MasterDB 初始化好的数据库连接
	err = fake_orm.GetDb(MasterDB).Tab("project").Select("id","name").Where("id", ">", 100).OrderBy("id", "desc").First(&project.id, &project.name)
	if err != nil {
		panic(err)
	}
	fmt.Println(project)
```


### 获取多条数据
```

	project := make([]*project, 0, 10)
	//MasterDB 初始化好的数据库连接
	project := make([]*raisefund, 0, 10)
    err = fake_orm.GetDb(MasterDB).Tab("project").Select("id","name").Where("id", ">", 30287016).OrderBy("id", "desc").Get(func(rows *sql.Rows) {
        obj := new(project)
        _ = rows.Scan(&obj.id, &obj.name)
        project = append(project, obj)
    })
    if err != nil {
        panic(err)
    }
    fmt.Println(project)
```

### 条件组合
```
    //多条件
    err = fake_orm.GetDb(MasterDB).Tab("project")
            .Select("id", "name")
            .Where("id", ">", 30287016)
            .Where("name", "=", "张三")
            .OrderBy("id", "desc")
    
    //多条件
    err = fake_orm.GetDb(MasterDB).Tab("project")
            .SelectRaw("id, name")
            .WhereRaw("id > 30287016")
            .Where("name", "=", "张三")
            .OrderBy("id", "desc")
    
    //join
    err = fake_orm.GetDb(MasterDB).Tab("users").join("group","user.group_id = group.id")
            .Select("users.id","users.name","group.name")
            .Where("users.id", ">", 30287016)
            .Where("users.name", "=", "张三")
            .OrderBy("users.id", "desc")
    
    //join as
    err = fake_orm.GetDb(MasterDB).Tab("users as u").join("group as g","u.group_id = g.id")
            .Select("u.id","u.name","g.name")
            .Where("u.id", ">", 30287016)
            .Where("u.name", "=", "张三")
            .OrderBy("u.id", "desc")
 
    //更多功能还在开发中
    
```
