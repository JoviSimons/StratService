# StratService documentation

## API Endpoints

### Endpoint: Get All Strategies
- URL: /all
- Method: GET
- Description: Retrieves all strategies stored in the database.
- Handler Function: returnAll

Example Response:
```
[
  {
    "_id": "645a46d8d9335b1665d85ef2",
    "created": "2023-05-09T15:12:56.67+02:00",
    "ex": "test",
    "name": "test.ex5",
    "userid": "n1BOnD9Ip2hk7Ua6XG83vrkqqxk1"
  },
  {
    "_id": "645a46d8d9335b1665d85ef3",
    "created": "2023-06-09T15:12:56.67+02:00",
    "ex": "test1",
    "name": "test1.ex5",
    "userid": "n1BOnD9Ip2hk7Ua6XG83vrkqqxk1"
  },
]
```

### Endpoint: Get Strategy by ID
- URL: /get/{id}
- Method: GET
- Description: Retrieves a specific strategy based on the provided ID.
- Handler Function: returnStrat
- Path Parameters:
id (string): The ID of the strategy to retrieve.

Example Response:
```
  {
    "_id": "645a46d8d9335b1665d85ef2",
    "created": "2023-05-09T15:12:56.67+02:00",
    "ex": "test",
    "name": "test.ex5",
    "userid": "n1BOnD9Ip2hk7Ua6XG83vrkqqxk1"
  }
```
### Endpoint: Use Strategy
- URL: /use/{id}
- Method: GET
- Description: Performs an action using a specific strategy based on the provided ID.
- Handler Function: useStrat
- Path Parameters:
id (string): The ID of the strategy to use.
- Action:
The strategy with the specified ID is retrieved and assigned to a Strategy object.
The name of the strategy is printed.
- Optional: The strategy can be sent as a message using RabbitMQ.

Example Response: None (prints the name of the strategy)

### Endpoint: Create Strategy
- URL: /create
- Method: POST
- Description: Stores a new strategy in the database and sends it to the TestManager.
- Handler Function: storeStrat
- Request Body: JSON object representing the strategy.

Example Request Body:
```
{
    "id": "test",
    "userId": "userID",
    "name": "ExpertMACD",
    "mq": "strat.mq5",
    "ex": "strat.ex"
}
```

Example Response: The ID of the newly created strategy.

## Running app

Start the server by running the Go application in the terminal.

```
go run .
``` 
Remember to adjust the host and port (http.ListenAndServe(":10000", myRouter)) based on your deployment environment or preferred configuration.

## Code documentation

`Main()`
```
The main function is the entry point of the Go program. 
It invokes the handleRequests function to start the HTTP server and handle incoming requests.
```

`handleRequest()`
```
The handleRequests function sets up the Gorilla Mux router, 
defines the HTTP endpoints, 
and starts the HTTP server to handle incoming requests.
```

`returnAll()`
```
The returnAll function is an HTTP request handler that returns all 
the strategies stored in the database as a JSON response.
```

`returnStrat()`
```
The returnStrat function is an HTTP request handler that returns a specific 
strategy based on the provided ID as a JSON response.
```

`useStrat()`
```
The useStrat function is an HTTP request handler that retrieves a specific strategy based 
on the provided ID and performs some action with it. 
In this case, it prints the name of the strategy and can potentially send a message using RabbitMQ.
```

`storeStrat()`
```
The storeStrat function is an HTTP request handler that stores a new strategy in the database. 
It receives the strategy data as a JSON payload in the request body, 
decodes it into a Strategy struct, inserts the strategy into the database, 
and sends the strategy to the TestManager using RabbitMQ.
```

`getAllStrats()`
```
The getAllStrats function retrieves all the strategies from the database 
and returns them as an array of primitive.M values.
```

`CORS()`
```
The CORS function is a middleware that adds Cross-Origin Resource Sharing (CORS) headers to the HTTP response. 
It allows cross-origin requests from any origin (*), 
supports all HTTP methods, and handles the preflight OPTIONS requests.
```

`FailOnError()`
```
The FailOnError function is a utility function that panics with an error message
if the provided error is not nil.
```