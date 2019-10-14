# Dice
dice notation processor

## Building

```bash
cd dice
go build .
```

## Interactive

```bash
$ cd dice
$ go run ./main
Dice Roller Shell
---------------------
-> 3d6+2
<- 18 : ((3 d 6 -> 16 [6,6,4]) + 2 -> 18)
-> d6
<- ERROR failed to parse d6
-> 1d6
<- 1 : (1 d 6 -> 1 [1])
-> 3 * ( 1 d 100 / 2)
<- 120 : (3 * ((1 d 100 -> 80 [80]) / 2 -> 40) -> 120)
-> 3b(4d6)
<- 13 : (3 b (4 d 6 -> 15 [2,3,4,6]) -> 13 [3,4,6])
-> 2w(4d6)
<- 5 : (2 w (4 d 6 -> 13 [2,3,4,4]) -> 5 [2,3])
-> exit
$
```

## Code

```
test := `3d6`
roller := dice.NewRoller()
result, plan, err := roller.Roll(test)
if err != nil {
  log.Errorf("%s : %v", test, err)
}else{
  log.Printf("%s : %d %d", test, result, plan)
}
```


## Remaining Work

* implement multi character operators (ie dF - a 3 sided die which yields [-1, 0, +1])
* support default operands ie: d6 == 1d6

