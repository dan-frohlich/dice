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
-> d6+2
<- 4 : ((1d6 [2])+2 [])
-> 3d6+2
<- 14 : ((3d6 [3 5 4])+2 [])
-> d6
<- 2 : (1d6 [2])
-> 3 * ( 1 d 100 / 2)
<- 99 : (3*((1d100 [66])/2 []) [])
-> 3b4d6
<- 13 : (3b(4d6 [5 2 6 1]) [2 5 6])
-> 2w4d6
<- 3 : (2w(4d6 [2 1 2 5]) [1 2])
-> d%
<- 7 : (1d% [7])
-> 4dF
<- 0 : (4dF [1 0 -1 0])
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

* for Savage Worlds we need a few more things to help support wild dice
* support generating a slice of ASTs via `,` delimiter for example:
```
-> d6,d8
<- 1 : (1d6 [1]) , 3 : (1d8 [3])
```

* it would be great if this could support best/worst...
```
-> 1b(d6,d8)
<- 3 : (1d8 [3]) (best of 1 : (1d6 [1]) , 3 : (1d8 [3]))
```
