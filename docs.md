Since this is the first time I will be working with HTTP framework in Go, i need to explore existing options. There is the ocassional
builtin net package but there exists better options

1. https://medium.com/@hasanshahjahan/a-deep-dive-into-gin-chi-and-mux-in-go-33b9ad861e1b
2. https://www.golang.company/blog/gin-vs-mux
3. https://www.reddit.com/r/golang/comments/1bg87hz/title_after_gorillamuxs_revival_stick_with_it_or/

We will be going with Gin, I found its syntax a little more simple, Mux would not be necessary for this project.

Next we need to see some standard project structure for Gin projects for best practices

1. https://golang.withcodeexample.com/blog/chapter-3-project-structure-in-gin-framework/
2. https://github.com/eddycjy/go-gin-example/tree/master
3. https://medium.com/@wahyubagus1910/build-scalable-restful-api-with-golang-gin-gonic-framework-43793c730d10

I would usually keep my `main.go` inside `cmd/` directory but since this dosnt seem usual, we will keep it in root.
We will be keeping the `middleware/` (for authentication) and  `handlers/` (for handling request logic for clients for Ax and Px)

As for database, after some digging, there exists a simple method to access databases called GORM, which can be used to interface with
databases, the GORM library has support for mysql/postgres/sqlite etc. We will be using sqlite since it dosn't need a password/server setup

1. https://www.pingcap.com/article/building-robust-go-applications-with-gorm-best-practices/#:~:text=One%20of%20the%20most%20popular,design%20and%20comprehensive%20feature%20set.
2. https://medium.com/@itskenzylimon/getting-started-on-golang-gorm-af49381caf3f