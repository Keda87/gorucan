package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const legalAge = 25

const csvDummyText = `
name,born_year,wallet_amount,is_verified
Alice Johnson,1985,1543.27,True
Bob Smith,1992,328.90,False
Carol Lee,1978,8732.45,True
David Brown,2000,54.12,False
Emily White,1995,1200.00,True
Frank Green,1982,5023.88,False
Grace Hall,1988,764.00,True
Henry Young,1975,927.56,False
Ivy King,1999,405.75,True
Jack Turner,1983,7801.50,True
Karen Scott,1990,124.20,False
Leo Adams,2002,58.33,True
Mia Baker,1998,349.99,False
Noah Carter,1993,10500.00,True
Olivia Diaz,1981,679.45,False
Paul Evans,2021,2374.00,True
Quinn Foster,1996,95.50,False
Ruby Garcia,1989,4250.75,True
Samuel Harris,1984,188.80,False
Tina Johnson,2001,3023.15,True
`

type Data struct {
	Name         string
	BornYear     int
	WalletAmount float64
	IsVerified   bool
}

func httpGetter(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil { // check if http request error.
		if ctx.Err() != nil { // check if the error caused by context timeout.
			return "", ctx.Err()
		}
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func ParsingCSV(ctx context.Context, content string) <-chan Data {
	c := make(chan Data)
	go func(text string) {
		defer close(c)

		rows := strings.Split(text, "\n")
		for i := 0; i < len(rows); i++ {
			if i == 0 || rows[i] == "" {
				continue
			}

			cols := strings.Split(rows[i], ",")
			if cols[0] == "name" { // skip header.
				continue
			}

			yearInt, _ := strconv.ParseInt(cols[1], 10, 32)
			wallet, _ := strconv.ParseFloat(cols[2], 64)
			data := Data{Name: cols[0],
				BornYear:     int(yearInt),
				WalletAmount: wallet,
				IsVerified:   cols[3] == "True",
			}

			select {
			case <-ctx.Done():
				return
			case c <- data:
			}
		}
	}(content)
	return c
}

func FilterVerifiedAdult(ctx context.Context, row <-chan Data) <-chan Data {
	c := make(chan Data)
	go func() {
		defer close(c)

		for i := range row {
			age := time.Now().Year() - i.BornYear
			isLegalAgeAndVerified := age >= legalAge && i.IsVerified
			if isLegalAgeAndVerified {
				select {
				case <-ctx.Done():
					return
				case c <- i:
				}
			}
		}
	}()
	return c
}

func VerifyLegalAPI(ctx context.Context, row <-chan Data, workerCount int) <-chan Data {
	c := make(chan Data)

	var wg sync.WaitGroup
	wg.Add(workerCount)

	for w := 0; w < workerCount; w++ {
		go func(workedID int) {
			defer wg.Done()

			for i := range row {
				res, _ := httpGetter(ctx, "https://httpbin.org/anything/background-check")
				if res != "" {
					select {
					case <-ctx.Done():
						return
					case c <- i:
					}
				}
			}
		}(w)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	return c
}

func Display(ctx context.Context, row <-chan Data) {
	select {
	case <-ctx.Done():
		return
	default:
		for o := range row {
			fmt.Println(o.Name, ": successfully checked")
		}
	}
}

func main() {
	start := time.Now()
	ctx := context.Background()

	var (
		people   = ParsingCSV(ctx, csvDummyText)
		adults   = FilterVerifiedAdult(ctx, people)
		verified = VerifyLegalAPI(ctx, adults, 3)
	)

	Display(ctx, verified)
	fmt.Println("execution: ", time.Since(start))
}
