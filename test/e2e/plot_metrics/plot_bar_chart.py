import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

# File path to the CSV
CSV_FILE = '../monitoring/exported_data/transactions_vs_size.csv'

def generate_barchart(csv_file):
    # Load the CSV into a DataFrame
    try:
        df = pd.read_csv(csv_file)
    except FileNotFoundError:
        print(f"Error: File {csv_file} not found. Ensure the file exists and the path is correct.")
        return

    # Ensure the required columns are present
    required_columns = {'validator', 'tx_count', 'tx_size'}
    if not required_columns.issubset(df.columns):
        print(f"Error: CSV file must contain columns {required_columns}.")
        return

    # Convert `tx_count` and `tx_size` to numeric values
    df['tx_count'] = pd.to_numeric(df['tx_count'])
    df['tx_size'] = pd.to_numeric(df['tx_size'])

    # Plot the bar chart
    plt.figure(figsize=(12, 6))
    sns.barplot(x='tx_count', y='tx_size', hue='validator', data=df)

    # Set labels and title
    plt.xlabel('Transaction Count')
    plt.ylabel('Transaction Size (Bytes)')
    plt.title('Transaction Size by Validator for Each Transaction Count')
    plt.legend(title='Validator')
    plt.grid(axis='y', linestyle='--', alpha=0.7)

    # Show the plot
    plt.tight_layout()
    plt.show()

# Call the function
generate_barchart(CSV_FILE)
