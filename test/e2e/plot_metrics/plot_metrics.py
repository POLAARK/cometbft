import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

df = pd.read_csv("../monitoring/exported_data/cometbft_mempool_tx_size_bytes_sum.csv")

df["timestamp"] = pd.to_datetime(df["timestamp"], unit="s")

mean_values = df.groupby("timestamp")["value"].mean().reset_index()
mean_values["job"] = "mean"
# sum_values = df.groupby("timestamp")["value"].sum().reset_index()
# sum_values["job"] = "sum"

combined_df = pd.concat([df, mean_values])
combined_df.head()
plt.figure(figsize=(12, 6))
sns.lineplot(data=combined_df, x="timestamp", y="value", hue="job", errorbar=None)

plt.title("Transaction Size Over Time by Validator")
plt.xlabel("Time")
plt.ylabel("Size (bytes)")

plt.xticks(rotation=45)

plt.tight_layout()

plt.show()
