package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/entitymanager"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type FullTable struct {
	PK          string    `dynamodbav:"PK,partition"`
	SK          string    `dynamodbav:"SK,sort" dynamodbindex:"GSI,global,sort"`
	Result      any       `dynamodbav:"result,coverter|json"`
	LastUpdated time.Time `dynamodbav:"LastUpdated"`
	GSI         int64     `dynamodbav:"GSI" dynamodbindex:"GSI,global,partition;GSI_PK_Index,local,sort"`
	TS          int64     `dynamodbav:"ts"`
}

type Race struct {
	RaceID  string `dynamodbav:"PK,partition"` // race-%d
	ClassID string `dynamodbav:"SK,sort"`      // class-%d
}

type RaceResult struct {
	RaceID  string `dynamodbav:"PK,partition"` // race-%d
	RacerID string `dynamodbav:"SK,sort"`      // racer-%d
	TimeMs  int64  `dynamodbav:"time_ms"`      // race time in milliseconds
}

type Racer struct {
	RaceID string `dynamodbav:"PK,partition"` //. race-%d
	Name   string `dynamodbav:"SK,sort"`      // racer-%d
}

func getTable[T any](client entitymanager.Client, tableName string) *entitymanager.Table[T] {
	sch, err := entitymanager.NewSchema[T]()
	if err != nil {
		panic(err)
	}

	sch = sch.WithTableName(aws.String(tableName))

	tbl, err := entitymanager.NewTable(client, func(options *entitymanager.TableOptions[T]) {
		options.Schema = sch
	})
	if err != nil {
		panic(err)
	}

	return tbl
}

func createTable(client entitymanager.Client, tableName string) context.CancelFunc {
	// create the full table with gsi and lsi and all
	log.Println("Creating full table")
	tbl := getTable[FullTable](client, tableName)

	if exists, err := tbl.Exists(context.Background()); !exists || err != nil {
		if err != nil {
			panic(err)
		}

		if err := tbl.CreateWithWait(context.Background(), time.Minute*2); err != nil {
			panic(err)
		}
		log.Println("Created full table")

		return func() {
			log.Println("Deleting full table")
			if err := tbl.DeleteWithWait(context.Background(), time.Minute*2); err != nil {
				panic(err)
			}
			log.Println("Deleted full table")
		}
	}

	return func() {}
}

func generateData(client entitymanager.Client, tableName string, count int) {
	generateRaces(client, tableName, count)
	generateRacers(client, tableName, count)
	generateRaceResults(client, tableName, count)
}

func generateRaces(client entitymanager.Client, tableName string, count int) {
	racesTbl := getTable[Race](client, tableName)

	for c := range count {
		_, err := racesTbl.PutItem(context.Background(), &Race{
			RaceID:  fmt.Sprintf("race-%d", c),
			ClassID: fmt.Sprintf("class-%d", c),
		})

		if err != nil {
			log.Printf("Error writing race")
		}
	}
}

func generateRacers(client entitymanager.Client, tableName string, count int) {
	racersTbl := getTable[Racer](client, tableName)

	for c := range count {
		_, err := racersTbl.PutItem(context.Background(), &Racer{
			RaceID: fmt.Sprintf("race-%d", c),
			Name:   fmt.Sprintf("name-%d", c),
		})

		if err != nil {
			log.Printf("Error writing racers")
		}
	}
}
func generateRaceResults(client entitymanager.Client, tableName string, count int) {
	raceResultsTbl := getTable[RaceResult](client, tableName)

	rand.Seed(time.Now().UnixNano())

	for c := range count {
		for d := range count {
			// Simulate a race time in milliseconds (e.g. 60s–80s)
			timeMs := int64(60000 + rand.Intn(20000))

			_, err := raceResultsTbl.PutItem(context.Background(), &RaceResult{
				RaceID:  fmt.Sprintf("race-%d", c),
				RacerID: fmt.Sprintf("racer-%d", d),
				TimeMs:  timeMs,
			})

			if err != nil {
				log.Printf("Error writing race result")
			}
		}
	}
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}

	tableName := fmt.Sprintf("table_%s", time.Now().Format("2006_01_02_15_04_05"))
	println(tableName)

	ddb := dynamodb.NewFromConfig(cfg)
	count := 10

	cancel := createTable(ddb, tableName)
	defer cancel()

	generateData(ddb, tableName, count)

	racerTbl := getTable[Racer](ddb, tableName)
	_ = racerTbl
	raceTbl := getTable[Race](ddb, tableName)
	_ = raceTbl
	raceResultTbl := getTable[RaceResult](ddb, tableName)

	ctx := context.Background()

	// show leaderboard table - max 10 rows
	showLeaderboard(ctx, raceResultTbl, 10)

	// show race results from random race
	showRandomRaceResults(ctx, raceResultTbl, count)

	// show rankings
	showRaceRankings(ctx, raceResultTbl, 10)
}

