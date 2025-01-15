import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

df = pd.read_csv(".monitoring/exported_data/cometbft_mempool_tx_size_bytes_bucket.csv")

df["timestamp"] = pd.to_datetime(df["timestamp"], unit="s")

plt.figure(figsize=(12, 6))
sns.lineplot(data=df, x="timestamp", y="value", hue="job", ci=None)

plt.title("Transaction Size Over Time by Validator")
plt.xlabel("Time")
plt.ylabel("Size (bytes)")

plt.xticks(rotation=45)

plt.tight_layout()

plt.show()
