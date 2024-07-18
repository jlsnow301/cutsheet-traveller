package travel

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"googlemaps.github.io/maps"
)

// getDirections fetches directions using Google Maps API.
func getDirections(origin, destination string, event *time.Time) *maps.Route {
	client, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_MAPS_API_KEY")))
	if err != nil {
		color.Red(fmt.Sprintf("Error creating Google Maps client: %v", err))
		return nil
	}

	// If the time is in the past, just say "now"
	var eventTime string
	if event.Before(time.Now()) {
		eventTime = "now"
	} else {
		eventTime = strconv.FormatInt(event.Unix(), 10)
	}

	request := &maps.DirectionsRequest{
		Origin:        origin,
		Destination:   destination,
		DepartureTime: eventTime,
	}

	routes, _, err := client.Directions(context.Background(), request)
	if err != nil {
		color.Red(fmt.Sprintf("Error fetching directions: %v", err))
		return nil
	}

	if len(routes) > 0 {
		return &routes[0]
	}
	return nil
}

func getDistanceText(directionsResult *maps.Route) string {
	if directionsResult == nil {
		color.Red("No directions found.")
		return ""
	}

	leg := directionsResult.Legs[0]
	distanceText := leg.Distance.HumanReadable

	return distanceText
}

func getDoubleDistance(distanceText string) (float64, error) {
	// Split the miles off the end of the string
	splitText := strings.SplitN(distanceText, " ", 2)
	if len(splitText) < 2 {
		return 0, errors.New("no distance found")
	}

	distance, err := strconv.ParseFloat(splitText[0], 64)
	if err != nil {
		return 0, err
	}

	double := 2.0
	// Multiply by 2 for round trips
	distance = distance * double

	return distance, nil
}

// GetBaseTravelTime gets the base travel time based on the origin and destination.
func GetBaseTravelDistance(origin, destination string, event *time.Time) (float64, error) {
	directionsResult := getDirections(origin, destination, event)
	distanceText := getDistanceText(directionsResult)
	if distanceText == "" {
		return 0, errors.New("no directions found")
	}

	roundTripMiles, err := getDoubleDistance(distanceText)
	if err != nil {
		return 0, err
	}

	return roundTripMiles, nil
}
