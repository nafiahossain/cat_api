package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/server/web"
	"io/ioutil"
	"net/http"
	"log"
	"bytes"
	"fmt"
	"strings"
)

type CatController struct {
	web.Controller
}

func (c *CatController) Get() {
	c.TplName = "index.tpl"
}

func (c *CatController) GetCatData() {
	apiKey, err := web.AppConfig.String("cat_api_key")
	if err != nil {
		log.Println("Error fetching API key from configuration:", err)
		c.Data["json"] = map[string]interface{}{"error": "Internal server error"}
		c.ServeJSON()
		return
	}

	// Use a channel for API calls
	dataChan := make(chan map[string]interface{})

	go func() {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "https://api.thecatapi.com/v1/images/search?size=med&mime_types=jpg&format=json&has_breeds=true&order=RANDOM&page=0&limit=1", nil)
		req.Header.Add("x-api-key", apiKey)
		
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error making request:", err)
			dataChan <- map[string]interface{}{"error": "Failed to fetch cat data"}
			return
		}
		defer resp.Body.Close()
		
		body, _ := ioutil.ReadAll(resp.Body)
		var data []map[string]interface{}
		json.Unmarshal(body, &data)
		
		if len(data) > 0 {
			dataChan <- data[0]
		} else {
			dataChan <- map[string]interface{}{"error": "No cat data found"}
		}
	}()

	result := <-dataChan
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *CatController) GetBreeds() {
	// Retrieve the API key from the configuration
	apiKey, err := web.AppConfig.String("cat_api_key")
	if err != nil {
		log.Println("Error fetching API key from configuration:", err)
		c.Data["json"] = map[string]interface{}{"error": "Internal server error"}
		c.ServeJSON()
		return
	}
	
	dataChan := make(chan []map[string]interface{})
	
	go func() {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "https://api.thecatapi.com/v1/breeds", nil)
		req.Header.Add("x-api-key", apiKey)
		
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error fetching breeds:", err)
			dataChan <- []map[string]interface{}{{"error": "Failed to fetch breeds"}}
			return
		}
		defer resp.Body.Close()
		
		body, _ := ioutil.ReadAll(resp.Body)
		var data []map[string]interface{}
		json.Unmarshal(body, &data)
		
		dataChan <- data
	}()

	result := <-dataChan
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *CatController) GetBreedInfo() {
	// Get the breed ID from the URL parameters
	breedId := c.Ctx.Input.Param(":id")

	// Retrieve the API key from the configuration
	apiKey, err := web.AppConfig.String("cat_api_key")
	if err != nil {
		log.Println("Error fetching API key from configuration:", err)
		c.Data["json"] = map[string]interface{}{"error": "Internal server error"}
		c.ServeJSON()
		return
	}

	// Create a channel to handle the API response data
	dataChan := make(chan map[string]interface{})

	// Fetch data using a goroutine
	go func() {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "https://api.thecatapi.com/v1/images/search?breed_ids="+breedId+"&limit=10", nil)
		req.Header.Add("x-api-key", apiKey)

		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error fetching breed info:", err)
			dataChan <- map[string]interface{}{"error": "Failed to fetch breed info"}
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var data []map[string]interface{}
		json.Unmarshal(body, &data)

		if len(data) > 0 {
			result := map[string]interface{}{
				"images":     data,
				"breed_info": data[0]["breeds"].([]interface{})[0],
			}
			dataChan <- result
		} else {
			dataChan <- map[string]interface{}{"error": "No data found for this breed"}
		}
	}()

	// Send the result as JSON response
	result := <-dataChan
	c.Data["json"] = result
	c.ServeJSON()
}

type FavoriteRequest struct {
    ImageID string `json:"image_id"`
}

func (c *CatController) AddFavorite() {
    apiKey, err := web.AppConfig.String("cat_api_key")
    if err != nil {
        c.Ctx.Output.SetStatus(500)
        c.Data["json"] = map[string]string{"error": "Failed to get API key"}
        c.ServeJSON()
        return
    }

    rawBody, err := ioutil.ReadAll(c.Ctx.Request.Body)
    if err != nil {
        c.Ctx.Output.SetStatus(500)
        c.Data["json"] = map[string]string{"error": "Failed to read request body"}
        c.ServeJSON()
        return
    }
    c.Ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))

    var requestBody FavoriteRequest
    if err := json.Unmarshal(rawBody, &requestBody); err != nil {
        c.Ctx.Output.SetStatus(400)
        c.Data["json"] = map[string]string{"error": "Invalid request body: " + err.Error()}
        c.ServeJSON()
        return
    }

    if requestBody.ImageID == "" {
        c.Ctx.Output.SetStatus(400)
        c.Data["json"] = map[string]string{"error": "Missing or empty image_id in request body"}
        c.ServeJSON()
        return
    }

    catAPIRequestBody := map[string]string{
        "image_id": requestBody.ImageID,
        "sub_id":   "user43",
    }

    jsonData, err := json.Marshal(catAPIRequestBody)
    if err != nil {
        c.Ctx.Output.SetStatus(500)
        c.Data["json"] = map[string]string{"error": "Failed to marshal JSON data"}
        c.ServeJSON()
        return
    }

    errChan := make(chan error, 1)
    resultChan := make(chan map[string]interface{}, 1)

    go func() {
        req, err := http.NewRequest("POST", "https://api.thecatapi.com/v1/favourites", bytes.NewBuffer(jsonData))
        if err != nil {
            errChan <- err
            return
        }
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("x-api-key", apiKey)

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            errChan <- err
            return
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            errChan <- err
            return
        }

        var result map[string]interface{}
        if err := json.Unmarshal(body, &result); err != nil {
            errChan <- err
            return
        }

        if resp.StatusCode != http.StatusOK {
            errChan <- fmt.Errorf("API Error: %s", string(body))
            return
        }

        resultChan <- result
    }()

    select {
    case err := <-errChan:
        c.Ctx.Output.SetStatus(500)
        c.Data["json"] = map[string]string{"error": "Failed to add favorite: " + err.Error()}
    case result := <-resultChan:
        c.Data["json"] = result
    }
    c.ServeJSON()
}


