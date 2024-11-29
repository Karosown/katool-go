# Katool - Go

## Stream

支持像Java一样的Stream流，但是由于Go不支持方法泛型，在使用过程中需要自行处理

包含map、reduce、filter、groupBy、distinct、sort、flatMap、orderBy等操作, 支持foreach遍历，也支持collect进行自定义函数逻辑采集

具体使用查看test

todo: 异步流 parallelStream

## list.Partition
进行分片操作，可以在分片后转换为stream流，同时也支持对分片进行foreach遍历（内部采用协程、同时可以控制协程大小）

## lockSupport
类似Java LockSupport 对协程进行控制，但是目前不支持自动回复协程

## remote
二次封装的resty请求库，适用于工作流处理，例如请求链处理

## convert
一些数据转换的util

