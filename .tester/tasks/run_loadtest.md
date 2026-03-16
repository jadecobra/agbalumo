# Distributed Load Testing with k6

This guide explains how to use the existing `benchmark_users.js` script to run a massive 100k user load test without crashing the host machine, by utilizing a "poor man's botnet" approach across multiple old machines (like Debian laptops).

The `run_loadtest.sh` script currently seeds 100k database entries and starts the server. The actual `k6 run` line inside of it has been commented out to facilitate running the clients remotely.

## 1. Setup the Server (Main Host)

1. Find the local network IP address of the machine hosting the server (e.g., `192.168.1.50`).
2. Execute the setup script on the host to seed the database and spin up the server:
   ```bash
   ./scripts/run_loadtest.sh
   ```
   *(Note: The server will start in the background and wait. It will not run the local benchmark because the `k6` command is commented out.)*

## 2. Setup the Clients (Debian Laptops)

On every machine that will generate traffic, do the following:

1. **Install k6**:
   ```bash
   sudo gpg -k
   sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
   echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
   sudo apt-get update
   sudo apt-get install k6
   ```

2. **Copy the Script**:
   Transfer `scripts/benchmark_users.js` from the main host to each laptop.

3. **Increase File Descriptors**:
   Increase the OS file descriptor limit so the network stack can handle the massive concurrency:
   ```bash
   ulimit -n 65000
   ```

## 3. Execute the Distributed Test

k6 handles distributed load generation using "Execution Segments". You tell each node exactly what fraction of the total load it should handle. No complex networking or master-slave configuration is needed.

Assuming you have **4 laptops**, you will split the 100k users evenly (25k Virtual Users per laptop).

Run these commands roughly simultaneously on each respective laptop, replacing `192.168.1.50` with your server's actual IP address:

**Laptop 1 (0% to 25% of load):**
```bash
TARGET_URL="https://192.168.1.50:8443" k6 run --insecure-skip-tls-verify --execution-segment "0:1/4" --execution-segment-sequence "0,1/4,2/4,3/4,1" benchmark_users.js
```

**Laptop 2 (25% to 50% of load):**
```bash
TARGET_URL="https://192.168.1.50:8443" k6 run --insecure-skip-tls-verify --execution-segment "1/4:2/4" --execution-segment-sequence "0,1/4,2/4,3/4,1" benchmark_users.js
```

**Laptop 3 (50% to 75% of load):**
```bash
TARGET_URL="https://192.168.1.50:8443" k6 run --insecure-skip-tls-verify --execution-segment "2/4:3/4" --execution-segment-sequence "0,1/4,2/4,3/4,1" benchmark_users.js
```

**Laptop 4 (75% to 100% of load):**
```bash
TARGET_URL="https://192.168.1.50:8443" k6 run --insecure-skip-tls-verify --execution-segment "3/4:1" --execution-segment-sequence "0,1/4,2/4,3/4,1" benchmark_users.js
```

### Note on Cleanup
Because you commented out the local k6 execution in `run_loadtest.sh`, the script on the host will just hang until you kill it or the laptops finish. Once the benchmark completes, return to the host machine and run:
```bash
kill $(pgrep -f "tmp_harness serve")
rm tmp_harness
```
