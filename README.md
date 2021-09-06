<font size="4">

# 项目介绍
基于微服务架构开发电商系统和后台管理系统，使用Go语言完成微服务的架构和实现，系统包括完整的商品展示、购物车、订单、收货地址管理等。
# 技术栈使用
1. 数据库:GORM/MYSQ/Redis/ElasticSearch
2. 日志库:zap  
3. 配置中心:nacos/viper  
4. 服务发现: consul  
5. NoSQL:redis(短信发送/分布式锁)  
6. 搜索:elasticsearch  
7. 消息队列:rocketMq(商品扣减分布式事务实现)

# 微服务模块划分
1. 商品模块(goods)
2. 订单模块(order)
3. 用户模块(user)
4. 用户行为模块(userop)
5. 库存模块(inventory)