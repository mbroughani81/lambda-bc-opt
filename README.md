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


## Dependencies

- **wrk2**: A powerful HTTP benchmarking tool that supports constant request rates.
- **Python**: The benchmarking Python notebook for running benchmark, parsing outputs, generating plots, and exporting data.
- **OpenWhisk**: An open-source serverless platform where the functions (handlers) are deployed.

## How to work

1. Ensure you have all necessary dependencies installed (`wrk2`, Python, `zip`, OpenWhisk CLI, etc.).
2. Use the `Makefile` to package and deploy the OpenWhisk actions.
   ```bash
   make devel
   ```
3. Code cells in `bench.py` notebook will run the tests and generate results. After that Check the generated CSV files and plots for benchmarking results.

## License
MIT License

## Authors
[Your Name]
