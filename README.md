# Forensics RADAR: Automated DFIR Reporting Tool for Real-time Insight into the Industry

> **Listening to our customers by analyzing what the world is saying.**


## Introduction

The field of Digital Forensics and Incident Response (DFIR) is rapidly evolving, and keeping up with the latest developments is essential for staying ahead of the curve. Traditional methods of information gathering, such as reading articles, newsletters, and customer feedback, can be time-consuming and unreliable. This is where our automated DFIR reporting tool comes in. Our tool uses Artificial Intelligence (AI) to extract keywords from multiple sources and generate quadrant-style reports that provide real-time insight into the industry.

## Objective

Our objective is to create an automated tool that collects data from various sources, extracts relevant keywords, categorizes them, and generates quadrant-style reports. The tool aims to help DFIR professionals stay up-to-date with the latest trends, news, and developments in the industry. With our tool, users will be able to quickly identify the emerging trends, predict future developments, and make informed decisions based on real-time data.

### System Design

![system](https://github.com/ozbekburak/dfir-radar/blob/main/img/system.png?raw=true)

## Methodology

Our tool is built using **Go** programming language and **ChatGPT** to extract keywords from different sources such as twitter, newsletters, articles, papers, etc. These keywords are then categorized based on their relevance to the DFIR domain. The categories are then used to generate quadrant-style reports, which provide real-time insights into the industry. The reports are generated using **Python** and its libraries, including matplot, and are saved as CSV files.

### NOTE

Our application currently works with ready-to-use .txt files, and it does not gather data directly from online sources at this time. Users will need to provide the necessary .txt files to the program for analysis. We are continually updating and improving our tool and may add online data gathering capabilities in the **near** future.


## Usage

User needs to set the [OPENAI API KEY](https://platform.openai.com/account/api-keys) environment variable. After that, they can generate the CSV file by running the following command in the terminal:

```go
export OPENAI_API_KEY={your_api_key}
git clone github.com/ozbekburak/dfir-radar
cd dfir-radar
go run .
```

The generated CSV file can be found in the reports directory along with a sample CSV file for reference. To generate a quadrant-style report, the user can run the following command in the terminal:


```python
cd quadrant
python3 main.py --csv ../reports/report.csv 
```

## Screenshot

![execute](https://github.com/ozbekburak/dfir-radar/blob/main/img/run.png?raw=true)
![report](https://github.com/ozbekburak/dfir-radar/blob/main/img/report.png?raw=true)


## Benefits

Our tool offers several benefits to DFIR professionals, including:

**Real-time insights:** Our tool provides real-time insights into the industry, allowing users to stay up-to-date with the latest trends and developments.

**Data-driven decision making:** Our reports are based on data, enabling users to make informed decisions based on the latest industry trends.

**Time-saving:** Our tool saves time by automating the data collection and categorization process, reducing the need for manual reading and analysis.

**Enhanced productivity:** With real-time insights and data-driven decision making, users can enhance their productivity by focusing on the most important tasks.

## Conclusion

Our automated DFIR reporting tool offers a new and innovative way for DFIR professionals to stay up-to-date with the latest industry trends and developments. By leveraging AI and real-time data, our tool provides users with the insights they need to make informed decisions and stay ahead of the curve. We believe our tool will become an essential part of any DFIR professional's toolkit and look forward to bringing it to the market.

## Future Work

In the future, we aim to create Gartner-like charts to enhance the reporting capabilities of our tool further. These charts will provide even more detailed insights into the DFIR industry by allowing users to compare and contrast different technologies, products, and services.

