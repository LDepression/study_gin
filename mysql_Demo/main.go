package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" //这里其实只需要(init)
	"time"
)

var db *sql.DB

type user struct {
	id   int
	age  int
	name string
}

func initMySql() error {
	dsn := "root:zxz123456@tcp(127.0.0.1:3306)/mysql_demo"
	//去初始化全局的db对象而不是声明一个db对象
	var err error
	db, err = sql.Open("mysql", dsn) //填写的信息的格式的校验
	if err != nil {
		fmt.Printf("connect to db failed,err:%v", err)
		panic(err)
	}
	//做完错误检查之后,确保db不为nil
	//尝试与数据库建立连接(校验dsn是否正确)
	err = db.Ping()
	if err != nil {
		fmt.Printf("connect to db failed,err:%v", err)
		return err
	}
	//defer db.Close() //释放掉连接相关的资源mysql_demo

	db.SetConnMaxLifetime(10 * time.Second) //
	db.SetMaxOpenConns(200)                 //最大连接数
	db.SetMaxIdleConns(10)                  //最大空闲连接数
	return nil
}

// 查询单条数据示例
func queryRowDemo() {
	sqlStr := "select id, name, age from user where id=?"
	var u user
	// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	//如果不去释放的话,在设置最大连接数的时候如果只设置1的话,程序就会在此处卡住
	err := db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
}

// 查询多条数据示例
func queryMultiRowDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	// 非常重要：关闭rows释放持有的数据库链接
	defer rows.Close()

	// 循环读取结果集中的数据
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
	}
}

// 插入数据
func insertRowDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	ret, err := db.Exec(sqlStr, "李四", 38)
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

//更新数据
func updateRowDemo() {
	sqlStr := "update user set age=? where id = ?"
	ret, err := db.Exec(sqlStr, 39, 2)
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

// 预处理查询示例
func prepareQueryDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	//先将命令行传过去
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err:%v\n", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	defer rows.Close()
	// 循环读取结果集中的数据
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return
		}
		fmt.Printf("id:%d name:%s age:%d\n", u.id, u.name, u.age)
	}
}

//解决sql注入的问题:我们在任何时候都不要去拼接一个sql语句
// sql注入示例
func sqlInjectDemo(name string) {
	sqlStr := fmt.Sprintf("select id, name, age from user where name='%s'", name)
	fmt.Printf("SQL:%s\n", sqlStr)
	var u user
	err := db.QueryRow(sqlStr).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("exec failed, err:%v\n", err)
		return
	}
	fmt.Printf("user:%#v\n", u)
}

// 事务操作示例
func transactionDemo() {
	tx, err := db.Begin() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback() // 回滚
		}
		fmt.Printf("begin trans failed, err:%v\n", err)
		return
	}
	sqlStr1 := "Update user set age=30 where id=?"
	ret1, err := tx.Exec(sqlStr1, 2)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql1 failed, err:%v\n", err)
		return
	}
	affRow1, err := ret1.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}

	sqlStr2 := "Update user set age=40 where id=?"
	ret2, err := tx.Exec(sqlStr2, 3)
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec sql2 failed, err:%v\n", err)
		return
	}
	affRow2, err := ret2.RowsAffected()
	if err != nil {
		tx.Rollback() // 回滚
		fmt.Printf("exec ret1.RowsAffected() failed, err:%v\n", err)
		return
	}

	fmt.Println(affRow1, affRow2)
	if affRow1 == 1 && affRow2 == 1 {
		fmt.Println("事务提交啦...")
		tx.Commit() // 提交事务
	} else {
		tx.Rollback()
		fmt.Println("事务回滚啦...")
	}

	fmt.Println("exec trans success!")
}
func main() {
	if err := initMySql(); err != nil {
		fmt.Printf("connect to db failed,err:%v", err)
	}
	fmt.Println("connect to db suc cess")
	//insertRowDemo()
	//queryRowDemo()
	//updateRowDemo()
	//prepareQueryDemo()
	sqlInjectDemo("xxx' or 1=1#")
}
