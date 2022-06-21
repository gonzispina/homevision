# Homevision

The implementation is written in golang, so you will need golang 1.8 to run it.


### Considerations
To solve the problem of the API responding non-200 status codes I've just
implemented a local retrier with a big treshold.

Ideally the program should persist the state to reprocess just the things that failed later, but that is beyond the scope of the test.

### How to run? 
```
make run
```

It will create a tmp folder in the working directory with all the downloaded images




