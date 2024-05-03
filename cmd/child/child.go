package child

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nerijusro/scootinAboot/types"
	"github.com/nerijusro/scootinAboot/types/enums"
)

type MobileClientDummy struct {
	basePath string
}

func NewMobileClientDummy(basePath string) *MobileClientDummy {
	return &MobileClientDummy{basePath: basePath}
}

func (c *MobileClientDummy) Run(wg *sync.WaitGroup) {
	clientStaticApiKey, err := c.getStaticApiKey("client")
	if err != nil {
		log.Fatal(err)
	}

	adminStaticApiKey, err := c.getStaticApiKey("admin")
	if err != nil {
		log.Fatal(err)
	}

	err = c.createScooters(*adminStaticApiKey)
	if err != nil {
		log.Fatal(err)
	}

	clientOne := c.spawnUser(*clientStaticApiKey, "Ferrucio Lamborghini")
	clientTwo := c.spawnUser(*clientStaticApiKey, "Enzo Ferrari")
	clientThree := c.spawnUser(*clientStaticApiKey, "Carroll Shelby")

	wg.Add(1)
	go c.simulateClient(wg, *clientStaticApiKey, clientOne)

	wg.Add(1)
	go c.simulateClient(wg, *clientStaticApiKey, clientTwo)

	wg.Add(1)
	go c.simulateClient(wg, *clientStaticApiKey, clientThree)
}

func (c *MobileClientDummy) simulateClient(wg *sync.WaitGroup, staticApiKey string, client *types.MobileClient) error {
	defer wg.Done()
	for i := 0; i > -1; i++ {
		err := c.simulateTrip(staticApiKey, client.ID.String())
		if err != nil {
			break
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}

func (c *MobileClientDummy) simulateTrip(staticApiKey string, clientID string) error {
	scootersResponse, err := c.getScootersInArea(staticApiKey, 24.0, 26.0, 53.0, 55.0)
	if err != nil {
		log.Fatal(err)
		return err
	}

	scooter := scootersResponse.Scooters[rand.Intn(len(scootersResponse.Scooters))]
	startTripResponse, err := c.startTrip(staticApiKey, clientID, scooter.ID)
	if err != nil {
		return err
	}

	if startTripResponse == nil {
		return nil
	}

	var newLocation types.Location
	for i := 0; i < 4; i++ {
		newLocation = c.getNewRandomLocation()
		err = c.updateTrip(
			staticApiKey,
			clientID,
			startTripResponse.TripID,
			types.TripUpdateRequest{
				Location:    newLocation,
				CreatedAt:   time.Now(),
				IsFinishing: false,
				Sequence:    i + 2,
			})

		if err != nil {
			log.Fatal(err)
			return err
		}
		time.Sleep(3 * time.Second)
	}

	newLocation = c.getNewRandomLocation()
	err = c.updateTrip(
		staticApiKey,
		clientID,
		startTripResponse.TripID,
		types.TripUpdateRequest{
			Location:    newLocation,
			CreatedAt:   time.Now(),
			IsFinishing: true,
			Sequence:    6,
		})

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (c *MobileClientDummy) updateTrip(staticApiKey string, clientID string, tripID uuid.UUID, updatedTrip types.TripUpdateRequest) error {
	requestBody := updatedTrip

	marshalledRequestBody, _ := json.Marshal(requestBody)
	request, err := http.NewRequest(http.MethodPut, c.basePath+"/client/trips/"+tripID.String(), bytes.NewBuffer(marshalledRequestBody))
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("x-api-key", staticApiKey)
	request.Header.Set("client-id", clientID)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println("Error updating trip", err)
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("Expected status code 200 but got", resp.StatusCode)
		return nil
	}

	return nil
}

func (c *MobileClientDummy) startTrip(staticApiKey string, clientID string, scooterID uuid.UUID) (*types.TripEvent, error) {
	requestBody := types.StartTripRequest{
		ScooterID: scooterID,
		CreatedAt: time.Now(),
	}

	marshalledRequestBody, _ := json.Marshal(requestBody)
	request, err := http.NewRequest(http.MethodPost, c.basePath+"/client/trips", bytes.NewBuffer(marshalledRequestBody))
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("x-api-key", staticApiKey)
	request.Header.Set("client-id", clientID)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println("Error creating trip", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		log.Println("Expected status code 201 but got", resp.StatusCode)

		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			log.Fatal(err)
		}

		log.Print("Message", response)
		return nil, nil
	}

	var tripEvent types.TripEvent
	if err := json.NewDecoder(resp.Body).Decode(&tripEvent); err != nil {
		log.Fatal(err)
	}

	return &tripEvent, nil
}

func (c *MobileClientDummy) getScootersInArea(staticApiKey string, x1, x2, y1, y2 float64) (*types.GetScootersResponse, error) {
	requestBody := types.GetScootersQueryParameters{
		Availability: string(enums.Available),
		X1:           x1,
		X2:           x2,
		Y1:           y1,
		Y2:           y2,
	}

	request, err := http.NewRequest(http.MethodGet, getScootersByAreaUrl(c.basePath, requestBody), nil)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("x-api-key", staticApiKey)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println("Error getting scooters", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("Expected status code 200 but got", resp.StatusCode)
		return nil, err
	}

	var scootersResponse types.GetScootersResponse
	if err := json.NewDecoder(resp.Body).Decode(&scootersResponse); err != nil {
		log.Fatal(err)
	}

	return &scootersResponse, nil
}

func (c *MobileClientDummy) createScooters(staticApiKey string) error {
	requestBody := types.CreateScooterRequest{
		Location: types.Location{
			Longitude: 25.279651,
			Latitude:  54.687157,
		},
		IsAvailable: true,
	}

	marshalledRequestBody, _ := json.Marshal(requestBody)
	for i := 0; i < 3; i++ {
		request, err := http.NewRequest(http.MethodPost, c.basePath+"/admin/scooters", bytes.NewBuffer(marshalledRequestBody))
		if err != nil {
			log.Fatal(err)
		}

		request.Header.Set("x-api-key", staticApiKey)
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Println("Error creating scooter", err)
			return nil
		}

		if resp.StatusCode != http.StatusCreated {
			log.Println("Expected status code 201 but got", resp.StatusCode)
			return nil
		}
	}

	return nil
}

