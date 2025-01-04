#!/bin/bash
PROMETHEUS_URL="http://localhost:9090"
EXPORT_DIR=".monitoring/exported_data"
METRICS=("cometbft_mempool_tx_size_bytes_sum" "cometbft_mempool_tx_size_bytes_bucket")

mkdir -p "$EXPORT_DIR"

export_metric() {
    METRIC=$1

    if [[ "$OSTYPE" == "darwin"* ]]; then
        START=$(date -v-1H -u +"%Y-%m-%dT%H:%M:%SZ")
        END=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    else
        START=$(date --utc -d "1 hour ago" +"%Y-%m-%dT%H:%M:%SZ")
        END=$(date --utc +"%Y-%m-%dT%H:%M:%SZ")
    fi

    STEP="1s"

    echo "Exporting metric: $METRIC"

    curl -G "${PROMETHEUS_URL}/api/v1/query_range" \
        --data-urlencode "query=${METRIC}" \
        --data-urlencode "start=${START}" \
        --data-urlencode "end=${END}" \
        --data-urlencode "step=${STEP}" \
        -o "${EXPORT_DIR}/${METRIC}.json"

    if ! command -v jq &> /dev/null; then
        echo "Error: jq is not installed. Install jq to continue."
        exit 1
    fi

    cat "${EXPORT_DIR}/${METRIC}.json" | jq -r '
        .data.result[] |
        .metric as $metric |
        .values[] | [$metric.__name__, .[0], .[1]] | @csv' > "${EXPORT_DIR}/${METRIC}.csv"

    echo "Saved to ${EXPORT_DIR}/${METRIC}.csv"
}

for METRIC in "${METRICS[@]}"; do
    export_metric "$METRIC"
done

echo "All metrics exported successfully!"
