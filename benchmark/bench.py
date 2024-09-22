# In[]:
import subprocess
import re
import matplotlib.pyplot as plt

# In[]:
# Function to run wrk and get the output
def run_wrk(rps, duration=30):
    """Run wrk2 for a specific RPS and return the latency data."""
    command = f"wrk -t10 -c20 -d{duration}s -R{rps} --latency -s visitorcounter_request.lua http://127.0.0.1:3001/visitorCounter"
    result = subprocess.run(command, shell=True, capture_output=True, text=True)
    return result.stdout

# In[]:
# Openwhisk run wrk
# Function to run wrk and get the output
def run_wrk(rps, duration=30):
    """Run wrk2 for a specific RPS and return the latency data."""
    command = f"wrk -t3 -c6 -d{duration}s -R{rps} --latency -s visitorcounter_request_openwhisk.lua http://10.10.0.1:3233"
    result = subprocess.run(command, shell=True, capture_output=True, text=True)
    return result.stdout

# In[]:
# Function to extract latency percentiles from wrk output
def parse_latency_output(output):
    """Extract 50th, 90th, and 99th percentile latencies from wrk output."""
    latencies = {}
    match_50 = re.search(r'50.000%\s+([0-9\.]+)([a-z]+)', output)
    match_90 = re.search(r'90.000%\s+([0-9\.]+)([a-z]+)', output)
    match_99 = re.search(r'99.000%\s+([0-9\.]+)([a-z]+)', output)

    # Convert to milliseconds
    if match_50:
        latencies['50th'] = convert_to_milliseconds(float(match_50.group(1)), match_50.group(2))
    if match_90:
        latencies['90th'] = convert_to_milliseconds(float(match_90.group(1)), match_90.group(2))
    if match_99:
        latencies['99th'] = convert_to_milliseconds(float(match_99.group(1)), match_99.group(2))

    return latencies

# Helper function to convert latency to milliseconds
def convert_to_milliseconds(value, unit):
    """Convert the latency values to milliseconds."""
    if unit == 'ms':
        return value
    elif unit == 's':
        return value * 1000
    elif unit == 'us':
        return value / 1000
    return value

# In[]:
# Define the RPS values to test
rps_values = [1]

# Lists to store results
latency_50th = []
latency_90th = []
latency_99th = []

# Loop through different RPS values, run wrk, and collect latencies
for rps in rps_values:
    print(f"Running wrk2 for {rps} requests per second...")
    output = run_wrk(rps)
    latencies = parse_latency_output(output)
    print(f"laaatt => {output}")
    print(f"50th percentile: {latencies.get('50th', 'N/A')} ms")
    print(f"90th percentile: {latencies.get('90th', 'N/A')} ms")
    print(f"99th percentile: {latencies.get('99th', 'N/A')} ms")
    # Append the results
    latency_50th.append(latencies.get('50th', None))
    latency_90th.append(latencies.get('90th', None))
    latency_99th.append(latencies.get('99th', None))

# In[]:
# Plot the latency percentiles
plt.figure(figsize=(10, 6))
plt.plot(rps_values, latency_50th, marker='o', label='50th Percentile')
plt.plot(rps_values, latency_90th, marker='o', label='90th Percentile')
plt.plot(rps_values, latency_99th, marker='o', label='99th Percentile')

plt.xlabel('Requests per Second (RPS)')
plt.ylabel('Latency (ms)')
plt.title('Latency Percentiles vs RPS')
plt.legend()
plt.grid(True)
plt.savefig('latency_plot.png')

# %%