func (c *MobileClientDummy) spawnUser(staticApiKey string, fullName string) *types.MobileClient {
	requestBody := types.CreateUserRequest{
		FullName: fullName,
	}

	marshalledRequestBody, _ := json.Marshal(requestBody)
	request, err := http.NewRequest(http.MethodPost, c.basePath+"/client/users", bytes.NewBuffer(marshalledRequestBody))
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("x-api-key", staticApiKey)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println("Error creating user", err)
		return nil
	}

	if resp.StatusCode != http.StatusCreated {
		log.Println("Expected status code 201 but got", resp.StatusCode)
		return nil
	}

	var userResponse types.MobileClient
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		log.Fatal(err)
	}

	return &userResponse
}

func (c *MobileClientDummy) getStaticApiKey(role string) (*string, error) {
	var url strings.Builder
	url.WriteString(c.basePath)
	url.WriteString("/")
	url.WriteString(role)
	url.WriteString("/auth")

	resp, err := http.Get(url.String())
	if err != nil {
		log.Println("Error getting static api key", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Expected status code 200 but got", resp.StatusCode)
		return nil, err
	}

	var authResponse types.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		log.Fatal(err)
	}

	return &authResponse.StaticApiKey, nil
}

func (c *MobileClientDummy) getNewRandomLocation() types.Location {
	return types.Location{
		Latitude:  53 + rand.Float64()*(55-53),
		Longitude: 24 + rand.Float64()*(26-24),
	}
}