func showLeaderboard(ctx context.Context, raceResultTbl *entitymanager.Table[RaceResult], maxRows int) {
	// Only consider RaceResult items (SK starts with "racer-")
	f := expression.Name("SK").BeginsWith("racer-")
	expr, err := expression.NewBuilder().WithFilter(f).Build()
	if err != nil {
		log.Printf("error building expression for leaderboard scan: %v", err)
		return
	}

	// Track best (smallest) time per racer across all races
	type racerBest struct {
		RacerID string
		TimeMs  int64
	}

	bestByRacer := map[string]int64{}
	for res := range raceResultTbl.Scan(ctx, expr) {
		if res.Error() != nil {
			log.Printf("Scan() error while building leaderboard: %v", res.Error())
			continue
		}

		item := res.Item()
		if item == nil {
			continue
		}

		t := item.TimeMs
		current, ok := bestByRacer[item.RacerID]
		if !ok || t < current {
			bestByRacer[item.RacerID] = t
		}
	}

	entries := make([]racerBest, 0, len(bestByRacer))
	for id, t := range bestByRacer {
		entries = append(entries, racerBest{RacerID: id, TimeMs: t})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].TimeMs < entries[j].TimeMs
	})

	fmt.Println("== Leaderboard (top racers by best time) ==")
	for i, e := range entries {
		if i >= maxRows {
			break
		}
		fmt.Printf("%2d. %-10s %d ms\n", i+1, e.RacerID, e.TimeMs)
	}
}

func showRandomRaceResults(ctx context.Context, raceResultTbl *entitymanager.Table[RaceResult], races int) {
	if races <= 0 {
		return
	}

	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(races)
	raceID := fmt.Sprintf("race-%d", idx)

	// Query a single race, but only RaceResult items (SK starts with "racer-")
	keyCond := expression.Key("PK").Equal(expression.Value(raceID)).And(
		expression.Key("SK").BeginsWith("racer-"),
	)
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		log.Printf("error building expression for random race query: %v", err)
		return
	}

	// Collect all results so we can sort by time
	results := make([]*RaceResult, 0)
	for res := range raceResultTbl.Query(ctx, expr) {
		if res.Error() != nil {
			log.Printf("Query() error for %s: %v", raceID, res.Error())
			continue
		}

		item := res.Item()
		if item == nil {
			continue
		}

		results = append(results, item)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].TimeMs < results[j].TimeMs
	})

	fmt.Printf("== Results for %s ==\n", raceID)
	for i, r := range results {
		fmt.Printf("%2d. %-10s %d ms\n", i+1, r.RacerID, r.TimeMs)
	}
}

func showRaceRankings(ctx context.Context, raceResultTbl *entitymanager.Table[RaceResult], maxRows int) {
	// Only consider RaceResult items (SK starts with "racer-")
	f := expression.Name("SK").BeginsWith("racer-")
	expr, err := expression.NewBuilder().WithFilter(f).Build()
	if err != nil {
		log.Printf("error building expression for rankings scan: %v", err)
		return
	}

	// For each race, track the best (smallest) time and the pilot
	type raceBest struct {
		RaceID    string
		BestTime  int64
		BestRacer string
	}

	bestByRace := map[string]raceBest{}
	for res := range raceResultTbl.Scan(ctx, expr) {
		if res.Error() != nil {
			log.Printf("Scan() error while building rankings: %v", res.Error())
			continue
		}

		item := res.Item()
		if item == nil {
			continue
		}

		t := item.TimeMs
		current, ok := bestByRace[item.RaceID]
		if !ok || t < current.BestTime {
			bestByRace[item.RaceID] = raceBest{
				RaceID:    item.RaceID,
				BestTime:  t,
				BestRacer: item.RacerID,
			}
		}
	}

	races := make([]raceBest, 0, len(bestByRace))
	for _, v := range bestByRace {
		races = append(races, v)
	}

	sort.Slice(races, func(i, j int) bool {
		return races[i].BestTime < races[j].BestTime
	})

	fmt.Println("== Race rankings (by best pilot time) ==")
	for i, e := range races {
		if i >= maxRows {
			break
		}
		fmt.Printf("%2d. %-10s %d ms (best: %s)\n", i+1, e.RaceID, e.BestTime, e.BestRacer)
	}
}
