import os
import subprocess
import requests
import time
import csv
from datetime import datetime

PROMETHEUS_URL = "http://localhost:9090/api/v1/query_range"
EXPORT_DIR = "test/e2e/monitoring/exported_data2"

def check_existing_results(payload, validators, threshold, load_time):
    """Check if a test has already been completed based on the presence of its CSV file."""
    filename = f"{EXPORT_DIR}/metrics_payload{payload}_validators{validators}_threshold{threshold}_loadtime{load_time}.csv"
    return os.path.isfile(filename)


def get_metrics(start_time, end_time, step="1s", validator_indices=None):
    """Fetch and aggregate metrics from Prometheus over a specified time range."""
    queries = {
        "bytes_sent": "cometbft_mempool_bytes_sent",
        "bytes_received": "cometbft_mempool_bytes_received",
        "transactions_sent": "cometbft_mempool_transactions_sent",
        "transactions_received": "cometbft_mempool_transactions_received",
        "signatures_received": "cometbft_mempool_signatures_sent_size",
        "signatures_sent": "cometbft_mempool_signatures_received_size",
        "total_transactions_in_consensus": "cometbft_consensus_total_txs"
    }

    all_metrics = {key: {} for key in queries}
    timestamps = []

    label_filter = ""
    if validator_indices is not None and len(validator_indices) > 0:
        pattern = "|".join(f"validator{i:02}" for i in validator_indices)
        label_filter = f'{{job=~"{pattern}"}}'

    for key, query in queries.items():
        if label_filter:
           query = f"{query}{label_filter}"

        response = requests.get(PROMETHEUS_URL, params={
            "query": query,
            "start": float(start_time),
            "end": float(end_time),
            "step": step
        })

        json_response = response.json()

        if response.status_code == 200 and "data" in json_response:
            results = json_response["data"].get("result", [])

            if results:
                for result in results:
                    for value in result["values"]:
                        timestamp, metric_value = value
                        timestamp = datetime.utcfromtimestamp(float(timestamp)).strftime('%Y-%m-%d %H:%M:%S')

                        if timestamp not in all_metrics[key]:
                            all_metrics[key][timestamp] = []

                        all_metrics[key][timestamp].append(float(metric_value))

                        if timestamp not in timestamps:
                            timestamps.append(timestamp)
            else:
                print(f"⚠️ No data found for query: {query}.")

    # Compute mean for each metric across all validators
    for key in all_metrics:
        for timestamp in all_metrics[key]:
            values = all_metrics[key][timestamp]
            all_metrics[key][timestamp] = sum(values) / len(values) if values else 0  # Compute mean

    return timestamps, all_metrics


def update_network_config(payload_size, num_validators):
    """Update the network configuration file with the payload size and number of validators."""
    load_tx_batch_size = 5

    print(f"Setting load_tx_size_bytes = {payload_size}")

    with open("test/e2e/networks/simple.toml", "w") as f:
        f.write(f"log_level = \"info\"\n")
        f.write(f"load_tx_size_bytes = {payload_size}\n")
        f.write(f"load_tx_batch_size = {load_tx_batch_size}\n")
        f.write("prometheus = true\n")
        for i in range(num_validators):
            f.write(f"[node.validator{i:02}]\n")


def run_tests(payload_size, num_validators, threshold, load_time):
    """Run the tests by updating the config and executing the workflow."""
    update_network_config(payload_size, num_validators)

    env = os.environ.copy()
    env["MEMPOOL_THRESHOLD_PERCENT"] = str(threshold)

    print(f"✅ Setting MEMPOOL_THRESHOLD_PERCENT = {threshold}")

    # Build and prepare the testnet
    subprocess.run(["make", "build"], check=True, env=env)  # Pass the modified environment

    # Change working directory for test execution
    subprocess.run(["make", "fast"], cwd="test/e2e", check=True, env=env)

    # Start network
    try:
        subprocess.run(["./build/runner", "-f", "networks/simple.toml", "cleanup"], cwd="test/e2e", check=True)
    except:
        print("Everything is already cleaned")


    subprocess.run(["./build/runner", "-f", "networks/simple.toml", "start"], cwd="test/e2e", check=True, env=env)

    # Start capturing time before load begins
    start_time = time.time() + 5

    # Load transactions with time parameter
    subprocess.run(["./build/runner", "-f", "networks/simple.toml", "load", "--time", str(load_time)], cwd="test/e2e", check=True)

    # Capture end time after load completes
    end_time = time.time() +2

    return start_time, end_time


