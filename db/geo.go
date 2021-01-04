package db

import (
	"app/types"
	"fmt"
	"strings"
)

//GetCityByName - get city model by name from geo table
func GetCityByName(name string) (types.CityModel, error) {
	row := database.QueryRow("select * from public.geo WHERE LOWER(city)=$1", strings.ToLower(name))
	city := types.CityModel{}

	err := row.Scan(&city.ID, &city.CountryEn, &city.RegionEn, &city.CityEn, &city.Country, &city.Region, &city.City, &city.Lat, &city.Lng, &city.Population)
	if err != nil {
		fmt.Println("city not found")
		return city, err
	}
	return city, nil
}

//GetRandomCityByLetter - get random city model from geo by first letter
func GetRandomCityByLetter(letter string) (types.CityModel, error) {
	query := fmt.Sprintf("select * from public.geo WHERE LOWER(city) LIKE '%s%s' ORDER BY random() LIMIT 1", strings.ToLower(letter), "%")
	fmt.Println(query)

	row := database.QueryRow(query)
	city := types.CityModel{}

	err := row.Scan(&city.ID, &city.CountryEn, &city.RegionEn, &city.CityEn, &city.Country, &city.Region, &city.City, &city.Lat, &city.Lng, &city.Population)
	if err != nil {
		fmt.Println("no random city")
		return city, err
	}
	return city, nil
}
