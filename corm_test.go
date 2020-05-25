package corm

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"time"
)

var masterDB *sql.DB
var err error

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
	masterDB, err = sql.Open("mysql", "root:g2q3g5p8@tcp(127.0.0.1:3306)/corm_demo")
	retErr(err)
}

func TestFilterZero(t *testing.T) {
	fmt.Println("-------------------查询零值-------------------")
	data := make([]*Users, 0)
	err = GetDb(masterDB).Tab("users").Select("name", "phone").Where("phone", "=", "").Get(func(rows *sql.Rows) {
		user := new(Users)
		_ = rows.Scan(&user.Name, &user.Phone)
		data = append(data, user)
	})
	retErr(err)
	for k, v := range data {
		fmt.Println(k, v.Name, v.Phone)
	}

	fmt.Println("-------------------过滤零值-------------------")
	data = make([]*Users, 0)
	err = GetDb(masterDB).Tab("users").Select("name", "phone").WhereFZ("phone", "=", "").Get(func(rows *sql.Rows) {
		user := new(Users)
		_ = rows.Scan(&user.Name, &user.Phone)
		data = append(data, user)
	})
	retErr(err)
	for k, v := range data {
		fmt.Println(k, v.Name, v.Phone)
	}
}

func TestSelect(t *testing.T) {
	fmt.Println("-------------------查询一条数据-------------------")
	name, phone := "", ""
	err = GetDb(masterDB).Tab("users").Select("name", "phone").Where("phone", "=", "13888888888").First(&name, &phone)
	fmt.Println(name, phone)

	fmt.Println("-------------------查询多条数据-------------------")
	data := make([]*Users, 0)
	err = GetDb(masterDB).Tab("users").SelectRaw("name, age").Select("phone").Where("age", ">=", 24).Get(func(rows *sql.Rows) {
		user := new(Users)
		_ = rows.Scan(&user.Name, &user.Age, &user.Phone)
		data = append(data, user)
	})
	retErr(err)
	for k, v := range data {
		fmt.Println(k, v.Name, v.Age, v.Phone)
	}
}

func TestWhere(t *testing.T) {
	fmt.Println("-------------------where多条件-------------------")
	data := make([]*Users, 0)
	err = GetDb(masterDB).Tab("users").Select("name", "age", "phone").
		Where("age", ">=", 24).
		Where("phone", "=", "18310953333").
		WhereRaw("id = 9").
		WhereLike("name", "王五").
		Get(func(rows *sql.Rows) {
			user := new(Users)
			_ = rows.Scan(&user.Name, &user.Age, &user.Phone)
			data = append(data, user)
		})
	retErr(err)
	for k, v := range data {
		fmt.Println(k, v.Name, v.Age, v.Phone)
	}
	fmt.Println("-------------------whereIn-------------------")
	data = make([]*Users, 0)
	err = GetDb(masterDB).Tab("users").Select("name", "age", "phone").
		WhereIn("phone", "18500082222", "13683624444").
		Get(func(rows *sql.Rows) {
			user := new(Users)
			_ = rows.Scan(&user.Name, &user.Age, &user.Phone)
			data = append(data, user)
		})
	retErr(err)
	for k, v := range data {
		fmt.Println(k, v.Name, v.Age, v.Phone)
	}
}

func TestJoin(t *testing.T) {
	fmt.Println("-------------------join-------------------")
	data := make([]*Users, 0)
	err = GetDb(masterDB).Tab("users u").Join("user_groups ug", "u.id = ug.user_id").
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
	retErr(err)
	for k, v := range data {
		fmt.Println(k, v.Name, v.Age, v.Phone, v.GroupId)
	}
	fmt.Println("-------------------join 原生语句-------------------")
	args := []interface{}{"users.num"}
	rows, err := masterDB.Query("SELECT users.name,users.age FROM users  INNER JOIN user_groups ON users.id = user_groups.user_id WHERE users.age > ?", args...)
	retErr(err)
	for rows.Next() {
		user := new(Users)
		_ = rows.Scan(&user.Name, &user.Age)
		data = append(data, user)
	}

	for k, v := range data {
		fmt.Println(k, v.Name, v.Age)
	}
}

func TestSelectPage(t *testing.T) {
	fmt.Println("------------------- 分页查询 -------------------")
	//当前页数
	page := 4
	//每页显示记录数
	pageCount := 2
	//总记录数
	var total int64
	//数据
	data := make([]*Users, 0)
	total, err = GetDb(masterDB).Tab("users").Select("name", "age", "phone").OrderBy("id", "desc").
		GetPage(page, pageCount, func(rows *sql.Rows) {
			user := new(Users)
			_ = rows.Scan(&user.Name, &user.Age, &user.Phone)
			data = append(data, user)
		})
	retErr(err)
	fmt.Println("总记录数：", total, "总页数：", total/int64(pageCount))
	for k, v := range data {
		fmt.Println(k, v.Name, v.Age, v.Phone)
	}
}

