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

For our usecase, we will avoid creating multiple directories and stick to a single directory (`models/`). We will be having a user model
for both Ax and Px, and another model for appointments. For our usecase we will refference a list of appointments (one-to-many) for
Px while this will stay empty for Ax. We will not be doing any transactional handeling and as such wont use transactions.

After defining our schema we will need to have logic for registering users and logging in first

With the database logic in place, we have to setup API end points which can call these database functionalities
we will be needing logical endpoints for
1. Authorization (login, register) POST methods
2. Middleware , this can reduce lines from our database code where we don't have to check for user type
3. endpoints for prof would be -> (check available/all slots, GET), (edit slot using start time, POST)
4. endpoints for students -> (check slots for a prof via mail, GET), (book a slot for prof via mail, POST), (check own slots, GET)

We start by hosting a server
1. https://www.naukri.com/code360/library/gin-middleware

Gin apparently allows you to define incoming JSON  requests, which is great much like gRPC, though this has its own disadvantage.
For maintaining sessions, we will be using cookies and not JWT for simplicity as this app will only run while server is on. We won't
be handling server restart logic.

The idea here becomes to use sessions which will act as a context and enable handling different clients.
Gin provides handlers which can be used as callbacks when a route is accessed, the signature of these handlers are in `*gin.context`
This becomes really useful and helps maintain code quality and we wont have to pass the database around.

To use sessions, add a store in the middleware which the handlers can access using Default(), this session will be set in clinet
(browser or Postman) but otherwise for a custom client this should be extracted and sent back with future headers.

To set up the routing logic, we will use the `routers/` directory, we will not be versioning our API's for this toy project.

Next we will create seperate handlers for profs and the students in the `handlers/` directory. Further we should only expose those endpoints
which are necessary via the middleware.


Finally, we have to write test cases. While I was able to test these endpoints on Postman, I will try to use Go tests module for writing tests
I expect the header to be stored to use sessions properly while making requests. In Go, the tests directory has to end file names wiht
`_test.go` and the function signature for entry point is TestMain(m *testing.M) which is then called using `m.Run()`

This invokes all the functions having `t *testing.T` as their argument, inside this function we will be putting examples from our
Postman, the tricky part here would be to set the right cookies for the right request. To this end we will be storing each of them
`var A1Cookie, A2Cookie, P1Cookie string` which come into play right after the registration.
