package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/rumyantseva/go-velobike/velobike"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:     "favourite",
		Aliases: []string{"f"},
		Short:   "List favourite parking spots",
		Run: func(cmd *cobra.Command, args []string) {
			favourite(cmd, args)
		},
	})
}

func favourite(cmd *cobra.Command, args []string) {
	client := mustClientFromEnv()

	park, _, err := client.Parkings.List()
	if err != nil {
		log.Fatalf("failed to get the list of parking spots: %v", err)
	}

	var favs []velobike.Parking
	for _, p := range park.Items {
		if p.IsFavourite != nil && *p.IsFavourite {
			favs = append(favs, p)
		}
	}
	if len(favs) == 0 {
		fmt.Printf("You don't have any favourite parking spots!\n")
		os.Exit(0)
	}

	sort.Slice(favs, func(i, j int) bool {
		return *favs[i].Id < *favs[j].Id
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprint(w, "ADDRESS\tFREE BIKES\tFREE PLACES\tIN SERVICE\n")
	for _, f := range favs {
		fmt.Fprintf(w, "%s\t%d\t%d\t", *f.Address, *f.TotalPlaces-*f.FreePlaces, *f.FreePlaces)
		if f.IsLocked == nil {
			fmt.Fprintf(w, "?\n")
		} else if *f.IsLocked == true {
			fmt.Fprintf(w, "✘\n")
		} else {
			fmt.Fprintf(w, "✔\n")
		}
	}
	w.Flush()
}
