#!/usr/bin/env python3
import os
import glob
import pandas as pd
import matplotlib.pyplot as plt
from matplotlib.patches import Patch
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
            last_row = temp_df.tail(1)
            df_list.append(last_row)
        except Exception as e:
            print(f"Error reading {file}: {e}")
    if not df_list:
        print("No valid CSV data.")
        return

    df = pd.concat(df_list, ignore_index=True)

    # Filter only configurations with Payload == 500 and Validators > 10
    df = df[(df["Payload"] == 250) & (df["Validators"] > 10)]

    # Create a combined Config column
    df["Config"] = df.apply(lambda row: f"PL={row['Payload']}, VAL={row['Validators']}", axis=1)

    # Group by Config and Threshold, summing Bytes Sent and Signatures Sent
    agg_df = df.groupby(["Config", "Threshold"], as_index=False)[["Bytes Sent", "Signatures Sent"]].sum()

    # Compute the non-signature part: Bytes Sent - Signatures Sent
    agg_df["NonSigSent"] = agg_df["Bytes Sent"] - agg_df["Signatures Sent"]

    # Pivot the data so that each threshold becomes a column.
    # We create two pivot tables: one for the non-signature bytes and one for the signature bytes.
    non_sig_pivot = agg_df.pivot(index="Config", columns="Threshold", values="NonSigSent")
    sig_pivot = agg_df.pivot(index="Config", columns="Threshold", values="Signatures Sent")

    # Sort the configs and thresholds
    configs = sorted(non_sig_pivot.index)
    x = range(len(configs))
    thresholds = sorted(non_sig_pivot.columns.astype(int))
    bar_width = 0.1

    fig, ax = plt.subplots(figsize=(12, 7))
    legend_handles = []
    for i, thresh in enumerate(thresholds):
        color = plt.cm.tab10(i)
        positions = [pos + i * bar_width for pos in x]
        non_sig = non_sig_pivot[thresh].reindex(configs).fillna(0)
        sig = sig_pivot[thresh].reindex(configs).fillna(0)
        ax.bar(positions, non_sig, width=bar_width, color=color)
        ax.bar(positions, sig, width=bar_width, color=color, hatch='//', bottom=non_sig)
        if thresh == 100:
            label = f"vanilla cometbft"
        else:
            label = f"Threshold={thresh}"
        legend_handles.append(Patch(facecolor=color,  label=label))


    signature_patch = Patch(facecolor='white', hatch='//', edgecolor='black', label=' Signature bytes portion')
    legend_handles.append(signature_patch)
    ax.legend(handles=legend_handles)
    # Adjust x-axis labels so that each config cluster is centered
    ax.set_xticks([pos + bar_width * (len(thresholds) / 2) for pos in x])
    ax.set_xticklabels(configs, rotation=45, ha="right")
    ax.set_ylabel("Bytes")
    ax.set_title("Message Composition: Signature Overhead vs. Other Bytes\n(Payload = 500, Validators > 10)")
    plt.tight_layout()
    plt.show()

if __name__ == "__main__":
    main()
