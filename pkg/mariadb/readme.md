// 1. 最简单的方式
db, err := mariadb.NewSimple("localhost", "root", "password", "testdb")

// 2. 指定端口
db, err := mariadb.NewWithPort("localhost", 3306, "root", "password", "testdb")

// 3. 使用默认配置
db, err := mariadb.NewWithDefaults()

// 4. 从DSN字符串
db, err := mariadb.NewFromDSN("root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local")

// 5. 使用选项模式（最灵活）
db, err := mariadb.NewWithOptions(
	mariadb.WithHost("localhost"),
	mariadb.WithPort(3306),
	mariadb.WithUser("root"),
	mariadb.WithPassword("password"),
	mariadb.WithDatabase("testdb"),
	mariadb.WithPool(100, 20, 60*time.Minute),
)

// 6. 使用完整配置
db, err := mariadb.New(mariadb.Config{
	Host:     "localhost",
	Port:     3306,
	User:     "root",
	Password: "password",
	DBName:   "testdb",
	Params:   "charset=utf8mb4&parseTime=True&loc=Local",
	MaxOpen:  100,
	MaxIdle:  20,
	MaxLife:  60 * time.Minute,
})

// 7. 失败时panic（适合初始化时）
db := mariadb.MustNewSimple("localhost", "root", "password", "testdb")