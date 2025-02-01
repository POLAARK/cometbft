#!/usr/bin/env python3
import os
import glob
import pandas as pd
import matplotlib.pyplot as plt

EXPORT_DIR = "../monitoring/exported_data"

def main():
    # Gather all CSV files
    csv_files = glob.glob(os.path.join(EXPORT_DIR, "*.csv"))
    if not csv_files:
        print("No CSV files found in the export directory.")
        return

    # Concatenate all CSVs into one DataFrame
    df_list = []
    for file in csv_files:
        try:
            temp_df = pd.read_csv(file)
            df_list.append(temp_df)
        except Exception as e:
            print(f"Error reading {file}: {e}")

    if not df_list:
        print("No valid CSV data.")
        return

    df = pd.concat(df_list, ignore_index=True)

    # Summarize total bytes sent and received by (Threshold, Payload, Validators)
    group_cols = ["Threshold", "Payload", "Validators"]
    agg_df = df.groupby(group_cols, as_index=False)[["Bytes Sent", "Bytes Received"]].sum()
    # You can also do .mean() or other metrics if desired.

    # Create a combined column for (Payload, Validators) to pivot on
    agg_df["Config"] = agg_df.apply(
        lambda row: f"PL={row['Payload']}, VAL={row['Validators']}",
        axis=1
    )

    # We'll pivot so each threshold becomes a column, and each Config is a row
    pivot_df = agg_df.pivot(index="Config", columns="Threshold", values=["Bytes Sent", "Bytes Received"])
    # This gives a multi-level column structure: ("Bytes Sent", 50), ("Bytes Sent", 75)...

    # Sort row index for consistency
    pivot_df.sort_index(inplace=True, key=lambda idx: idx.str.extract(r'(\d+)', expand=False).astype(int))

    # We'll do a grouped bar chart for Bytes Sent + Bytes Received
    # (one cluster per "Config", multiple bars for thresholds)
    fig, ax = plt.subplots(figsize=(10, 6))

    # We want each "Config" on the X-axis, with bars for each threshold.
    # We'll flatten the pivot columns for easier referencing in a loop.
    # pivot_df.columns looks like MultiIndex: [("Bytes Sent", 50), ("Bytes Sent", 75), ...].
    # Let's convert that to something simpler.
    thresholds = pivot_df["Bytes Sent"].columns  # e.g. [50, 75, 90, 100] if present
    x_labels = pivot_df.index  # e.g. "PL=100, VAL=4", etc.

    x = range(len(x_labels))     # numeric positions of each config cluster
    bar_width = 0.1             # width of each bar
    offset_shift = 0

    for thresh in thresholds:
        # For each threshold, we get the 'Bytes Sent' and 'Bytes Received'
        bytes_sent = pivot_df["Bytes Sent"][thresh]
        bytes_recv = pivot_df["Bytes Received"][thresh]

        # If either series has NaN for missing combos, fill with 0
        bytes_sent = bytes_sent.fillna(0)
        bytes_recv = bytes_recv.fillna(0)

        # We'll sum them to see total bandwidth in one bar
        total_bandwidth = bytes_sent + bytes_recv

        # Plot them at x + offset_shift
        ax.bar(
            [pos + offset_shift for pos in x],
            total_bandwidth,
            width=bar_width,
            label=f"Threshold={thresh}"
        )
        offset_shift += bar_width

    # Formatting the axes
    ax.set_xticks([pos + bar_width*(len(thresholds)/2) for pos in x])  # center the labels
    ax.set_xticklabels(x_labels, rotation=45, ha="right")
    ax.set_ylabel("Total Bytes (Sent + Received)")
    ax.set_title("Bandwidth Comparison by Threshold (Grouped by Payload & Validators)")
    ax.legend()
    plt.tight_layout()
    plt.show()


if __name__ == "__main__":
    main()

