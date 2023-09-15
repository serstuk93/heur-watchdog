# WatchDog - Product Monitor 🐶




A simple Go application designed to monitor product listings from certain URLs and display the results in a window using the Fyne GUI framework. Currently in Slovak language only. Its possible to create apk file for Android devices.

![Screenshot](https://raw.github.com/serstuk93/heur-watchdog/master/screenshot1.PNG)

## Table of Contents
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [How It Works](#how-it-works)
- [Future Enhancements](#future-enhancements)
- [Acknowledgements](#acknowledgements)
- [Contributing](#contributing)
- [License](#license)

## Features

1. **Product Monitoring**: Watches specific URLs for product availability.
2. **GUI Display**: Uses the Fyne framework to display results in a clean user interface.
3. **Clickable Links**: Any product links in the results are clickable.
4. **Refresh Capability**: Can manually refresh results or set to auto-refresh every 1 minute.
5. **Notifications**: Notifies the user if a new product appears upon refreshing.


## Requirements

1. Go (at least version 1.19).
2. Fyne for GUI.
3. Logrus for enhanced logging capabilities.

## Installation

1. Ensure Go is installed.
2. Clone the repository.

```bash
git clone https://github.com/serstuk93/heur-watchdog.git
```
3. Navigate to the repository directory and run:

```bash
go get fyne.io/fyne/v2
go get github.com/sirupsen/logrus
```

4. Run the program
```bash
go run .
```

5. Compile Apk file for Android device

    - install mingw-w64-x86_64-toolchain
    - set path inside mingw

    ```bash 
    echo "export PATH=\$PATH:/c/Program\ Files/Go/bin:~/Go/bin" >> ~/.bashrc
    ```
    - verify env variables by typing commands:

    ```bash
    go version
    gcc--version
    fyne
    echo $ANDROID_NDK_HOME
    echo $ANDROID_HOME
    ```
    - compile apk

    ```bash
    fyne package -os android -appID com.serstuk93.watchdog -icon icon.png
    ```

    - copy and install apk to your android device or run it virtually via Android Studio 

## How It Works

Upon starting, the app fetches products from predefined URLs. The results are then displayed in a Fyne window where each product has its own clickable link. The user has the option to manually refresh the results using the "Refresh" button provided. Additionally, the application has been set to auto-refresh every hour. If there happen to be new products that appear after a refresh, the app will send a desktop notification to alert the user.

## Future Enhancements

1. Enhancement in the areas of detailed error handling and logging.
2. Provision for users to add/edit/delete URLs that they wish to monitor.

## Acknowledgements

- A big thank you to [Fyne](https://fyne.io/) for their exceptional GUI framework.
- Also, appreciation goes to [Logrus](https://github.com/sirupsen/logrus) for their enhanced logging capabilities.

## Contributing

If you've got suggestions, improvements, or any other feedback, I encourage you to submit issues or pull requests. Your contribution is highly valued.


## License 

[![License: CC BY 4.0](https://img.shields.io/badge/License-CC%20BY%204.0-lightgrey.svg)](https://creativecommons.org/licenses/by/4.0/)
This work is licensed under a Creative Commons Attribution 4.0 International License.