func TestCount(t *testing.T) {
	fmt.Println("------------------- Count -------------------")
	count, err := GetDb(masterDB).Tab("users").Where("age", ">", 20).Count()
	retErr(err)
	fmt.Println("总数：", count)
	fmt.Println("------------------- Max -------------------")
	max, err := GetDb(masterDB).Tab("users").Where("age", ">", 20).Max("age")
	retErr(err)
	fmt.Println("最大年龄：", max)
	fmt.Println("------------------- Min -------------------")
	min, err := GetDb(masterDB).Tab("users").Where("age", ">", 20).Min("age")
	retErr(err)
	fmt.Println("最小年龄：", min)
	fmt.Println("------------------- Sum -------------------")
	sum, err := GetDb(masterDB).Tab("users").Where("age", ">", 20).Sum("age")
	retErr(err)
	fmt.Println("年龄之和：", sum)
}

func TestInsert(t *testing.T) {
	fmt.Println("------------------- 插入数据 -------------------")
	insertId, err := GetDb(masterDB).Tab("users").
		Insert(map[string]interface{}{
			"nickname": "夏雨荷",
			"name":     "夏雨荷",
			"phone":    1231231234,
			"age":      30,
		})
	retErr(err)
	fmt.Println("插入ID:", insertId)
}

func TestUpdate(t *testing.T) {
	fmt.Println("------------------- 更新数据 -------------------")
	num, err := GetDb(masterDB).Tab("users").
		WhereIn("id", 9, 10).
		Update(map[string]interface{}{
			"age": 20,
		})
	retErr(err)
	//更新行数
	fmt.Println("影响行数：", num)
}

func TestExists(t *testing.T) {
	fmt.Println("------------------- 判断数据是否存在 -------------------")
	is, err := GetDb(masterDB).Tab("users").Where("id", "=", 19).Exists()
	retErr(err)
	//更新行数
	fmt.Println("数据是否存在：", is)
}

func TestTrans(t *testing.T) {
	fmt.Println("------------------- 事务 -------------------")
	err := GetDb(masterDB).Transaction(func(dbTrans *Db) error {
		user := new(Users)
		err := dbTrans.Tab("users").Select("id").WhereEqual("nickname", "大张伟").First(&user.Id)
		if err != nil {
			return err
		}

		_, err = dbTrans.Tab("users").Insert(map[string]interface{}{
			"nickname": "夏雨荷",
			"name":     "夏雨荷",
			"phone":    1231231234,
			"age":      30,
		})
		if err != nil {
			return err
		}

		_, err = dbTrans.Tab("users").WhereEqual("id", user.Id).Update(map[string]interface{}{
			"name": "666666",
			"age":  30,
		})
		if err != nil {
			return err
		}

		_, err = dbTrans.Tab("users").WhereEqual("nickname", "夏雨荷").Update(map[string]interface{}{
			"age":  30,
			"name": "5555555",
		})
		if err != nil {
			return err
		}
		return nil
	})
	retErr(err)
}

func TestForce(t *testing.T) {
	fmt.Println("------------------- 强制索引 -------------------")
	name, phone := "", ""
	err = GetDb(masterDB).Tab("users").Select("name", "phone").Where("phone", "=", "13888888888").First(&name, &phone)
	retErr(err)
	fmt.Println(name, phone)
}

func TestWhereLike(t *testing.T) {
	fmt.Println("------------------- 模糊查询 -------------------")
	data := make([]*Users, 0)
	err = GetDb(masterDB).Tab("users").Select("name", "age", "phone").
		WhereLikeLeft("phone", "136").
		Get(func(rows *sql.Rows) {
			user := new(Users)
			_ = rows.Scan(&user.Name, &user.Age, &user.Phone)
			data = append(data, user)
		})
	retErr(err)
	for k, v := range data {
		fmt.Println(k, v.Name, v.Age, v.Phone)
	}
}

func TestValue(t *testing.T) {
	fmt.Println("------------------- 获取字段值ValueStr -------------------")
	valStr, err := GetDb(masterDB).Tab("users").Where("id", "=", 10).ValueStr("name")
	retErr(err)
	fmt.Println("字符串：", valStr)

	fmt.Println("------------------- 获取字段值ValueInt -------------------")
	valInt, err := GetDb(masterDB).Tab("users").Where("id", "=", 10).ValueInt("age")
	retErr(err)
	fmt.Println("int值：", valInt)

	fmt.Println("------------------- 获取字段值ValueFloat -------------------")
	valFloat, err := GetDb(masterDB).Tab("users").Where("id", "=", 10).ValueFloat("age")
	retErr(err)
	fmt.Println("float值：", valFloat)
}

func retErr(err error) {
	if err != nil {
		panic(err)
	}
}