def save_metrics_to_csv(payload, validators, threshold, load_time, timestamps, all_metrics):
    """Save aggregated (mean) time-series metrics to CSV in the correct order."""
    directory = EXPORT_DIR

    # Ensure the directory exists
    os.makedirs(directory, exist_ok=True)

    # Define the file path correctly
    filename = f"{directory}/metrics_payload{payload}_validators{validators}_threshold{threshold}_loadtime{load_time}.csv"

    # Ensure header is only written once
    write_header = not os.path.isfile(filename)

    with open(filename, mode="w" if write_header else "a", newline="") as file:
        writer = csv.writer(file)
        if write_header:
            writer.writerow(["Timestamp", "Payload", "Validators", "Threshold", "Load Time",
                            "Bytes Sent", "Bytes Received", "Tx Sent", "Tx Received",
                            "Signatures Received", "Signatures Sent", "Total Tx In Consensus"])

        for timestamp in timestamps:
            writer.writerow([
                timestamp,
                payload,
                validators,
                threshold,
                load_time,
                all_metrics["bytes_sent"].get(timestamp, 0),
                all_metrics["bytes_received"].get(timestamp, 0),
                all_metrics["transactions_sent"].get(timestamp, 0),
                all_metrics["transactions_received"].get(timestamp, 0),
                all_metrics["signatures_received"].get(timestamp, 0),
                all_metrics["signatures_sent"].get(timestamp, 0),
                all_metrics["total_transactions_in_consensus"].get(timestamp, 0)
            ])


if __name__ == "__main__":
    subprocess.run(["./build/runner", "-f", "networks/simple.toml", "monitor", "stop"], cwd="test/e2e", check=True)
    subprocess.run(["./build/runner", "-f", "networks/simple.toml", "monitor", "start"], cwd="test/e2e", check=True)
    for payload in [120]:  # Varying payload sizes
        for validators in [16]:  # Different validator configurations
            for threshold in [100]:  # Different mempool thresholds
                for load_time in [10]:  # Different load durations
                    print(f"Running test: Payload={payload}, Validators={validators}, Threshold={threshold}, Time={load_time}")

                    # Check if results already exist
                    if check_existing_results(payload, validators, threshold, load_time):
                        print(f"⚠️ Test already completed, skipping: Payload={payload}, Validators={validators}, Threshold={threshold}, Time={load_time}")
                        continue

                    start_time, end_time = run_tests(payload, validators, threshold, load_time)

                    # Ensure Prometheus collects data before querying
                    time.sleep(4)  # Short wait for Prometheus to process

                    validator_indices = list(range(validators))
                    timestamps, all_metrics = get_metrics(
                        start_time,
                        end_time,
                        step="1s",
                        validator_indices=validator_indices
                    )


                    print(f"Saving data for Payload={payload}, Validators={validators}, Threshold={threshold}, Time={load_time}")
                    save_metrics_to_csv(payload, validators, threshold, load_time, timestamps, all_metrics)

                    # subprocess.run(["./build/runner", "-f", "networks/simple.toml", "monitor", "stop"], cwd="test/e2e", check=True)
                    subprocess.run(["./build/runner", "-f", "networks/simple.toml", "stop"], cwd="test/e2e", check=True)
                    try:
                        subprocess.run(["./build/runner", "-f", "networks/simple.toml", "cleanup"], cwd="test/e2e", check=True)
                        time.sleep(2)
                    except:
                        print("Everything is already cleaned")