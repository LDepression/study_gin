package main

import (
	"database/sql/driver"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

//初始化数据库
func initDB() (err error) {
	dsn := "root:zxz123456@tcp(127.0.0.1:3306)/mysql_demo?charset=utf8mb4&parseTime=True"
	// 也可以使用MustConnect连接不成功就panic

	//这里内部就已经调用了ping
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10) //闲置连接数
	return
}

type user struct {
	ID   int    `db:"id"`
	Age  int    `db:"age"`
	Name string `db:"name"`
}

func (u user) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Age}, nil
}

//get是查询单条数据
func queryRowDemo() {
	sqlStr := "select id,name,age from user where id=?"
	var u user
	err := db.Get(&u, sqlStr, 1)
	if err != nil {
		fmt.Printf("get failed,err:%v", err)
		return
	}
	fmt.Printf("id:%v,name:%v,age:%v", u.ID, u.Name, u.Age)
}

//select是查询多条数据
func queryMoreDemo() {
	sqlStr := "select id,name,age from user ser where id> ?"
	var users []user
	db.Select(&users, sqlStr, 0)
	//for _, v := range u {
	//	fmt.Println(v)
	//}

	fmt.Printf("%#v", users)
}

func insertRowDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	ret, err := db.Exec(sqlStr, "沙河小王子", 19)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theID)
}

// 更新数据
func updateRowDemo() {
	sqlStr := "update user set age=? where id = ?"
	ret, err := db.Exec(sqlStr, 39, 6)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("update success, affected rows:%d\n", n)
}

// 删除数据
func deleteRowDemo() {
	sqlStr := "delete from user where id = ?"
	ret, err := db.Exec(sqlStr, 6)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete success, affected rows:%d\n", n)
}

//nameExec用来绑定SQL语句与结构体或map中的同名字段
func insertUserDemo() error {
	_, _ = db.NamedExec(`insert into user(name,age) values(:name,:age)`,
		map[string]interface{}{
			"name": "七米",
			"age":  28,
		})
	return nil
}

func namedQuery() {
	sqlStr := "select * from user where name = :name"
	//使用map做命名查询
	rows, err := db.NamedQuery(sqlStr, map[string]interface{}{"name": "李四"})
	if err != nil {
		fmt.Printf("query failed,err:%v", err)
		return
	}
	for rows.Next() {
		var u user
		//这里不能是scan而是应该是使用结构体的查询
		rows.StructScan(&u)
		fmt.Println(u)
	}
}

// BatchInsertUsers2 使用sqlx.In帮我们拼接语句和参数, 注意传入的参数是[]interface{}
func BatchInsertUsers2(users []interface{}) error {
	query, args, _ := sqlx.In(
		"INSERT INTO user (name, age) VALUES (?), (?), (?)",
		users..., // 如果arg实现了 driver.Valuer, sqlx.In 会通过调用 Value()来展开它
	)
	fmt.Println(query) // 查看生成的querystring
	fmt.Println(args)  // 查看生成的args
	_, err := db.Exec(query, args...)
	return err
}

// BatchInsertUsers3 使用NamedExec实现批量插入
func BatchInsertUsers3(users []*user) error {
	_, err := db.NamedExec("INSERT INTO user (name, age) VALUES (:name, :age)", users)
	return err
}

// QueryByIDs 根据给定ID查询
func QueryByIDs(ids []int) (users []user, err error) {
	// 动态填充id
	query, args, err := sqlx.In("SELECT name, age FROM user WHERE id IN (?)", ids)
	if err != nil {
		return
	}
	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)
	err = db.Select(&users, query, args...)
	return
}
func main() {
	if err := initDB(); err != nil {
		fmt.Printf("init Db failed, err:%v", err)
		return
	}
	fmt.Println("sqlx connect successfully")
	//queryRowDemo()
	//queryMoreDemo()
	//insertUserDemo()
	//namedQuery()
	user1 := user{
		Age:  18,
		Name: "ll",
	}
	user2 := user{
		Age:  28,
		Name: "lll",
	}
	user3 := user{
		Age:  38,
		Name: "llll",
	}
	users := []interface{}{user1, user2, user3}
	BatchInsertUsers2(users)
	//var users1 []*user
	//users1 = append(users1, &user1, &user2, &user3)
	//BatchInsertUsers3(users1)

	//var idxs []int
	//idxs = append(idxs, 3, 2, 1)
	//users, _ := QueryByIDs(idxs)
	//fmt.Println(users)
}

//单行查询是get,多行查询时select
