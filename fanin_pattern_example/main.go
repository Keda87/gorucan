package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const csv1 = `
email,saving_amount,debt_amount
alice.johnson@example.com,7845.22,1200.00
bob.smith@example.com,2350.10,540.75
carol.lee@example.com,12350.00,9823.55
david.brown@example.com,475.65,0.00
emily.white@example.com,9400.33,2500.00
frank.green@example.com,658.00,340.10
grace.adams@example.com,18000.55,9500.12
henry.evans@example.com,2750.30,470.00
irene.wilson@example.com,3100.44,1200.88
jack.hall@example.com,5800.99,780.50
karen.young@example.com,9050.12,4200.00
leo.king@example.com,430.00,150.00
mia.scott@example.com,12200.77,6200.55
noah.adams@example.com,6700.00,3400.00
olivia.lee@example.com,300.45,80.00`

const csv2 = `
email,saving_amount,debt_amount
peter.turner@example.com,15800.22,1100.00
quinn.baker@example.com,7200.80,3500.25
rachel.morris@example.com,5200.90,220.10
sam.carter@example.com,630.40,100.00
tina.morgan@example.com,8700.00,5000.00
ursula.hughes@example.com,4050.30,2000.50
victor.kelly@example.com,950.55,400.20
wendy.reed@example.com,11200.99,7500.75
xavier.james@example.com,2500.00,800.00
yasmin.bennett@example.com,780.10,300.00
zack.perry@example.com,13500.88,9800.50
abby.stewart@example.com,3700.25,2100.90
brian.ellis@example.com,1600.00,500.00
claire.ross@example.com,200.00,50.00
dylan.bell@example.com,4200.80,1400.00`

const csv3 = `email,saving_amount,debt_amount
ella.henderson@example.com,13200.50,7200.00
freddie.morgan@example.com,450.75,100.50
george.clark@example.com,8900.00,3100.10
hannah.wright@example.com,2200.20,850.00
isaac.walker@example.com,17500.00,9200.90
julia.evans@example.com,780.00,300.00
kevin.mitchell@example.com,6500.60,4000.00
lily.phillips@example.com,14200.30,7800.45
matthew.cook@example.com,310.90,50.00
natalie.ward@example.com,9900.10,6300.00
owen.foster@example.com,5700.00,2900.50
paula.bailey@example.com,1200.25,450.10
quentin.hughes@example.com,10400.99,8700.30
rebecca.patterson@example.com,8400.00,3200.80
steven.bennett@example.com,3500.40,1400.50`

type Data struct {
	Email    string
	Saving   float64
	Debt     float64
	NetWorth float64
}

func readCSV(source string) <-chan Data {
	out := make(chan Data)

	go func(txt string) {
		defer close(out)

		lines := strings.Split(txt, "\n")
		for _, row := range lines {
			if row == "" {
				continue
			}

			columns := strings.Split(row, ",")
			if columns[0] == "email" {
				continue
			}

			saving, _ := strconv.ParseFloat(columns[1], 64)
			debt, _ := strconv.ParseFloat(columns[2], 64)
			out <- Data{
				Email:    columns[0],
				Saving:   saving,
				Debt:     debt,
				NetWorth: saving - debt,
			}

		}

	}(source)
	return out
}

func mergeFanIn(channels ...<-chan Data) <-chan Data {
	var (
		wg  sync.WaitGroup
		out = make(chan Data)
	)

	output := func(source <-chan Data) {
		defer wg.Done()

		for i := range source {
			out <- i
		}
	}

	for _, chans := range channels {
		wg.Add(1)
		go output(chans)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	start := time.Now()

	out1 := readCSV(csv1)
	out2 := readCSV(csv2)
	out3 := readCSV(csv3)

	// gather multiple channels into one processing function.
	merged := mergeFanIn(out1, out2, out3)

	for i := range merged {
		fmt.Printf("output: %s saving: %.2f but has networth %.2f\n", i.Email, i.Saving, i.NetWorth)
	}

	fmt.Println("execution: ", time.Since(start))
}
