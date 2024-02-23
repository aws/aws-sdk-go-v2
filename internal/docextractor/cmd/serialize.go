package main

import (
	"fmt"
	"encoding/json"
	"os"
)



func serialize(packageName string, items map[string]JewelryItem) {
	// there are single files getting through here
	clientData := map[string]JewelryItem{}
	for k, v := range items {
		if k == "packageDocumentation" {
			clientData[k] = v
			continue
		}
		clientData[v.Name] = v
		if v.Name == "Client" {
			for _, m := range v.Members {
				if m.Name == "Options" {
					continue
				}
				// m.Input = fmt.Sprintf("%vInput", m.Name)
				// m.Output = fmt.Sprintf("%vOutput", m.Name)
				clientData[m.Name] = m
			}
		}
	}
	content, err := json.Marshal(clientData)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile(fmt.Sprintf("../clients/%v.json", packageName), content, 0644)
	if err != nil {
		panic(err)
	}

	for _, item := range clientData {
		if item.Name == "" || item.Name == "packageDocumentation" {
			continue
		}
		content, err := json.Marshal(item)
		if err != nil {
			fmt.Println(err)
		}
		err = os.WriteFile(fmt.Sprintf("../public/members/-aws-sdk-client-%v.%v.%v.json", packageName, item.Name, string(item.Type)), content, 0644)
		if err != nil {
			panic(err)
		}
	}


	typeData := map[string][]string{}
	for _, item := range clientData {
		if item.Name == "" || item.Name == "packageDocumentation" {
			continue
		}
		val, ok := typeData[string(item.Type)]
		if !ok {
			val = []string{}
		}
		val = append(val, item.Name)
		typeData[string(item.Type)] = val
	}
	content, err = json.Marshal(typeData)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile(fmt.Sprintf("../public/members/-aws-sdk-client-%v.json", packageName), content, 0644)
	if err != nil {
		panic(err)
	}
}