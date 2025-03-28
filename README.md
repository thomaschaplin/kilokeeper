<img src="web/assets/logo.png" alt="logo" width="500" height="250" />

# KiloKeeper

A simple weight tracking application with a web-based interface that allows users to log their weight over time and visualize it using a chart.

## Features
- Log weight entries with date validation.
- View a historical graph of weight changes over time.
- Weight over time is displayed in a chart.
- Data is persisted in a JSON file (`weights.json`).
- Calculate and display BMI and BMI status.
- Display age based on date of birth.
- Display weight in kilograms and stones, pounds, and ounces.
- Display goal weight.

## Demo

See below screenshot for example

<img src="web/assets/demo.png" alt="demo" width="435" height="435" />

## Technologies Used
- **Backend:** Go (net/http)
- **Frontend:** HTML, JavaScript, Chart.js
- **Data Storage:** JSON files

## Setup and Installation

### Prerequisites
- Go installed (1.16+ recommended)
- A working web browser

### Environment Variables
- `DOB`: Date of birth in `DD/MM/YYYY` format.
- `HEIGHT`: Height in meters. e.g. `1.75`
- `GOAL`: Goal weight in kilograms. e.g. `70.2`

### Steps to Run
1. **Clone the repository:**
   ```sh
   git clone https://github.com/thomaschaplin/kilokeeper.git
   cd kilokeeper
   ```
2. **Create required files and directories:**
   ```sh
   mkdir data
   touch data/weights.json
   ```
3. **Initialize the `weights.json` file:**
   ```json
   []
   ```
4. **Set environment variables:**
   ```sh
   export DOB="01/01/1990"
   export HEIGHT="1.75"
   export GOAL="70.2"
   ```
5. **Run the application:**
   ```sh
   go run main.go
   ```
6. **Access the web app:**
   Open `http://localhost:8080` in your browser.

## API Endpoints

### Get All Weights
**GET** `/weights`

### Add a New Weight Entry
**POST** `/weights/add`

#### Request Body (JSON):
```json
{
  "date": "DD/MM/YYYY",
  "kilograms": 75.5
}
```

## Future Improvements
- [X] Add BMI calculation.
- [X] Add age to chart based on global date of birth
- [ ] Store data in a database instead of JSON files.
- [ ] Add multiple users.
- [ ] Add user authentication.

## License
MIT License

---
Made with ❤️ by Thomas Chaplin

