simple struct of etcd

0、启动阶段

1）、raftNode在将库中的节点启动后，异步启动一个定时器，每隔一个周期，向Node tickc
    channel发送一个信号
2）、node节点启动后，for循环select 监听tickc channel.一旦收到一个信号，则调用
    raft的tick函数。敲黑板，划重点：注意，这个tick在不同的角色阶段，不一样。比如
    在follow阶段，becomeFollow函数将它置成tickElection函数。so，raft中的tick
    属性是一个func()

1、vote阶段

1）、将身份由follow转变成candidate
2）、创建MsgVote类型消息，发送给集群内的所有节点
3）、MsgVote的消息将会发送到raft的msgs channel
    r.msgs = append(r.msgs, m)
    
4）、在node.go里的node.run方法中，构建了Ready对象，这个对象里面包含了被赋值的msgs,
    并最终写到node.readyc这个channel中









/**
    常识备注：
**/

1、每一个类型都占用了内存的若干字节数目。
摘抄参考：https://studygolang.com/articles/04106
eg:
var s string
fmt.Println(unsafe.Sizeof(s))  // prints 8

var a [3]uint32
fmt.Println(unsafe.Sizeof(a)) // prints 12

其中，空结构体 struct{}是不占用的存储空间的
var x [1000000000]struct{}
fmt.Println(unsafe.Sizeof(x)) // prints 0

//仅仅消耗slice的头部空间
var x = make([]struct{}, 1000000000)
fmt.Println(unsafe.Sizeof(x)) // prints 12 in the playground

//两个空struct值的地址可能是相同的
var a, b struct{}
fmt.Println(&a == &b) // true

a := struct{}{} // not the zero value, a real new struct{} instance
b := struct{}{}
fmt.Println(a == b) // true

2、所以空struct可以作为信号而存在


3、一个有意思的地方，之前忘了。如果一个函数的返回类型是接口类型。那我可以返回接口
    实现类，也可以返回接口实现类的指针。因为接口实现类指针，也实现了接口的两个方法。

4、sync.Mutex Lock会阻塞协程的