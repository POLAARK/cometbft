#!/usr/bin/env python3

import os
import glob
import pandas as pd
import matplotlib.pyplot as plt

EXPORT_DIR = "../monitoring/exported_data"

def main():
    # Récupère tous les CSV du répertoire
    all_files = glob.glob(os.path.join(EXPORT_DIR, "*.csv"))
    if not all_files:
        print("Aucun fichier CSV trouvé dans le répertoire export.")
        return

    # Charge et concatène tous les CSV dans un seul DataFrame
    df_list = []
    for file in all_files:
        try:
            df_temp = pd.read_csv(file, parse_dates=["Timestamp"])
            df_list.append(df_temp)
        except Exception as e:
            print(f"Impossible de lire {file} : {e}")

    if not df_list:
        print("Aucune donnée CSV valide n'a pu être lue.")
        return

    df = pd.concat(df_list, ignore_index=True)

    #------------------------------
    # 1) Impact du nombre de validateurs sur la bande passante
    #------------------------------
    validators_agg = df.groupby("Validators", as_index=False)[["Bytes Sent", "Bytes Received"]].sum()
    validators_agg.sort_values(by="Validators", inplace=True)

    plt.figure(figsize=(6, 4))
    x_vals = validators_agg["Validators"]
    plt.bar(x_vals, validators_agg["Bytes Sent"], width=0.4, label="Bytes Sent")
    # On empile Bytes Received au-dessus de Bytes Sent pour une vue "stacked"
    plt.bar(x_vals, validators_agg["Bytes Received"], width=0.4, bottom=validators_agg["Bytes Sent"], label="Bytes Received")
    plt.title("Impact du nombre de validateurs sur la bande passante")
    plt.xlabel("Nombre de validateurs")
    plt.ylabel("Total Bytes (Sent + Received)")
    plt.legend()
    plt.tight_layout()
    plt.show()

    #------------------------------
    # 2) Impact de la taille du payload sur la bande passante
    #------------------------------
    payload_agg = df.groupby("Payload", as_index=False)[["Bytes Sent", "Bytes Received"]].sum()
    payload_agg.sort_values(by="Payload", inplace=True)

    plt.figure(figsize=(6, 4))
    x_vals = payload_agg["Payload"]
    # On peut ajuster la largeur de barre (width) pour mieux visualiser si besoin
    plt.bar(x_vals, payload_agg["Bytes Sent"], width=5, label="Bytes Sent")
    plt.bar(x_vals, payload_agg["Bytes Received"], width=5, bottom=payload_agg["Bytes Sent"], label="Bytes Received")
    plt.title("Impact de la taille du payload sur la bande passante")
    plt.xlabel("Taille du payload (octets)")
    plt.ylabel("Total Bytes (Sent + Received)")
    plt.legend()
    plt.tight_layout()
    plt.show()

    #------------------------------
    # 3) Impact du seuil (threshold) sur la bande passante
    #------------------------------
    threshold_agg = df.groupby("Threshold", as_index=False)[["Bytes Sent", "Bytes Received"]].sum()
    threshold_agg.sort_values(by="Threshold", inplace=True)

    plt.figure(figsize=(6, 4))
    x_vals = threshold_agg["Threshold"]
    plt.bar(x_vals, threshold_agg["Bytes Sent"], width=3, label="Bytes Sent")
    plt.bar(x_vals, threshold_agg["Bytes Received"], width=3, bottom=threshold_agg["Bytes Sent"], label="Bytes Received")
    plt.title("Impact du seuil (Threshold) sur la bande passante")
    plt.xlabel("Seuil du mempool (%)")
    plt.ylabel("Total Bytes (Sent + Received)")
    plt.legend()
    plt.tight_layout()
    plt.show()

if __name__ == "__main__":
    main()
