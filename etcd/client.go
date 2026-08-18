package etcd

import (
	"context"
	"fmt"
	"time"
	"user_service/config"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var Client *clientv3.Client

func InitEtcd() error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.Cfg.Etcd.Host},
		DialTimeout: 5 * time.Second, // etcd连接超时
	})
	if err != nil {
		return err
	}

	Client = cli
	fmt.Println("✅ Etcd 启动成功")
	return nil
}

// 服务注册
func RegisterService(serviceName, addr string) error {
	// 判断 etcd 客户端是否已经初始化
	// 如果没有先调用 InitEtcd()，这里就不能继续注册服务
	if Client == nil {
		return fmt.Errorf("etcd client not initialized")
	}

	// 创建一个上下文对象
	// context.Background() 表示这里没有设置超时，也没有主动取消
	ctx := context.Background()

	// 向 etcd 申请一个租约，过期时间是 10 秒
	// 如果 10 秒内没有续约，etcd 会自动删除绑定这个租约的 key
	lease, err := Client.Grant(ctx, 10)
	if err != nil {
		return err
	}

	// 拼接服务注册用的 key
	// 比如：
	// serviceName = "user-service"
	// addr = "127.0.0.1:8081"
	// key = "/services/user-service/127.0.0.1:8081"
	key := fmt.Sprintf("/services/%s/%s", serviceName, addr)

	// 把服务地址写入 etcd
	// key:   /services/user-service/127.0.0.1:8081
	// value: 127.0.0.1:8081
	//
	// clientv3.WithLease(lease.ID) 表示这个 key 绑定到上面的租约
	// 只要租约过期，这个 key 就会被 etcd 自动删除
	_, err = Client.Put(ctx, key, addr, clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}

	// 开启自动续约
	// etcd 会返回一个 channel，只要程序还活着，就会不断收到续约响应
	keepAliveChan, err := Client.KeepAlive(ctx, lease.ID)
	if err != nil {
		return err
	}

	fmt.Printf("✅ %s 服务注册成功：%s\n", serviceName, addr)

	// 启动一个 goroutine 持续读取 keepAliveChan
	// 必须读取这个 channel，否则自动续约可能无法正常持续
	go func() {
		for keepAliveResp := range keepAliveChan {
			if keepAliveResp == nil {
				fmt.Println("etcd 续约失败")
				return
			}

			fmt.Println("etcd 续约成功")
		}
	}()

	return nil
}

// 获取所有服务
func GetService(serviceName string) ([]string, error) {
	// 判断 etcd 客户端是否已经初始化
	// 如果没有先调用 InitEtcd()，就不能从 etcd 查询服务
	if Client == nil {
		return nil, fmt.Errorf("etcd client not initialized")
	}

	// 创建上下文对象
	// 这里暂时不用管它，可以理解成本次请求的控制对象
	ctx := context.Background()

	// 拼接查询前缀
	// 比如 serviceName = "user-service"
	// prefix = "/services/user-service/"
	//
	// etcd 里面可能有：
	// /services/user-service/127.0.0.1:8081
	// /services/user-service/127.0.0.1:8082
	// /services/user-service/127.0.0.1:8083
	prefix := fmt.Sprintf("/services/%s/", serviceName)

	// 按前缀查询 etcd
	// clientv3.WithPrefix() 的意思是：
	// 不是只查某一个完整 key，而是查询所有以 prefix 开头的 key
	resp, err := Client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	// 创建一个字符串切片，用来保存查询到的服务地址
	// len(resp.Kvs) 是查询到的 key-value 数量
	addrs := make([]string, 0, len(resp.Kvs))

	// 遍历 etcd 返回的所有 key-value
	for _, kv := range resp.Kvs {
		// kv.Value 是 etcd 中保存的 value
		// 注册服务时 value 存的是服务地址，比如 "127.0.0.1:8081"
		// kv.Value 类型是 []byte，所以需要转成 string
		addrs = append(addrs, string(kv.Value))
	}

	// 返回所有查询到的服务地址
	return addrs, nil
}

func CloseEtcd() error {
	if Client == nil {
		return nil
	}
	return Client.Close()
}
