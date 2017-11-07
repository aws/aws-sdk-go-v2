package modeledendpoints_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/modeledendpoints"
)

func ExampleResolver_Partitions() {
	partitions := modeledendpoints.NewDefaultResolver().Partitions()

	for _, p := range partitions {
		fmt.Println("Regions for", p.ID())
		for id := range p.Regions() {
			fmt.Println("*", id)
		}

		fmt.Println("Services for", p.ID())
		for id := range p.Services() {
			fmt.Println("*", id)
		}
	}
}
