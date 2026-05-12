# wgsl — Leukemia Pattern Discovery CLI

A command-line tool for medical professionals to run machine learning-based leukemia subtype detection from gene expression data. Results are delivered securely to the registered doctor's email.

> **Disclaimer:** This tool is intended to aid medical professionals in making informed decisions. It is not a substitute for clinical judgment or professional diagnosis.

---

## Overview

`wgsl` orchestrates a full ML pipeline — from data preprocessing and PCA-based dimensionality reduction to SVM classification and visualization — and emails the results directly to the registered practitioner.

The pipeline uses the [Leukemia Gene Expression (CuMiDa) dataset](https://www.kaggle.com/datasets/brunogrisci/leukemia-gene-expression-cumida) and produces three visualisations:

- **PCA Cluster Analysis** — Scatter plot of leukemia subtypes in principal component space
- **Confusion Matrix** — Classification performance heatmap
- **Scree Plot** — Explained variance per principal component

---

## Prerequisites

- **Go** 1.25+
- **Python 3** with `pip3`
- A configured `.env` file (see [Environment Setup](#environment-setup))
- Kaggle API access (to download the dataset)

---

## Installation

### 1. Clone the repository

```bash
git clone <your-repo-url>
cd wgsl
```

### 2. Grant execution permissions

```bash
bash getPermission.sh
```

### 3. Download the dataset and install Python dependencies

```bash
bash GetDataset.sh
```

This will:
- Download and unzip the leukemia gene expression dataset from Kaggle
- Rename it to `Leukemia.csv`
- Install all required Python packages

### 4. Build the CLI

```bash
go build -o wgsl .
```

---

## Environment Setup

Create a `.env` file in the project root with the following variables:

```env
FROM_EMAIL=your-sender@gmail.com
FROM_EMAIL_PASSWORD=your-app-password
FROM_EMAIL_SMTP=smtp.gmail.com
SMTP_ADDR=smtp.gmail.com:587
```

> **Note:** For Gmail, use an [App Password](https://support.google.com/accounts/answer/185833) rather than your account password.

---

## Usage

All commands follow the pattern:

```
wgsl <command>
```

### Commands

| Command  | Description |
|----------|-------------|
| `help`   | Display all available commands and their descriptions |
| `init`   | Register a doctor's email with OTP verification and initialise the working directory |
| `train`  | Preprocess the dataset, train the SVM model, and generate visualisations |
| `result` | Email all generated visualisations to the registered doctor's address |

### Typical Workflow

```bash
# 1. Initialise and register the doctor's email
./wgsl init

# 2. Train the model and generate plots
./wgsl train

# 3. Send results to the registered email
./wgsl result
```

---

## How It Works

### Initialisation (`init`)

- Prompts for the doctor's email (input is hidden for privacy)
- Validates the domain (Gmail, Outlook, iCloud, Yahoo supported)
- Sends a 6-digit OTP to the provided address
- OTP expires after **5 minutes**
- On success, stores the email in a local `.wgsl` config file

### Training (`train`)

Executes `main.py`, which:

1. Loads `Leukemia.csv`
2. Encodes leukemia type labels and fills missing values
3. Resamples the dataset with `RandomOverSampler` to handle class imbalance
4. Applies PCA (retaining 95% of variance) and normalises features
5. Trains an SVM classifier (`SVC`) on the processed data
6. Prints accuracy to the console
7. Saves three visualisation plots to the `Images/` directory

### Results (`result`)

- Reads the registered email from `.wgsl`
- Attaches all `.png` files from the `Images/` folder
- Sends a multipart email with the visualisations as attachments
- Deletes the `.wgsl` config file after sending

---

## Project Structure

```
wgsl/
├── main.go             # CLI entry point and command routing
├── cmds/
│   ├── init.go         # Email registration and OTP verification
│   ├── result.go       # Email delivery of results
│   └── help.go         # Help text
├── main.py             # ML pipeline (preprocessing, training, visualisation)
├── Leukemia.csv        # Gene expression dataset (downloaded via GetDataset.sh)
├── Images/             # Output visualisations (auto-created)
├── GetDataset.sh       # Downloads dataset and installs Python deps
├── dataset.sh          # Alternative dataset download (white blood cells)
├── getPermission.sh    # Sets execute permissions on shell scripts
├── requirments.txt     # Python dependencies
├── go.mod / go.sum     # Go module files
└── .env                # SMTP credentials (not committed)
```

---

## Python Dependencies

Listed in `requirments.txt`:

```
pandas
numpy
scikit-learn
matplotlib
imblearn
```

Install manually with:

```bash
pip3 install -r requirments.txt
```

---

## Go Dependencies

- [`github.com/joho/godotenv`](https://github.com/joho/godotenv) — `.env` file loading
- [`golang.org/x/term`](https://pkg.go.dev/golang.org/x/term) — Hidden terminal input for email and OTP
- [`golang.org/x/crypto`](https://pkg.go.dev/golang.org/x/crypto) — Cryptographic utilities

---

## License

See [LICENSE](./LICENSE) for details.
