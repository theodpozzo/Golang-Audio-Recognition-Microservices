# Golang Audio Recognition Microservices


A modular Go microservice ecosystem for recognising audio fragments, querying an external music-identification API, and retrieving or storing track audio locally.  
The system is composed of three independent services - Cooltown, Search, and Tracks - that communicate through simple REST APIs to identify songs from base64 encoded audio.


## Overview


This project implements an audio recognition pipeline using a microservice architecture in Go. When given a base64 encoded audio fragment:


1. Cooltown accepts the fragment and orchestrates the lookup.  
2. Search queries the external Audd.io API to identify the track.  
3. Tracks retrieves or persists the audio in a local SQLite database.  
4. Cooltown returns the resolved audio back to the client.


Each service is isolated, maintainable, and responsible for a single part of the workflow.


## Architecture


```
Client
   ↓ POST /cooltown
Cooltown Service (3002)
   ↓ POST /search
Search Service (3001) -> Audd.io API
   ↓ GET /tracks/{id}
Tracks Service (3000) -> SQLite DB
```


## Services


### 1. Cooltown Service (Port 3002)


**Endpoint**


- `POST /cooltown`  
  **Input**
  ```json
  { "Audio": "<base64>" }
  ```
  **Output**
  ```json
  { "Audio": "<base64>" }
  ```


**Responsibilities**


- Validate incoming audio payloads.
- Call the Search service to resolve track ID or title.
- Fetch full audio from the Tracks service.
- Return appropriate HTTP status codes: 200, 400, 404, 500.


### 2. Search Service (Port 3001)


Interacts with the Audd.io API to recognise the audio fragment.


**Endpoint**


- `POST /search`  
  **Output**
  ```json
  { "Id": "<song-title>" }
  ```


**Responsibilities**


- Read API token from `search/api_key.txt`.
- Post audio + token to Audd.io and parse the response.
- Map Audd.io error states to sensible HTTP responses:
  - Invalid base64 -> 404
  - Invalid API key -> 400
  - Other failures -> 400 or 500
- Normalise track names by replacing spaces with `+` for use as an identifier.


### 3. Tracks Service (Port 3000)


Local storage for full track audio using SQLite.


**Endpoints**


- `PUT /tracks/{id}` - Insert or update a track
- `GET /tracks/{id}` - Retrieve audio
- `GET /tracks` - List all track IDs
- `DELETE /tracks/{id}` - Delete a track


**Database Schema**


| Column | Type |
|--------|------|
| Id     | TEXT |
| Audio  | TEXT |


**Responsibilities**


- Initialise SQLite database on startup.
- Provide CRUD operations using prepared statements.
- Return appropriate HTTP codes for success and failure.


## Key Features


- Decoupled Go microservice architecture.
- REST based HTTP communication between services.
- SQLite backed track storage with prepared SQL statements.
- External integration with Audd.io for audio fingerprinting.
- Consistent error semantics propagated across services.
- Router implementation using `gorilla/mux`.


## Running the System


### Prerequisites


- Go 1.20 or later
- SQLite3
- Audd.io API token saved at:
```
search/api_key.txt
```


### Start Each Service


Open three terminals and run:


```bash
go run ./cooltown
```


```bash
go run ./search
```


```bash
go run ./tracks
```


Each service prints a simple startup message and listens on its configured port.


## API Usage


### Identify and Retrieve Audio


**POST** `http://localhost:3002/cooltown`


**Request**
```json
{
  "Audio": "<base64-fragment>"
}
```


**Responses**


- `200 OK` - track audio returned
- `400 Bad Request` - invalid audio or malformed request
- `404 Not Found` - track not recognised or missing
- `500 Internal Server Error` - downstream service failure


### Tracks Service Examples


**Insert or Update Track**
```
PUT /tracks/My+Song
```
Body:
```json
{
  "Id": "My+Song",
  "Audio": "<base64-audio>"
}
```


**Retrieve Track**
```
GET /tracks/My+Song
```


**List All Tracks**
```
GET /tracks
```


**Delete Track**
```
DELETE /tracks/My+Song
```


## Directory Structure


```
cooltown/
    main.go
    cooltown/
        cooltown.go
search/
    main.go
    api_key.txt
    search/
        search.go
tracks/
    main.go
    tracks/
        tracks.go
    repository/
        repository.go
        track.go
    tmp/
        test.db
```


## Error Handling Notes


The code uses internal numeric return codes for some internal routines. The mapping used in this project is:


- 1 -> success
- 0 -> client error
- -1 -> not found
- -2 -> internal/server error


These values are mapped to HTTP status codes when producing responses. Consider replacing numeric codes with typed errors for clarity in future refactors.


## Suggested Improvements


- Add Dockerfiles and a Docker Compose file for unified local development.
- Introduce structured logging with a logging library such as Zap or Zerolog.
- Add retry and backoff logic for inter-service HTTP calls.
- Replace base64 payloads with streamed uploads or multipart form data for large audio fragments.
- Add unit tests and integration tests using Go's `httptest` package.
- Store richer metadata in the Tracks database - for example artist, album, and confidence score.
- Consider asynchronous processing using a message queue to improve resilience and throughput.


## Tests


- Add unit tests for each service handler.
- Add integration tests that start lightweight instances of each service and verify end-to-end flows.


## Security Considerations


- Store the Audd.io API token securely and avoid committing it to version control.
- Validate and limit input sizes to avoid memory exhaustion via large base64 payloads.
- Run the SQLite database with appropriate file permissions.
- Consider rate limiting and authentication on public endpoints.


## License


MIT
