gostudy
=======

Some people learning [Go](https://golang.org).

Setup
-----

Install Go:

```sh
brew install go
```

Make sure the Go bin directory is in your path (you'll probably want to put this in a shell startup script):

```sh
export PATH="$HOME/go/bin:$PATH"
```

Install some useful tools:

```sh
brew install httpie
go get golang.org/x/tools/cmd/goimports
go get golang.org/x/lint/golint
```

Install or enable Go support in your editor/IDE; I use https://github.com/fatih/vim-go

Clone this repo:

```sh
git clone https://github.com/cultureamp/gostudy
```

The `master` branch is empty, ready for you to proceed. You'll find various other branches with example code.

Session 01 — HTTP server & middleware
-------------------------------------

Initialize your codebase as a Go module: `go mod init gostudy`

This produces a `go.mod` that looks like this:

```
module gostudy

go 1.12
```

Create a basic HTTP server:

```diff
commit 29016c8b4efa07c5578a32dd468613f5ffbc565a
Author: Paul Annesley <paul@annesley.cc>
Date:   Thu Sep 5 14:39:17 2019 +1000

    hello world HTTP server

diff --git a/main.go b/main.go
new file mode 100644
index 0000000..a3eb9bf
--- /dev/null
+++ b/main.go
@@ -0,0 +1,21 @@
+package main
+
+import (
+	"io"
+	"net/http"
+	"os"
+)
+
+func main() {
+	port := os.Getenv("PORT")
+	if port == "" {
+		port = "1234"
+	}
+
+	app := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
+		w.WriteHeader(200)
+		io.WriteString(w, "hello world")
+	})
+
+	http.ListenAndServe(":"+port, app)
+}
```

Add a basic logging middleware:

```diff
commit 19874cb54d6c008835e709f56e998aaf5ca202d2
Author: Paul Annesley <paul@annesley.cc>
Date:   Thu Sep 5 14:39:57 2019 +1000

    httplog middleware

diff --git a/main.go b/main.go
index a3eb9bf..0a98801 100644
--- a/main.go
+++ b/main.go
@@ -2,6 +2,7 @@ package main
 
 import (
 	"io"
+	"log"
 	"net/http"
 	"os"
 )
@@ -17,5 +18,12 @@ func main() {
 		io.WriteString(w, "hello world")
 	})
 
-	http.ListenAndServe(":"+port, app)
+	http.ListenAndServe(":"+port, httplog(app))
+}
+
+func httplog(next http.Handler) http.Handler {
+	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
+		log.Println(r.Method, r.URL.Path)
+		next.ServeHTTP(w, r)
+	})
 }
```

Add an authentication middleware:

```diff
commit a3cf02a1c291b72d58ea82bca717e9f3197842cb
Author: Paul Annesley <paul@annesley.cc>
Date:   Thu Sep 5 14:40:26 2019 +1000

    secure middleware for authentication

diff --git a/main.go b/main.go
index 0a98801..221ec07 100644
--- a/main.go
+++ b/main.go
@@ -18,7 +18,7 @@ func main() {
 		io.WriteString(w, "hello world")
 	})
 
-	http.ListenAndServe(":"+port, httplog(app))
+	http.ListenAndServe(":"+port, httplog(secure(app)))
 }
 
 func httplog(next http.Handler) http.Handler {
@@ -27,3 +27,15 @@ func httplog(next http.Handler) http.Handler {
 		next.ServeHTTP(w, r)
 	})
 }
+
+// secure verifies that all requests have an extremely strong password
+func secure(next http.Handler) http.Handler {
+	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
+		// auth the request
+		if r.Header.Get("Authorization") != "hunter2" {
+			w.WriteHeader(401)
+			return
+		}
+		next.ServeHTTP(w, r)
+	})
+}
```

Session 02 — JSON request & response
------------------------------------

```diff
commit 9d023bce7b232df9942e6e6029e122ec57480ef3
Author: Paul Annesley <paul@annesley.cc>
Date:   Thu Sep 5 14:47:13 2019 +1000

    JSON request & response data

diff --git a/main.go b/main.go
index 221ec07..40bc586 100644
--- a/main.go
+++ b/main.go
@@ -1,12 +1,22 @@
 package main
 
 import (
-	"io"
+	"encoding/json"
+	"io/ioutil"
 	"log"
 	"net/http"
 	"os"
 )
 
+type RequestThing struct {
+	Name string
+	Camp string
+}
+
+type ResponseThing struct {
+	Greeting string
+}
+
 func main() {
 	port := os.Getenv("PORT")
 	if port == "" {
@@ -14,8 +24,25 @@ func main() {
 	}
 
 	app := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
+		data, err := ioutil.ReadAll(r.Body)
+		if err != nil {
+			w.WriteHeader(500)
+			return
+		}
+		r.Body.Close()
+
+		var rt RequestThing
+		json.Unmarshal(data, &rt)
+
+		respBytes, err := json.Marshal(ResponseThing{Greeting: "hello " + rt.Name})
+		if err != nil {
+			w.WriteHeader(500)
+			return
+		}
+
+		w.Header().Set("Content-Type", "application/json")
 		w.WriteHeader(200)
-		io.WriteString(w, "hello world")
+		w.Write(respBytes)
 	})
 
 	http.ListenAndServe(":"+port, httplog(secure(app)))
```
