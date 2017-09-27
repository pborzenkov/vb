package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:     "stats",
		Aliases: []string{"s"},
		Short:   "Display velobike.ru usage statistics",
		Run: func(cmd *cobra.Command, args []string) {
			stats(cmd, args)
		},
	})
}

func parseDuration(ds string) (time.Duration, error) {
	var d, h, m, s int
	var dur time.Duration

	if _, err := fmt.Sscanf(ds, "%02d.%02d:%02d:%02d", &d, &h, &m, &s); err != nil {
		return 0, err
	}

	dur += time.Duration(d) * 24 * time.Hour
	dur += time.Duration(h) * time.Hour
	dur += time.Duration(m) * time.Minute
	dur += time.Duration(s) * time.Second

	return dur, nil
}

func stats(cmd *cobra.Command, args []string) {
	client := mustClientFromEnv()

	history, _, err := client.History.Get()
	if err != nil {
		log.Fatalf("failed to get the history: %v", err)
	}

	var payment, distance float64
	var rides int
	var duration time.Duration
	for _, h := range history.Items {
		switch *h.Type {
		case "Pay":
			payment += *h.Price
		case "Ride":
			if h.Rejected != nil && *h.Rejected {
				continue
			}

			rides++
			payment += *h.Price
			distance += float64(*h.Distance) / 1000
			d, err := parseDuration(*h.Duration)
			if err != nil {
				log.Fatalf("failed to parse duration: %v", err)
			}
			duration += d
		}
	}
	fmt.Printf("Total spendings: %.2f RUB\n", payment)
	fmt.Printf("Total rides:     %d\n", rides)
	fmt.Printf("Total distance:  %.1f km\n", distance)
	fmt.Printf("Total ride time: %v\n\n", duration)

	fmt.Printf("Average cost per ride: %.2f RUB/ride\n", payment/float64(rides))
	fmt.Printf("Cost per km:           %.2f RUB/km\n", payment/distance)
	fmt.Printf("Cost per minute:       %.2f RUB/min\n", payment/duration.Minutes())
}
