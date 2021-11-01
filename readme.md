## marin-blc 智能合约
marin-blc 是一种允许在没有第三方的情况下进行可信交易，并且这些交易可追踪不可逆转的合约。

## 如何使用
1. 设置环境变量
> set NODE_ID=3000

![](https://github.com/marin-man/marin-blc/blob/master/img/1.png)

2. 创建钱包
> bc.exe createwallet

![](https://github.com/marin-man/marin-blc/blob/master/img/2.png)

3. 查看钱包集合
> bc.exe accounts

![](https://github.com/marin-man/marin-blc/blob/master/img/3.png)

4. 创建区块链
> bc.exe createblockchain -address 选择一个钱包地址

![](https://github.com/marin-man/marin-blc/blob/master/img/4.png)

5. 查看当前钱包余额
> bc.exe getbalance -address 选择一个钱包地址

![](https://github.com/marin-man/marin-blc/blob/master/img/5.png)

6. 再创建一个钱包
> bc.exe createwallet

![](https://github.com/marin-man/marin-blc/blob/master/img/6.png)

7. 转账操作
> bc.exe send -from "[\\"付款地址\\"]" -to "[\\"收款地址\\"]" -ammount "[\\"金额\\"]"

![](https://github.com/marin-man/marin-blc/blob/master/img/7.png)

8. 启动服务（将此服务上线）
> bc.exe start

![](https://github.com/marin-man/marin-blc/blob/master/img/8.png)

9. 开启另一个服务（会从主节点中同步数据）
> bc.exe set NODE_ID=3001
> bc.exe start

![](https://github.com/marin-man/marin-blc/blob/master/img/9.png)
![](https://github.com/marin-man/marin-blc/blob/master/img/10.png)
