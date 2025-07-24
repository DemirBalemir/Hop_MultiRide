# 🛴 MultiRide

MultiRide is a backend command-line application written in Go that calculates the optimal multi-hop travel path between electric scooters. The goal is to find a scooter route from a starting point to a destination, considering real-world constraints like battery percentage, elevation, travel time, and distance — even if it requires switching scooters along the way.

---

## 🚀 Features

- 🗺️ **Graph-based routing engine** that finds the best scooter-to-scooter path
- 🔄 **Multi-hop switching**: Finds routes that require one or more scooter changes
- ⚡ **Battery-aware routing**: Takes battery percentage into account to estimate max range
- 🧭 **Elevation-aware logic**: Applies elevation penalties to estimated range
- ⏱️ **Time-based discount support**: Can adjust final pricing based on duration
- 🧵 **Concurrency**: Graph building is done concurrently using goroutines, mutexes, and semaphores to scale route computation
- 🧪 **Full testing suite**: Includes table-driven unit tests and mocks for external services

---


---

## ⚙️ How It Works

### 1. Load scooter data from JSON  
Scooters are stored in a static `scooters.json` file with fields like `ID`, `Latitude`, `Longitude`, `Battery`, and optionally `Elevation`.

### 2. Build a graph  
Using concurrency (20 goroutines max), a graph is constructed between scooters:
- Each node is a scooter
- An edge exists if another scooter is reachable within range (after applying elevation penalty)

### 3. Find optimal path  
A custom routing algorithm based on Dijkstra’s principle selects the path with:
- Minimum number of scooter switches  
- Then minimum time  
- Destination (-1) is added when the final destination is reachable directly from the last scooter

---

## 🧪 Testing

The project includes robust unit tests in the `algorithm` package, featuring:
- ✅ Table-driven test cases for different routing scenarios
- ✅ Custom mock function to simulate API responses (distance, errors, etc.)
- ✅ Coverage for:
  - Direct reachability
  - Multi-hop routing
  - No available paths
  - External API failures

 ## CLI Flags

| Flag         | Description                                 |
| ------------ | ------------------------------------------- |
| `-generate`  | Generates random scooter data (300 entries) |
| `-elevation` | Updates scooter JSON with elevation values  |


 ## Tech Stack

Chi (for possible HTTP expansion)

OSRM (routing engine)

Google Maps Elevation API

sync.WaitGroup / sync.Mutex / semaphores for concurrency control

testing & testify for unit testing

## Example scooter data structure

{
  "ID": 42,
  "Latitude": 39.9408,
  "Longitude": 32.8541,
  "Battery": 75,
  "Elevation": 950.2
}
