
#Sample exercise for "Wanna See It?"#

##Local environment setup:##
#####For gin, an HTTP request framework:#####
* `go get github.com/gin-gonic/gin`

#####For gorp, used for ORM/DB interaction:#####
* `go get github.com/coopernurse/gorp`
* `go get github.com/mattn/go-sqlite3`

##Usage:##

#####Create user:#####
`curl --request POST 'http://localhost:8000/login' --data-urlencode User_id=saghevli --data-urlencode Pwd=hunter2`

#####Create a post:#####
`curl --request POST 'http://localhost:8000/post'  --data-urlencode Text=hello --data-urlencode Author=1234  --data-urlencode Img=picture`

#####Consume posts:#####
`curl --request GET 'http://localhost:8000/posts/5'`

#####Request permalink:#####
`curl --request GET 'http://localhost:8000/perma/*id of an image you added*'`
