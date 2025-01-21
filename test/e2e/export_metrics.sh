#!/bin/bash
PROMETHEUS_URL="http://localhost:9090"
EXPORT_DIR="./monitoring/exported_data"
METRICS=("cometbft_mempool_tx_size_bytes_sum" "cometbft_mempool_tx_size_bytes_count")

mkdir -p "$EXPORT_DIR"

OUTPUT_CSV="${EXPORT_DIR}/transactions_vs_size.csv"

if ! command -v jq &>/dev/null; then
	echo "Error: jq is not installed. Install jq to continue."
	exit 1
fi

export_metrics() {
	# Determine time range based on OS
	if [[ "$OSTYPE" == "darwin"* ]]; then
		START=$(date -v-1H -u +"%Y-%m-%dT%H:%M:%SZ")
		END=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
	else
		START=$(date --utc -d "1 hour ago" +"%Y-%m-%dT%H:%M:%SZ")
		END=$(date --utc +"%Y-%m-%dT%H:%M:%SZ")
	fi

	STEP="1s"

	# Initialize temporary files for each metric
	TMP_SUM="${EXPORT_DIR}/metric_sum.tmp"
	TMP_COUNT="${EXPORT_DIR}/metric_count.tmp"

	for METRIC in "${METRICS[@]}"; do
		echo "Exporting metric: $METRIC"

		# Fetch data from Prometheus
		if ! curl -G "${PROMETHEUS_URL}/api/v1/query_range" \
			--data-urlencode "query=${METRIC}" \
			--data-urlencode "start=${START}" \
			--data-urlencode "end=${END}" \
			--data-urlencode "step=${STEP}" \
			-o "${EXPORT_DIR}/${METRIC}.json" \
			--fail --silent; then
			echo "Error: Could not connect to Prometheus at ${PROMETHEUS_URL}"
			echo "Please ensure Prometheus is running before executing this script"
			exit 1
		fi

		# Parse data and save to temporary file
		if [[ "$METRIC" == "cometbft_mempool_tx_size_bytes_sum" ]]; then
			cat "${EXPORT_DIR}/${METRIC}.json" | jq -r '
            .data.result[] |
            .metric as $metric |
            .values[] | [.[0], $metric.job, .[1]] | @tsv' >"$TMP_SUM"
		elif [[ "$METRIC" == "cometbft_mempool_tx_size_bytes_count" ]]; then
			cat "${EXPORT_DIR}/${METRIC}.json" | jq -r '
            .data.result[] |
            .metric as $metric |
            .values[] | [.[0], $metric.job, .[1]] | @tsv' >"$TMP_COUNT"
		fi
	done

	# Write the header to the output CSV
	echo "timestamp,validator,tx_count,tx_size" >"$OUTPUT_CSV"

	# Combine the data from both metrics
	join -t$'\t' -o 1.1,1.2,1.3,2.3 "$TMP_COUNT" "$TMP_SUM" | while IFS=$'\t' read -r timestamp validator tx_count tx_size; do
		echo "$timestamp,$validator,$tx_count,$tx_size" >>"$OUTPUT_CSV"
	done

	# Clean up temporary files
	rm -f "$TMP_SUM" "$TMP_COUNT"

	echo "Combined metrics saved to $OUTPUT_CSV"
}

export_metrics

echo "All metrics exported and combined successfully!"
