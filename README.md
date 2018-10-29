A pet project, does not really work for now.

Made in Go, using [Evolutionary optimization library for Go](https://github.com/MaxHalford/eaopt)

# What is this

The goal of this code is to use [Genetic Algorithms](https://en.wikipedia.org/wiki/Genetic_algorithm) to create ASCII art images from a source `png` file. In GA terms, the _genome_ is a grid of characters mutated randomly, and the _fitness evaluation_ is done comparing the rendered text with the goal image.

Maybe it could work with a good combination of parameters for the genetic algorithm. For now it can only produce satisfactory results with the `--simple` flag and small images.

```
Usage:
  main [OPTIONS]

This command reads ./image.png and tries to create ASCII representation using Genetic Algorithms

Application Options:
  -f, --file=   Path to the input .png image (default: ./image.png)
  -s, --simple  Uses # and ' ' instead of any character

Help Options:
  -h, --help    Show this help message

```

```bash
$ go run main.go -s -f test-img/halfpounds-small.png
```
