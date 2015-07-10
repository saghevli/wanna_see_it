
#Sample exercise for "Wanna See It?"#

##Local environment setup:##
#####For gin, an HTTP request framework:#####
* `go get github.com/gin-gonic/gin`

#####For gorp, used for ORM/DB interaction:#####
* `go get github.com/coopernurse/gorp`
* `go get github.com/mattn/go-sqlite3`

##Usage:##

#####Create user:#####
`URL: /login`

`$ curl --request POST 'http://localhost:8000/login' --data-urlencode User_id=saghevli --data-urlencode Pwd=hunter2`

#####Create a post:#####
`URL: /post Params: Text (string), Author (int), Img (string)`

`$ curl --request POST 'http://localhost:8000/postimg'  --data-urlencode Text=hello --data-urlencode Author=1234  --data-urlencode Img=picture`

#####Consume posts:#####
`URL: /posts/num_posts/post_offset`

`$ curl --request GET 'http://localhost:8000/posts/5/0'`

#####Request permalink:#####
`URL: /perma/image_id` (returned from a post call)

`$ curl --request GET 'http://localhost:8000/perma/image_id'`
