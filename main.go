package main

import (
	"image"
	"os"

	"github.com/MaxHalford/eaopt"
	tm "github.com/buger/goterm"
	"github.com/carlosms/ga-ascii-art/genome"
	"github.com/carlosms/ga-ascii-art/img"
	flags "github.com/jessevdk/go-flags"
)

type options struct {
	File      string `short:"f" long:"file" description:"Path to the input .png image" default:"./image.png"`
	PoundOnly bool   `short:"s" long:"simple" description:"Uses # and ' ' instead of any character"`
}

var opt options

var parser = flags.NewParser(&opt, flags.Default)

func main() {
	parser.LongDescription =
		`This command reads ./image.png and tries to create ASCII representation using Genetic Algorithms`
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	goal := img.ToGray(img.ReadPNG(opt.File))

	img.SavePNG(goal, "gray.png")

	ga(goal)
}

func ga(goal *image.Gray) error {
	var conf = eaopt.NewDefaultGAConfig()

	conf.NPops = 1
	conf.PopSize = 150
	conf.NGenerations = 250
	/*
		conf.Model = eaopt.ModDownToSize{
			NOffsprings: conf.PopSize * 2,
			SelectorA: eaopt.SelTournament{
				NContestants: 3,
			},
			SelectorB: eaopt.SelTournament{
				NContestants: 3,
			},
			MutRate:   0.5,
			CrossRate: 0.7,
		}
	*/
	/*
		conf.Model = eaopt.ModSteadyState{
			Selector: eaopt.SelTournament{
				NContestants: 3,
			},
			KeepBest:  true,
			MutRate:   0.5,
			CrossRate: 0.7,
		}
	*/

	//conf.RNG = rand.New(rand.NewSource(42))

	var ga, err = conf.NewGA()
	if err != nil {
		return err
	}

	// Add a custom print function to track progress
	ga.Callback = func(ga *eaopt.GA) {
		tm.Clear()
		tm.MoveCursor(1, 1)
		best := ga.HallOfFame[0]
		tm.Println(best.Genome.(*genome.ASCIIGenome).String())
		tm.Printf("Generation   %d\nBest fitness %f\n", ga.Generations, best.Fitness)
		tm.Flush()
	}

	// Run the GA
	err = ga.Minimize(genome.NewASCIIGenome(goal, opt.PoundOnly))
	if err != nil {
		return err
	}

	// Best genome
	asciiGenome := ga.HallOfFame[0].Genome.(*genome.ASCIIGenome)
	img.SavePNG(asciiGenome.Image(), "best.png")

	return nil
}
