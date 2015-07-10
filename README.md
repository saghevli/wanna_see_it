
#Sample exercise for "Wanna See It?"#

Simple API for "Wanna See It?". Users are presented with a text description of a photo. If they choose, they can view it (once). This API currently supports:

* User login
* Post creation
* Post consumption
* Permalinks

##Notes##

User login is currently mostly mocked. A full sytstem would require OAUTH or other fully fledged authentication scheme, without personally implemented crypto. Images are posted at /postimg, with the specified parameters passed as urlencoded values. Images are consumed at the /posts endpoint, which is passed a count and an offset. This would enable a simple, descending order by date API, which would allow a user to scroll through a chronological feed of posts. Supporting other consumption patters, such as viewing posts by user, or searches for data contained in a posts text field would require additional database indexes on those fields (to operate at scale).

####Possible optimizations for bandwidth constraints:####

Since posts initially only expose their text fields, the consumption endpoint could receive everything but the image data, which would be replaced by a permalink. When and if users want to see the image, it would be loaded at the permalink location, replacing the need for every image to be transmitted to the mobile app.

For a real production system, I'd include the following that were outside the scope of this small example.

* Real images! (currently mocked out with strings)
* Real UUID's, collision checks, etc. Currently using a random, 63 bit integer.
* SSL/TLS sitewide. Any http requests should be redirected to HTTPS. 
* Testing plan: dependency inject a mock database, unit test functions for success and error conditions.
* OAUTH, no self-implemented crypto.

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


