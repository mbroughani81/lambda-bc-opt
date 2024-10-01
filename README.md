# Lambda Batch-Call Optimization

## Project Structure

### Handlers and Benchmarking
- **lambda-bc-opt/handlers/wsk**: Contains handler functions written in Go that are executed on OpenWhisk. Each subdirectory represents a different function.
- **benchmark**: Contains benchmarking scripts and related data, including Lua scripts for `wrk`, CSV files for benchmarking results, and Python scripts for running the tests and parsing the output.

### Important Scripts

1. **bench.py**:
   - This Python script orchestrates the benchmarking process. It runs the `wrk2` tool to simulate different request rates (RPS) against specific URLs (usually pointing to OpenWhisk functions or services).
   - It collects latency percentiles (50th, 90th, and 99th) and exports the results into CSV files.
   - The script also generates plots of latency percentiles against RPS and can read from CSV files to generate additional insights.

2. **visitorcounter_request.lua**:
   - A Lua script used by `wrk2` to simulate requests against the `visitorCounter` service.

3. **Makefile**:
   - The Makefile provides automation for packaging and deploying OpenWhisk actions.
   - The `devel` target automates the creation of zip files for each handler in the `lambda-bc-opt/handlers/wsk` directory and updates the corresponding OpenWhisk actions using Docker-based Golang actions.
   - The `zip` target creates individual zip files from handler subdirectories.
   - The `clean` target cleans up the generated zip directories.

### OpenWhisk Integration

- **OpenWhisk Functions**: Each handler subdirectory under `lambda-bc-opt/handlers/wsk` represents an OpenWhisk action written in Go. These functions are zipped and deployed using the Makefile.
- **OpenWhisk Compiler**: The OpenWhisk functions are compiled and deployed using a Docker image for Go actions, with memory settings of 1024MB by default.

### Benchmarking

- The `wrk2` tool is used to simulate requests at varying rates (RPS) and measure latency percentiles.
- Latency data is collected and exported into CSV files, and graphs are generated to visualize the relationship between RPS and latency percentiles.
- Different sets of tests are conducted for Redis (e.g., `redis-batched` and `redis-naive`) and OpenWhisk actions (e.g., `gencnt1` and `gencnt2`).

## Dependencies

- **wrk2**: A powerful HTTP benchmarking tool that supports constant request rates.
- **Python**: The benchmarking script uses Python for running the tests, parsing outputs, generating plots, and exporting data.
- **OpenWhisk**: An open-source serverless platform where the functions (handlers) are deployed.

## How to Run the Project

1. Ensure you have all necessary dependencies installed (`wrk2`, Python, `zip`, OpenWhisk CLI, etc.).
2. Use the `Makefile` to package and deploy the OpenWhisk actions.
   ```bash
   make devel
   ```
3. Run the benchmarking script `bench.py` to conduct tests and generate results.
   ```bash
   python3 benchmark/bench.py
   ```
4. Check the generated CSV files and plots for benchmarking results.

## License
MIT License

## Authors
[Your Name]
