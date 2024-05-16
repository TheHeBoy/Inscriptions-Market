# 添加功能模块的操作顺序
1. 创建模型 `go run main.go make model [category]`
2. 创建迁移 `go run main.go make migration [add_categories_table] [StructName]`
3. 执行迁移 `go run main.go migrate up`

# 用于.sol 文件生成 go代码
```shell
    solc --abi Store.sol
    solc --bin Store.sol
    abigen --bin=Store_sol_Store.bin --abi=Store_sol_Store.abi --pkg=store --out=Store.go
```
