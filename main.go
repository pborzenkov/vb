package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/rumyantseva/go-velobike/velobike"
)

func main() {
	id := os.Getenv("VELOBIKE_ID")
	if id == "" {
		log.Fatal("VELOBIKE_ID env var is not set")
	}
	pass := os.Getenv("VELOBIKE_PASS")
	if pass == "" {
		log.Fatal("VELOBIKE_PASS env var is not set")
	}

	tp := velobike.BasicAuthTransport{
		Username: id,
		Password: pass,
	}
	client := velobike.NewClient(tp.Client())
	auth, _, err := client.Authorization.Authorize()
	if err != nil {
		log.Fatalf("failed to authorize with velobike: %v", err)
	}
	client.SessionId = auth.SessionId

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
	fmt.Fprint(w, "ADDRESS\tFREE BIKES\tFREE PLACES\n")
	for _, f := range favs {
		fmt.Fprintf(w, "%s\t%d\t%d\n", *f.Address, *f.TotalPlaces-*f.FreePlaces, *f.FreePlaces)
	}
	w.Flush()
}
