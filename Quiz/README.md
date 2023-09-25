# GO Quiz

A console program that does quiz game from the selected .csv file

## Prerequisites

Go 1.20+

## Usage

Download the project and run it:

```
go run main.go
```

This command can take up to three arguments:

- `--csv` to specify the path to the .csv file to be the source of the quiz (default is `./problems.csv`).
- `--shuffle` to specify should the problems be shuffled or not (`true` or `false`, default is `false`).
- `--time` to set the time limit for the whole quiz (`int`, seconds, default is `30`).

*Example to create a quiz with a different .csv file with shuffled problems and a time limit of a minute:*

```
go run main.go --csv /home/location/quiz.csv --shuffle true --time 60
```

