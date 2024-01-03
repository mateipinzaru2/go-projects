# CHANNELS
Channels are a typed, thread-safe queue. Channels allow different goroutines to communicate with each other.

## CREATE A CHANNEL
Like maps and slices, channels must be created *before* use. They also use the same `make` keyword:

```go
ch := make(chan int)
```

## SEND DATA TO A CHANNEL
```go
ch <- 69
```

The `<-` operator is called the *channel operator*. Data flows in the direction of the arrow. This operation will *block* until another goroutine is ready to receive the value.

## RECEIVE DATA FROM A CHANNEL
```go
v := <-ch
```

This reads and removes a value from the channel and saves it into the variable `v`. This operation will *block* until there is a value in the channel to be read.

## BLOCKING AND DEADLOCKS
A [deadlock](https://yourbasic.org/golang/detect-deadlock/#:~:text=yourbasic.org%2Fgolang,look%20at%20this%20simple%20example.) is when a group of goroutines are all blocking so none of them can continue. This is a common bug that you need to watch out for in concurrent programming.

## ASSIGNMENT
Run the program. You'll see that it deadlocks and never exits. The `filterOldEmails` function is trying to send on a channel, but no other goroutines are running that can accept the value *from* the channel.

Fix the deadlock by spawning a goroutine to send the "is old" values.