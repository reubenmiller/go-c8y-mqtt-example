# go-c8y-mqtt-example

A simple example of an MQTT client for the Cumulocity Platform. It creates a new device and pushes measurements periodically into the platform.

## Requirements

* go v1.13 or newer

## Getting started

1. Clone the repository

    ```sh
    git clone https://github.com/reubenmiller/go-c8y-mqtt-example.git
    cd go-c8y-mqtt-example
    ```

2. Activate go modules

    ```sh
    # Powershell
    $env:GO111MODULE = "on"

    # Windows cmd.exe
    SET GO111MODULE=on

    # bash
    export GO111MODULE=on
    ```

3. Configure the Cumulocity target settings

    ```sh
    # Powershell
    $env:MQTT_PROTOCOL = "wss"                  # Optional
    $env:C8Y_DEVICE_NAME = "customDeviceName"   # Optional
    $env:C8Y_HOST = "https://cumulocity.com"
    $env:C8Y_TENANT = "mytenant"
    $env:C8Y_USER = "myuser@mail.com"
    $env:C8Y_PASSWORD = "myPassw0rd"

    # Windows cmd.exe
    SET MQTT_PROTOCOL=wss
    $env:C8Y_DEVICE_NAME = "customDeviceName"
    SET C8Y_HOST=https://cumulocity.com
    SET C8Y_TENANT=mytenant
    SET C8Y_USER=myuser@mail.com
    SET C8Y_PASSWORD=myPassw0rd

    # bash
    export MQTT_PROTOCOL=wss                        # Optional
    export C8Y_DEVICE_NAME = "customDeviceName"     # Optional
    export C8Y_HOST=https://cumulocity.com
    export C8Y_TENANT=mytenant
    export C8Y_USER=myuser@mail.com
    export C8Y_PASSWORD=myPassw0rd
    ```

4. Run the client

    ```sh
    go run main.go
    ```


## Building

```sh
go build main.go
```