func (c *CatController) GetFavorites() {
    apiKey, err := web.AppConfig.String("cat_api_key")
    if err != nil {
        c.Ctx.Output.SetStatus(500)
        c.Data["json"] = map[string]string{"error": "Failed to get API key"}
        c.ServeJSON()
        return
    }

    // Create channels for concurrent processing
    errChan := make(chan error, 1)
    resultChan := make(chan []map[string]interface{}, 1)

    go func() {
        req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/favourites", nil)
        if err != nil {
            errChan <- err
            return
        }
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("x-api-key", apiKey)

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            errChan <- err
            return
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            errChan <- err
            return
        }

        var result []map[string]interface{}
        if err := json.Unmarshal(body, &result); err != nil {
            errChan <- err
            return
        }

        resultChan <- result
    }()

    // Handle the result or error from the channels
    select {
    case err := <-errChan:
        c.Ctx.Output.SetStatus(500)
        c.Data["json"] = map[string]string{"error": "Failed to get favorites: " + err.Error()}
    case result := <-resultChan:
        c.Data["json"] = result
    }
    c.ServeJSON()
}

type VoteRequest struct {
    ImageID string `json:"image_id"`
    SubID   string `json:"sub_id"`
    Value   int    `json:"value"`
}

func (c *CatController) SubmitVote() {
    var voteReq VoteRequest

    // Helper function to send JSON response with status code
    sendJSONResponse := func(statusCode int, data interface{}) {
        c.Ctx.ResponseWriter.WriteHeader(statusCode)
        c.Data["json"] = data
        c.ServeJSON()
    }

    // Read and log the raw request body
    rawBody, err := ioutil.ReadAll(c.Ctx.Request.Body)
    if err != nil {
        sendJSONResponse(http.StatusInternalServerError, map[string]string{"error": "Failed to read request body: " + err.Error()})
        return
    }

    // Check if the body is empty
    if len(rawBody) == 0 {
        sendJSONResponse(http.StatusBadRequest, map[string]string{"error": "Received an empty request body"})
        return
    }

    // Reset the request body for further reading
    c.Ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))

    // Unmarshal the request body into the VoteRequest struct
    if err := json.Unmarshal(rawBody, &voteReq); err != nil {
        sendJSONResponse(http.StatusBadRequest, map[string]string{"error": "Invalid request body: " + err.Error()})
        return
    }

    // Validate the request data
    if voteReq.ImageID == "" || voteReq.SubID == "" {
        sendJSONResponse(http.StatusBadRequest, map[string]string{"error": "Missing required fields in request body"})
        return
    }

    apiKey, err := web.AppConfig.String("cat_api_key")
    if err != nil {
        sendJSONResponse(http.StatusInternalServerError, map[string]string{"error": "Failed to get API key"})
        return
    }

    // Prepare the request body for The Cat API
    catAPIRequestBody := map[string]interface{}{
        "image_id": voteReq.ImageID,
        "sub_id":   voteReq.SubID,
        "value":    voteReq.Value,
    }
    jsonData, err := json.Marshal(catAPIRequestBody)
    if err != nil {
        sendJSONResponse(http.StatusInternalServerError, map[string]string{"error": "Failed to marshal JSON data"})
        return
    }

    req, err := http.NewRequest("POST", "https://api.thecatapi.com/v1/votes", bytes.NewBuffer(jsonData))
    if err != nil {
        sendJSONResponse(http.StatusInternalServerError, map[string]string{"error": "Failed to create request"})
        return
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-api-key", apiKey)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        sendJSONResponse(http.StatusInternalServerError, map[string]string{"error": "Failed to submit vote: " + err.Error()})
        return
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        sendJSONResponse(http.StatusInternalServerError, map[string]string{"error": "Failed to read response body"})
        return
    }

    // Log the response body for debugging
    //log.Printf("Response body: %s", string(body))

    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        sendJSONResponse(http.StatusInternalServerError, map[string]string{"error": "Failed to unmarshal response body"})
        return
    }

    // Check if the response contains an error message
    if errorMsg, ok := result["message"].(string); ok && strings.ToLower(errorMsg) != "success" {
        sendJSONResponse(http.StatusInternalServerError, map[string]string{"error": "API Error: " + string(body)})
        return
    }

    // If we've reached this point, it's a successful response
    sendJSONResponse(http.StatusOK, result)
}