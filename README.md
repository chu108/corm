![](static/corm.png) 
<br>
corm是一个不使用反射的orm, golang版本: > 1.13
<br><br>
![](https://img.shields.io/badge/corm-tools-orange?style=flat-square&logo=appveyor)
![](https://img.shields.io/badge/release-linux%2Cmac%2Cwin-blue?style=flat-square&logo=appveyor)
![](https://img.shields.io/badge/corm-V1.1-Stable)
![](https://img.shields.io/github/stars/chu108/corm)
![](https://img.shields.io/github/issues/chu108/corm)
![](https://img.shields.io/github/forks/chu108/corm)
![](https://img.shields.io/github/license/chu108/corm)

## 支持的操作
- 查询
    - Select 普通查询
    - SelectRaw 原生查询
- where 条件
    - Where
    - WhereRaw 原生where条件
    - WhereIn
    - WhereNotIn
    - WhereLike
    - WhereNotLike
    - WhereBetween
    - WhereStrToInt
    - WhereInt64ToStr
    - WhereIntToStr
- 关联查询 Join
    - Join
    - LeftJoin
    - RightJoin
- 排序
    - OrderBy
- 分组查询
    - GroupBy
- 数据限定
    - Limit
    - Offset
- 查询数据
    - First
    - Get
    - GetPage 分页
    - Exists 是否存在
    - ValueStr 以 string 形式返回指定字段
    - ValueInt 以 int 形式返回指定字段
    - ValueFloat 以 float64 形式返回指定字段
- 聚合查询
    - Sum
    - Max
    - Min
    - Count
- 插入更新
    - Insert
    - Update
    - Transaction 事务支持
- 打印SQL
    - PrintSql

## 示例
```go
package main

import (
	"database/sql"
	"fmt"
	"github.com/chu108/corm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var MasterDB *sql.DB
var cerr error
//用户
type Users struct {
	Id        int64
	Name      string
	Age       int64
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	GroupId   int64
}
//组
type Groups struct {
	Id        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
//用户与组关联
type UserGroups struct {
	Id        int64
	UserId    int64
	GroupId   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func init() {
	MasterDB, cerr = sql.Open("mysql", "corm:666666@tcp(127.0.0.1:3306)/corm_demo")
	echoErr(cerr)
}

func main() {
	//查询
	selectOne()
	selectTwo()
	//条件组合
	whereOne()
	whereTwo()
	//关联查询
	joinOne()
	joinTwo()
	//单条查询、分页查询
	first()
	getPage()
	//聚合查询
	count()
	//插入、更新
	insert()
	update()
	//验证结果是否存在
	exists()
	//事务
	trans()
}

func selectOne() {
	data := make([]*Users, 0)
	cerr = corm.GetDb(MasterDB).Tab("users").Select("name", "age", "phone").Where("age", ">=", 24).Get(func(rows *sql.Rows) {
		user := new(Users)
		_ = rows.Scan(&user.Name, &user.Age, &user.Phone)
		data = append(data, user)
	})
	echoErr(cerr)
	for k, v := range data {
		fmt.Println(k, v.Name, v.Age, v.Phone)
	}
}

func selectTwo() {
	data := make([]*Users, 0)
	cerr = corm.GetDb(MasterDB).Tab("users").SelectRaw("name, age").Select("phone").Where("age", ">=", 24).Get(func(rows *sql.Rows) {
		user := new(Users)
		_ = rows.Scan(&user.Name, &user.Age, &user.Phone)
		data = append(data, user)
	})
	echoErr(cerr)
	for k, v := range data {
		fmt.Println(k, v.Name, v.Age, v.Phone)
	}
}

func whereOne() {
	data := make([]*Users, 0)
	cerr = corm.GetDb(MasterDB).Tab("users").Select("name", "age", "phone").
		Where("age", ">=", 24).
		Where("phone", "=", "18310953333").
		WhereRaw("id = 9").
		WhereLike("name", "王五").
		Get(func(rows *sql.Rows) {
			user := new(Users)
			_ = rows.Scan(&user.Name, &user.Age, &user.Phone)
			data = append(data, user)
		})
	echoErr(cerr)
	for k, v := range data {
		fmt.Println(k, v.Name, v.Age, v.Phone)
	}
}

func whereTwo() {
	data := make([]*Users, 0)
	cerr = corm.GetDb(MasterDB).Tab("users").Select("name", "age", "phone").
		Where("name", "=", "王五").
		WhereBetween("age", 30, 34).
		WhereBetween("created_at", "2017-08-08 00:00:00", "2019-11-14 23:23:00").
		WhereIntToStr("phone", "=", 18310953333).
		Get(func(rows *sql.Rows) {
			user := new(Users)
			_ = rows.Scan(&user.Name, &user.Age, &user.Phone)
			data = append(data, user)
		})
	echoErr(cerr)
	for k, v := range data {
		fmt.Println(k, v.Name, v.Age, v.Phone)
	}
}

func joinOne() {
	data := make([]*Users, 0)
	cerr = corm.GetDb(MasterDB).Tab("users").Join("user_groups", "users.id = user_groups.user_id").
		Select("users.name", "users.age", "users.phone", "user_groups.group_id").
		Where("users.name", "=", "王五").
		WhereBetween("users.age", 30, 34).
		WhereBetween("users.created_at", "2017-08-08 00:00:00", "2019-11-14 23:23:00").
		WhereIntToStr("users.phone", "=", 18310953333).
		Get(func(rows *sql.Rows) {
			user := new(Users)
			_ = rows.Scan(&user.Name, &user.Age, &user.Phone, &user.GroupId)
			data = append(data, user)
		})
	echoErr(cerr)
	for k, v := range data {
		fmt.Println(k, v.Name, v.Age, v.Phone, v.GroupId)
	}
}

func joinTwo() {
	data := make([]*Users, 0)
	cerr = corm.GetDb(MasterDB).Tab("users u").Join("user_groups ug", "u.id = ug.user_id").
		Select("u.name", "u.age", "u.phone", "ug.group_id").
		Where("u.name", "=", "王五").
		WhereBetween("u.age", 30, 34).
		WhereBetween("u.created_at", "2017-08-08 00:00:00", "2019-11-14 23:23:00").
		WhereIntToStr("u.phone", "=", 18310953333).
		Get(func(rows *sql.Rows) {
			user := new(Users)
			_ = rows.Scan(&user.Name, &user.Age, &user.Phone, &user.GroupId)
			data = append(data, user)
		})
	echoErr(cerr)
	for k, v := range data {
		fmt.Println(k, v.Name, v.Age, v.Phone, v.GroupId)
	}
}

func first() {
	u := new(Users)
	cerr = corm.GetDb(MasterDB).Tab("users u").Join("user_groups ug", "u.id = ug.user_id").
		Select("u.name", "u.age", "u.phone", "ug.group_id").
		Where("u.name", "=", "王五").
		WhereBetween("u.age", 30, 34).
		WhereRaw("u.phone like '%18310953333%'").
		First(&u.Name, &u.Age, &u.Phone, &u.GroupId)
	echoErr(cerr)
	fmt.Println(u.Name, u.Age, u.Phone, u.GroupId)
}

func getPage() {
	//当前页数
	page := 1
	//每页显示记录数
	pageCount := 2
	//总记录数
	var total int64
	//数据
	data := make([]*Users, 0)

	total, cerr = corm.GetDb(MasterDB).Tab("users").Select("name", "age", "phone").
		GetPage(page, pageCount, func(rows *sql.Rows) {
			user := new(Users)
			_ = rows.Scan(&user.Name, &user.Age, &user.Phone)
			data = append(data, user)
		})
	echoErr(cerr)
	fmt.Println("总记录数：", total, "总页数：", total/int64(pageCount))
	for k, v := range data {
		fmt.Println(k, v.Name, v.Age, v.Phone)
	}
}

func count() {
	count, cerr := corm.GetDb(MasterDB).Tab("users").Where("age", ">", 20).Count()
	echoErr(cerr)
	fmt.Println("总数：", count)

	max, cerr := corm.GetDb(MasterDB).Tab("users").Where("age", ">", 20).Max("age")
	echoErr(cerr)
	fmt.Println("最大年龄：", max)

	min, cerr := corm.GetDb(MasterDB).Tab("users").Where("age", ">", 20).Min("age")
	echoErr(cerr)
	fmt.Println("最小年龄：", min)

	sum, cerr := corm.GetDb(MasterDB).Tab("users").Where("age", ">", 20).Sum("age")
	echoErr(cerr)
	fmt.Println("年龄之和：", sum)
}

func insert() {
	insertId, cerr := corm.GetDb(MasterDB).Tab("users").
		Insert(map[string]interface{}{
			"nickname": "夏雨荷",
			"name":     "夏雨荷",
			"phone":    1231231234,
			"age":      30,
		})
	echoErr(cerr)
	fmt.Println("插入ID:", insertId)
}

func update() {
	num, err := corm.GetDb(MasterDB).Tab("users").
		WhereIn("id", 9, 10).
		Update(map[string]interface{}{
			"age":      20,
		})
	echoErr(err)
	//更新行数
	fmt.Println("影响行数：", num)
}

func exists() {
	is, err := corm.GetDb(MasterDB).Tab("users").Where("id", "=", 19).Exists()
	echoErr(err)
	//更新行数
	fmt.Println("数据是否存在：", is)
}

func trans() {

	err := corm.GetDb(MasterDB).Transaction(func(dbTrans *corm.Db) error {

		user := new(Users)
		err := dbTrans.Tab("users").Select("id").Where("nickname", "=", "大张伟").First(&user.Id)
		if err != nil {
			return err
		}

		_, err = dbTrans.Tab("users").Where("id", "=", user.Id).Update(map[string]interface{}{
			"name": "张柏芝22",
			"age":  30,
		})
		if err != nil {
			return err
		}

		_, err = dbTrans.Tab("users").Where("nickname", "=", "张三").Update(map[string]interface{}{
			"age":  30,
			"name": "周星驰22",
		})
		if err != nil {
			return err
		}

		_, err = dbTrans.Tab("users").Where("nickname", "=", "李四").Update(map[string]interface{}{
			"age":  30,
			"name": "吴镇宇22",
		})
		if err != nil {
			return err
		}

		return nil
	})
	echoErr(err)

}

func value() {
	valStr, err := corm.GetDb(MasterDB).Tab("users").Where("id", "=", 10).ValueStr("name")
	echoErr(err)
	fmt.Println("字符串：", valStr)

	valInt, err := corm.GetDb(MasterDB).Tab("users").Where("id", "=", 10).ValueInt("age")
	echoErr(err)
	fmt.Println("int值：", valInt)

	valFloat, err := corm.GetDb(MasterDB).Tab("users").Where("id", "=", 10).ValueFloat("age")
	echoErr(err)
	fmt.Println("float值：", valFloat)
}

func echoErr(err error) {
	if err != nil {
		panic(err)
	}
}

```
