# APPEND
The built-in append function is used to dynamically add elements to a slice:

`func append(slice []Type, elems ...Type) []Type`

If the underlying array is not large enough, `append()` will create a new underlying array and point the slice to it.

Notice that `append()` is variadic, the following are all valid:

```go
slice = append(slice, oneThing)
slice = append(slice, firstThing, secondThing)
slice = append(slice, anotherSlice...)
```

## ASSIGNMENT
We've been asked to "bucket" costs per day, in a given period.

Complete the `getCostsByDay` function. It should return a slice of `float64`s, where each element is the total cost for that day. The length of the slice should be equal to the number of days represented in the `costs` slice up to the last day represented in the slice. Only one "bucket" of costs per day, and include a "bucket" for days without any costs.

Days are numbered and start at `0`.

### EXAMPLE
Input in `day, cost` format:

```go
[]cost{
    {0, 4.0},
    {1, 2.1},
    {5, 2.5},
    {1, 3.1},
}
```

I would expect this result:

```go
[]float64{
    4.0,
    5.2,
    0.0,
    0.0,
    0.0,
    2.5,
}
```