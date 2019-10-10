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
-> 3 * 1d100 /2
-> d6
<- ERROR failed to parse d6
-> 1d6
<- 1 : (1 d 6 -> 1 [1])
<- 115 : ((3 * (1 d 100 -> 77 [77]) -> 231) / 2 -> 115)
-> 3 * ( 1 d 100 / 2)
<- 120 : (3 * ((1 d 100 -> 80 [80]) / 2 -> 40) -> 120)
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
