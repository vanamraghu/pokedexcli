package pokeapis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokecache"
	"sync"
)

const LOCATION = "https://pokeapi.co/api/v2/location-area"

type LocationStructure struct {
	LocationCount int               `json:"count"`
	NextUrl       *string           `json:"next"`
	PreviousUrl   *string           `json:"previous"`
	Results       []locationDetails `json:"results"`
}

type locationDetails struct {
	LocationName string `json:"name"`
	Url          string `json:"url"`
}

var responseData *LocationStructure

//var c = pokecache.NewCache(5 * time.Second)

var c = pokecache.Cache{
	Mux:       &sync.Mutex{},
	CacheData: make(map[string]pokecache.CacheEntry),
}

func getLocationData(url string) (LocationStructure, error) {
	res, err := http.Get(url)
	if err != nil {
		return LocationStructure{}, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return LocationStructure{}, err
	}
	err = res.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return LocationStructure{}, err
	}
	return *responseData, nil
}

func DisplayBackwardLocations() error {
	var cachedData []byte
	if responseData != nil {
		if responseData.PreviousUrl == nil {
			fmt.Println("Previous url doesn't exist")
		} else {
			url := *responseData.PreviousUrl
			fmt.Printf("Previous Url is %s\n", url)
			cd, ok := c.Get(url)
			if ok {
				fmt.Printf("Reading from cache data")
				fmt.Printf("%s\n", cd)
			} else {
				locationData, err := getLocationData(url)
				if err != nil {
					return err
				}
				data := addCachedData(url, locationData, cachedData)
				displayCachedData(data)
				return nil
			}
		}
	} else {
		err := fmt.Errorf("no location data, pls run map command to get locations")
		fmt.Println(err)
	}
	return nil
}

func DisplayLocations() error {
	// if cache has already data, check for the url and display from cache
	var cachedData []byte
	if responseData != nil {
		nextUrl := *responseData.NextUrl
		cd, ok := c.Get(nextUrl)
		if ok {
			fmt.Printf("Reading from cache data")
			fmt.Printf("%s\n", cd)
		} else {
			locationData, err := getLocationData(nextUrl)
			if err != nil {
				return err
			}
			data := addCachedData(nextUrl, locationData, cachedData)
			displayCachedData(data)
			return nil
		}
	} else {
		locationData, err := getLocationData(LOCATION)
		if err != nil {
			return err
		}

		data := addCachedData(LOCATION, locationData, cachedData)
		displayCachedData(data)
		return nil
	}
	return nil
}

func addCachedData(url string, locationData LocationStructure, data []byte) []byte {
	for _, val := range locationData.Results {
		data = append(data, []byte(val.LocationName+"\n")...)
	}
	c.Add(url, data)
	return data
}

func displayCachedData(data []byte) {
	fmt.Printf("%s\n", data)
}
