# SupInscriptions

基于以太坊 ETHS https://docs.ethscriptions.com/  协议实现的 inscriptions market。

- 铭文索引器：https://github.com/TheHeBoy/indexer_msc20
- 市场合约：https://github.com/TheHeBoy/market-contract
- 前端demo:：https://github.com/TheHeBoy/market-front

## 后端接口

| 接口名称                 | 请求类型 | 接口路径             |
| ------------------------ | -------- | -------------------- |
| 获取签名message          | **GET**  | /auth/message        |
| 签名登录                 | **POST** | /auth/login          |
| 刷新token                | **POST** | /auth/refresh-token  |
| tokens 分页              | **GET**  | /tokens/page         |
| 出售中的token分页        | **GET**  | /tokens/page-listing |
| 通过地址查询token        | **GET**  | /tokens/{address}    |
| 创建订单                 | **PUT**  | /orders/create       |
| 通过tick查询出售中的订单 | **GET**  | /orders/listing      |
| 查看个人订单             | **GET**  | /orders/{address}    |
| 订单签名                 | **POST** | /orders/sign         |
| 通过地址查询msc20        | **GET**  | /msc20/{address}     |
| 获取最新铭文             | **GET**  | /inscriptions/latest |

## 订单状态

订单的状态流程如下：

![无标题-2024-04-23-0246](https://s2.loli.net/2024/05/19/nlJNI2vw7AzRjqg.png)

## 功能截图

### 0. 最新铭文

![image-20240520000848909](https://s2.loli.net/2024/05/20/oEUby4a9AYhC8F1.png)

### 1. Deploy 铭文

![image-20240520000024625](https://s2.loli.net/2024/05/20/91r5Y2BSIFQCfzA.png)

### 2. Tokens 列表
![image-20240519235959262](https://s2.loli.net/2024/05/20/iJH2z4t7crguSYQ.png)

#### 3. 个人界面

#### 3.1 持有的 Tokens

![image-20240520000109880](https://s2.loli.net/2024/05/20/KRATPBZogvs2lIj.png)

##### 3.1.1 出售 Tokens
![image-20240520000408349](https://s2.loli.net/2024/05/20/GEOPC1YVjZgkX2s.png)

##### 3.1.2 发送 Tokens

![image-20240520000540330](https://s2.loli.net/2024/05/20/5aApHgsdNoYhO9Z.png)

#### 3.2 铭刻的铭文
![image-20240520000229338](https://s2.loli.net/2024/05/20/zNqVtZTYQWaRbKD.png)

#### 3.3 个人订单

![image-20240520000801520](https://s2.loli.net/2024/05/20/DHSujEPdfU6NvLb.png)

#### 3.4 Tokens 市场

![image-20240520001638078](https://s2.loli.net/2024/05/20/Il7RDJGzumg1Ls9.png)

##### 3.4.1 出售的订单

![image-20240520001736347](https://s2.loli.net/2024/05/20/SjkeyhTND4cuUHK.png)

## License
[MIT](https://github.com/inscriptions-market/market-end?tab=MIT-1-ov-file)