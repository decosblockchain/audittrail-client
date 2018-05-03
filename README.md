# Decos Blockchain Audittrail Client
This client will sign blockchain transactions locally, containing the audittrail data, and transmit them to a server over HTTPS. The server will broadcast the transaction to the blockchain nodes for being added to a staging blockchain. The staging blockchain is periodically timestamped onto a public blockchain such as Bitcoin.

## Installation
The installation involves a few easy steps. 

### Download

First, download the appropriate release file for your platform from [our releases page](https://github.com/decosblockchain/audittrail-client/releases). 

### Unpack

NOTE: On Windows platforms make sure the download is unblocked before you unpack it.

Unpack the files into a folder of your preference. The log files are written into a subdirectory of where you place the application data. They are limited to 10MB (rolling logs). So you don't need excessive space on the device you are installing the software.

### Edit the configuration file

The file ```config.json``` contains two configuration parameters. Under normal circumstances, these don't need to be adjusted. You can however decide to alter them if your situation requires.

| Setting   | Description |
|-----------|-------------|
| ```ServerUrl``` | The HTTP endpoint where the transactions will be sent to |
| ```ListenPort``` | The port where the local API will be listening to. Default 8585 |

### Run for the first time

When you run the executable for the first time, it will generate a private key for signing the blockchain transactions. This file is placed in data/privatekey.hex. Make sure this file is backed up somewhere. Run the executable in console mode to make it print out your public key (address):

(For Windows:)

```
audittrail-client.exe console
```

(For Linux:)

```
./audittrail-client console
```

The output will be something like this:

```
INFO: 2018/05/03 11:32:50 client.go:69: Running in console mode
INFO: 2018/05/03 11:32:50 client.go:77: My address is 0x287a95CE0E1d4F7e83dFdfBf332e389F7f5F4b49
```

Before you can start sending blockchain transactions, the address needs to be funded by Decos; so convey this address to your contact and he will make sure the address is funded properly.

### Install as a service

Next, you can install the service by issuing the ```install``` command on the executable, like so:

(For Windows, make sure you do this from an Administrator command prompt:)

```
audittrail-client.exe install
```

(For Linux:)

```
sudo ./audittrail-client install
```

The output on the console will be:

```
2018/05/03 11:52:21 Service control action [install] executed succesfully
```

### Starting and stopping the service

You can start or stop the service by using the regular service manager (i.e. for Windows, the Services control panel). Or you can do it by issuing the ```start``` or ```stop``` command on the executable:

(For Windows, make sure you do this from an Administrator command prompt:)

```
audittrail-client.exe [start|stop]
```

(For Linux:)

```
sudo ./audittrail-client [start|stop]
```

The output on the console will be:

```
2018/05/03 11:52:34 Service control action [start|stop] executed succesfully
```

### Checking if the service is working

If you open ```http://localhost:[ListenPort]/``` in your browser (so for the default configuration this would be [http://localhost:8585/](http://localhost:8585), you will see a friendly message "It works!". This is a signal that the server is listening and ready to accept requests.

### Uninstall

If you wish to uninstall the service from your system, you can issue the ```uninstall``` command, like so:

(For Windows, make sure you do this from an Administrator command prompt:)

```
audittrail-client.exe uninstall
```

(For Linux:)

```
sudo ./audittrail-client uninstall
```

The output on the console will be:

```
2018/05/03 11:52:34 Service control action [uninstall] executed succesfully
```