### Log Collection

    This is a golang project to display log from /var/log

    To install and download golang go to  https://golang.org/

### Here are some of the features
    Display a specific log file in /var/log
    Show the last n events of a specified file in /var/log
    Basic text/keyword filtered of events

### To Build the server
    
    change directory to logcollection, to generate logcollection binary
        go build
       

    To run the server
        ./logcollection

    note: the server is hard-coded to port 8000
### Use log Collection
    Once the logcollection binary is running
    To show the results, a browser or curl can be used

    Display the system.log
    http://localhost:8000/logs/files/system.log
    curl http://localhost:8000/logs/files/system.log
    
    Display the last 2 events from the system.log
    http://localhost:8000/logs/files/system.log/lastevents/2
    curl http://localhost:8000/logs/files/system.log/lastevents/2


    Display the last 50 events from the system.log with word apple 
    http://localhost:8000/logs/files/system.log/lastevents/50?filter=apple
    curl "http://localhost:8000/logs/files/system.log/lastevents/50?filter=apple"



### Run Unit Test

    change the directory to handlers
        go test -v 


Note:
   Giving the time constrain, I displayed the log file as is this will have performance impacts if the log file is large.  One possible implementation is to implement pagination or limit the number of log events return to the user. 