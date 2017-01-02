# workout-log
Workout Log built with golang and google appengine

# Download and install the appengine golang SDK
https://cloud.google.com/appengine/docs/go/download

# from workout-log directory

go get golang.org/x/net/context
go get google.golang.org/appengine
go get google.golang.org/appengine/datastore
go get google.golang.org/appengine/user
go get google.golang.org/appengine/log
go get google.golang.org/appengine/search

# and then view at: http://localhost:8080/
goapp serve .

# Todo
1. Build workout struct for datastore
2. Add task to index workouts for searching
3. Add mailgun API for outgoing/incoming email to queue up workouts and notifications
4. Import datastore into big query for analysis
5. add cloud endpoints for mobile app development
6. Add firebase.io for push and realtime features