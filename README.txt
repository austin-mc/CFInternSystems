Cloudflare Internship Systems Assignment
Created by: Austin Christiansen

Endpoints: 
Root ("/"): Returns a 404 Not Found response
README ("/README.txt"): Returns the contents of this README in application/text format
Stats ("/stats): Queries the Cloudflare Radar API at cfisysapi.developers.workers.dev/stats to retrieve data. Calculates the Mean, Median, Min, and Max values from this data and returns them as strings in a JSON object. A query parameter "timestamp" can be provided to get Radar data from this timestamp instead of from the current Epoch time.     

Testing:
Run the following command to run unit tests:
go test