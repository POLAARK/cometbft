#!/usr/bin/env python3
import os
import glob
import pandas as pd
import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D  # noqa: F401
import numpy as np

EXPORT_DIR = "../monitoring/exported_data"

def main():
    # Exclude files with "threshold100" in their name
    csv_files = [f for f in glob.glob(os.path.join(EXPORT_DIR, "*.csv")) if "threshold100" not in f]
    if not csv_files:
        print("No CSV files found after exclusion.")
        return

    df_list = []
    for file in csv_files:
        try:
            temp_df = pd.read_csv(file)
            # Use only the last row of each CSV
            last_row = temp_df.tail(1).copy()
            last_row.loc[:, "Source"] = file
            df_list.append(last_row)
        except Exception as e:
            print(f"Error reading {file}: {e}")
    if not df_list:
        print("No valid CSV data.")
        return

    df = pd.concat(df_list, ignore_index=True)

    # Group by Validators and Threshold and average Total Tx In Consensus across different payloads
    grouped = df.groupby(["Validators", "Threshold"]).agg({"Total Tx In Consensus": "mean"}).reset_index()

    # Calculate Validation Percentage (100% when Total Tx In Consensus equals 200)
    grouped["Validation Percentage"] = (grouped["Total Tx In Consensus"] / 200) * 100

    # Prepare data for the 3D plot
    X = grouped["Validators"].values
    Y = grouped["Threshold"].values
    Z = grouped["Validation Percentage"].values

    fig = plt.figure(figsize=(10, 8))
    ax = fig.add_subplot(111, projection="3d")
    surf = ax.plot_trisurf(X, Y, Z, cmap="viridis", edgecolor="none")
    ax.set_xlabel("Validators")
    ax.set_ylabel("Threshold")
    ax.set_zlabel("Validation Percentage (%)")
    ax.set_title("Validation Percentage vs. Validators and Threshold\n(Averaged by Payload)")
    fig.colorbar(surf, ax=ax, label="Validation Percentage (%)")
    plt.tight_layout()
    plt.savefig("validation_percentage_3d_plot.png")
    plt.show()

if __name__ == "__main__":
    main()
