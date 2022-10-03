package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	rdb *redis.Client
)

var ctx = context.Background()

// 初始化连接
func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
		PoolSize: 100,      // 连接池大 小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	return err
}
func hGetDemo() {
	ctx := context.Background()
	if err := initClient(); err != nil {
		return
	}
	v, err := rdb.HGetAll(ctx, "user").Result()
	if err != nil {
		fmt.Println("hGetAll failed,err:", err)
		return
	}
	fmt.Println("user:", v)
	//rdb.HMGet()
}

func hsetVal() {
	ctx := context.Background()
	if err := initClient(); err != nil {
		return
	}
	rdb.HSet(ctx, "user", "Tom", "Jnny")
}

func redisDemo() {

	//最后面的参数是key存在的时间,如果为0的话表示永久的
	err := rdb.Set(ctx, "name", "lyc", time.Second*60).Err()
	if err != nil {
		fmt.Println("err:", err)
	} else {
		fmt.Println("数据库插入成功.......")
	}
}

//数据的查询
func redisQuery() {

	val, err := rdb.Get(ctx, "user").Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("key 值不存在")
		} else {
			panic(err)
		}
	} else {
		fmt.Println("val->", val)
	}
}

type Model struct {
	Str1 string `json:"str1"`
	Str2 string `json:"str2"`
}

//struct类型数据的查询
//为结构体添加方法,否则会报错
func (m *Model) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}
func (m *Model) UnMarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

//要注意的是这里需要给对象添加上MarshalBinary这个方法用于将struct类型转成 []byte 类型。
func InsertStructData() {
	err := rdb.Set(ctx, "k1", &Model{
		Str1: "lyc",
		Str2: "xxx",
	}, time.Minute).Err()
	if err != nil {
		fmt.Println("err:", err)
		panic(err)
	} else {
		fmt.Println("插入成功")
	}
}

//对于结构体数据的查询
func redisQueryStructData() {
	data, err := rdb.Do(ctx, "get", "user").Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("不存在这个key值")
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("data -> %v\n", data)
		var myModel Model
		json.Unmarshal([]byte(data.(string)), &myModel)
		fmt.Printf("data-> %v\n", myModel)
	}
}

//插入map类型数据的插入
func InsertMapData() {
	data := make(map[string]string)
	data["name"] = "lyc"
	data["age"] = "18"
	err := rdb.HSet(ctx, "testHash", data).Err()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("没有这个数据")
		} else {
			panic(err)
		}
	} else {
		fmt.Println("map类型数据插入成功........")
	}
}

func redisQueryMapData() {
	//查询整个map,返回的类型是map[string]
	data, err := rdb.HGetAll(ctx, "testHash").Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("没有查找到值")
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("data-->%v", data)
	}

	//去key中的field对应的值
	hashGet, _ := rdb.HGet(ctx, "testHash", "age").Result()
	fmt.Printf("hashGet-->value : %+v\n", hashGet)

}

//List数据的插入删除
func redisInsertListData() {
	//插入list数据
	//Lpush是插入到首部
	//rpush是插入到尾部
	rdb.LPush(ctx, "ListDemo", "首部").Err()
	rdb.RPush(ctx, "ListDemo", "尾部").Err()

	//去除并且一处左边第一个元素
	// 对应的有 RPop可以去除并且移除list中的最后一个元素
	first, err := rdb.LPop(ctx, "ListDemo").Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("未找到该数据值")
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("first: %v\n", first)
	}

	// 查询出list中指定位置的元素
	listIndexData, _ := rdb.LIndex(ctx, "ListDemo", 1).Result()
	fmt.Println("listData", listIndexData)

	// 查询出列表的长度
	listLen, _ := rdb.LLen(ctx, "ListDemo").Result()
	fmt.Println("listLen", listLen)

}

// (zset) 操作有序的set
func TestRestZset() {
	zSetKey := "zSetKey"
	people := []*redis.Z{
		&redis.Z{Score: 20, Member: "张三"},
		&redis.Z{Score: 50, Member: "老王"},
		&redis.Z{Score: 100, Member: "李四"},
		&redis.Z{Score: 20, Member: "小黄"},
	}

	// 想zset中插入数据
	rdb.ZAdd(ctx, zSetKey, people...)

	// 给指定的元素指定score
	//让张三的分数加5.0
	newScore, _ := rdb.ZIncrBy(ctx, zSetKey, 5.0, "张三").Result()
	fmt.Println("张三的newScore ->", newScore)

	// 查询zset中Score前两名的数据
	zsetListData, _ := rdb.ZRevRangeWithScores(ctx, zSetKey, 0, 1).Result()
	fmt.Println(zsetListData)

	// 移除分数最低的两个数据，返回的是成功移除的条数
	zsetListData2, _ := rdb.ZRemRangeByRank(ctx, zSetKey, 0, 0).Result()
	fmt.Println(zsetListData2)

}

func watchDemo() {
	// 监视watch_count的值，并在值不变的前提下将其值+1
	key := "watch_count"
	err := rdb.Watch(ctx, func(tx *redis.Tx) error {
		n, err := tx.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}
		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			//这里是用来加业务逻辑的
			time.Sleep(time.Second * 10)
			pipe.Set(ctx, key, n+1, 0)
			return nil
		})
		return err
	}, key)
	if err != nil {
		fmt.Println("tx exec failed,err:", err)
		return
	}
	fmt.Println("tx exec success")
}

//注意
//1. %v    只输出所有的值   &{ "18","lyc"}
//
//2. %+v 先输出字段名字，再输出该字段的值 &{age:"18",name:"lyc"}
//
//3. %#v 先输出结构体名字值，再输出结构体（字段名字+字段的值）&main.student{age:"18",name:"lyc"}
func main() {
	err := initClient()
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("redis connect successfully")
	//redisDemo()
	//hGetDemo()
	//redisQuery()
	//InsertStructData()
	//redisQueryStructData()
	//InsertMapData()
	//redisQueryMapData()
	//	redisInsertListData()
	watchDemo()
}
