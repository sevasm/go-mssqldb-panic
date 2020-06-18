`go-mssqldb-panic` can be used to trigger a panic in `mssql.Conn#checkBadConn(error)` by calling `sql.DB#BeginTx(context.Context, *sql.TxOptions)` after obtaining a valid connection, and the connection subsequently being closed by shutting down SQL server as per issue https://github.com/denisenkom/go-mssqldb/issues/536.

To reproduce the issue, launch a SQL server Docker container:

`docker run --name sqlserver -e "ACCEPT_EULA=Y" -e "MSSQL_SA_PASSWORD=Password1" -e "MSSQL_COLLATION=Latin1_General_CI_AS" -p 1433:1433 -d mcr.microsoft.com/mssql/server:latest`

Run the program:

`go run main.go`

Once the program starts printing output, stop the SQLServer Docker container:

`docker stop sqlserver`

The program will exit with the following output:

```
2020/06/18 17:05:22 initiating response reading
2020/06/18 17:05:22 ERROR: BeginRead failed EOF
2020/06/18 17:05:22 response finished
panic: driver.ErrBadConn in checkBadConn. This should not happen.

goroutine 1 [running]:
github.com/denisenkom/go-mssqldb.(*Conn).checkBadConn(0xc00032c000, 0x6f86e0, 0xc000040210, 0x6f86e0, 0xc000040210)
        C:/Users/sevasm/go/pkg/mod/github.com/denisenkom/go-mssqldb@v0.0.0-20200428022330-06a60b6afbbc/mssql.go:184 +0x1d1
github.com/denisenkom/go-mssqldb.(*Conn).begin(0xc00032c000, 0x6fc520, 0xc00000a038, 0x40a600, 0x663d80, 0x6940c0, 0xc00032e001, 0x26d60028)
        C:/Users/sevasm/go/pkg/mod/github.com/denisenkom/go-mssqldb@v0.0.0-20200428022330-06a60b6afbbc/mssql.go:288 +0xc2
github.com/denisenkom/go-mssqldb.(*Conn).BeginTx(0xc00032c000, 0x6fc520, 0xc00000a038, 0x0, 0xc00032c000, 0x1, 0x0, 0x0, 0x0)
        C:/Users/sevasm/go/pkg/mod/github.com/denisenkom/go-mssqldb@v0.0.0-20200428022330-06a60b6afbbc/mssql.go:928 +0x85
database/sql.ctxDriverBegin(0x6fc520, 0xc00000a038, 0x0, 0x6fbb60, 0xc00032c000, 0x0, 0x0, 0x0, 0xc000107d10)
        c:/go/src/database/sql/ctxutil.go:104 +0xa5
database/sql.(*DB).beginDC.func1()
        c:/go/src/database/sql/sql.go:1705 +0x72
database/sql.withLock(0x6fa560, 0xc00032e000, 0xc000107d10)
        c:/go/src/database/sql/sql.go:3184 +0x70
database/sql.(*DB).beginDC(0xc000126000, 0x6fc520, 0xc00000a038, 0xc00032e000, 0xc000041220, 0x0, 0x0, 0x0, 0x0)
        c:/go/src/database/sql/sql.go:1704 +0xb0
database/sql.(*DB).begin(0xc000126000, 0x6fc520, 0xc00000a038, 0x0, 0xc000107e01, 0x4519e8, 0x6b3cc0, 0xc00004e4b0)
        c:/go/src/database/sql/sql.go:1698 +0xda
database/sql.(*DB).BeginTx(0xc000126000, 0x6fc520, 0xc00000a038, 0x0, 0x2, 0x5e, 0x0)
        c:/go/src/database/sql/sql.go:1676 +0x90
main.loop(0xc000126000, 0x9, 0xc0000181e0)
        C:/Workspace-VisualStudioCode/go-mssqldb-panic/main.go:32 +0x154
main.main()
        C:/Workspace-VisualStudioCode/go-mssqldb-panic/main.go:23 +0x174
exit status 2
```
