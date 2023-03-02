import argparse
import matplotlib.pyplot as plt
import csv

# Define the command-line arguments
parser = argparse.ArgumentParser()
parser.add_argument('--csv', help='CSV file to read')
args = parser.parse_args()

# Read the CSV file
with open(args.csv, newline='') as csvfile:
    reader = csv.reader(csvfile, delimiter=',', quotechar='"')
    rows = [row for row in reader]

# Extract the headers and data
headers = rows[0]
data = rows[1:]

# Create a quadrant chart
fig, ax = plt.subplots(figsize=(16, 16))
ax.set_title("Forensics RADAR: Automated DFIR Report")

# Draw the quadrants
for i in range(4):
  ax.axvline(x=0.5, color='black', linestyle='--')
  ax.axhline(y=0.5, color='black', linestyle='--')

# Add the keywords to each quadrant
for i in range(len(headers)):
  x = 0.25 if i % 2 == 0 else 0.75
  y = 0.75 if i // 2 == 0 else 0.25
  ax.text(x, y, headers[i], fontsize=12, fontweight='bold', ha='center', va='center')

  keyword_list = []
  for j in range(len(data)):
    keyword_list.append(data[j][i])
  ax.text(x, y - 0.1, '\n'.join(keyword_list), fontsize=10, ha='center', va='center')       

# Adjust the axis limits and remove the ticks
ax.set_xlim(0, 1)
ax.set_ylim(0, 1)
ax.set_xticks([])
ax.set_yticks([])

plt.show()