<h1 style="text-align: center"> CoWitness </h1>

<h3 style="text-align: center">

![Images\CoWitness.png](https://github.com/officialjm/cowitness/blob/main/Images/CoWitness.png)

[![License: MIT](https://img.shields.io/badge/License-MIT-darkred.svg)](https://github.com/stolenusername/cowitness/blob/main/LICENSE)
[![made-with-python](https://img.shields.io/badge/Made%20with-GoLang-blue.svg)](https://go.dev/)

</h3>


**CoWitness** is a tool designed to function as an HTTP, HTTPS, and DNS server. It allows you to serve web pages, log HTTP requests, and customize DNS responses. This README provides an overview of the features and instructions on how to compile and use the CoWitness tool.

### Features 

- **HTTP Server**: CoWitness includes an HTTP server that listens on **port 80**. It can serve static files from the current working directory. Each HTTP request is logged, including the client's IP address, requested resource, and user agent.

- **HTTPS Server**: In addition to the HTTP server, CoWitness also provides an HTTPS server that listens on port 443. Similar to the HTTP server, it serves static files and logs each request.

- **DNS Server**: CoWitness functions as a DNS server, listening on **port 53**. It allows you to customize DNS responses, including NS and A records. DNS requests are logged, including the client's IP address and the requested domain.

- **Logging**: CoWitness creates log files for HTTP and DNS requests. The logs are appended to existing files, allowing you to track and analyze the server activity.

- **Quiet Mode**: CoWitness can run in quiet mode by passing the `-q` command-line argument. In this mode, the ASCII art banner will not be displayed.

## Prerequisites 📝

Before using CoWitness, ensure that you have the following requirements:

- Go programming language installed on your system, get it [HERE](https://go.dev/) go 1.21 required if building from source.
- Internet access to download Go dependencies.
- A remote server with a public IP address.

## Build 👨🏼‍🔧

Follow these steps to install and compile CoWitness:

1. Clone the CoWitness repository from GitHub:

```bash
git clone https://github.com/your-username/CoWitness.git
```

2. Change to the CoWitness directory:

```bash
cd CoWitness
```

3. Build the CoWitness executable:

```bash
make -f Makefile build
```
This command compiles the CoWitness source code and creates an executable file for amd64 and arm64.


## Installation 👨🏻‍💻

**To use CoWitness on a remote server, follow these steps**:

1. Choose a domain name for your testing environment.
2. Set up a remote server and obtain a public IP address for it.
3. Register your name servers to point to the public IP address. (_Refer to your domain registrar's documentation for instructions on how to set up name servers_.)
4. Create glue records to associate the IP address with your remote server. This step may vary dpending on your DNS provider. Consult their documentatin for guidance on creating glue records.
5. Ensure that ports 80 and 53 are available on the remote server. Configure any firewall or network settings to allow incoming connections on these ports.
6. Transfer the CoWitness executable to the remote server using a secure method such as SCP or SFTP.

Connect to the remote server via SSH:

```bash
ssh username@your-server-ip
```
Easy way to install cowitness is to get the automated package build from [here](https://github.com/stolenusername/cowitness/releases/latest)
which you can then curl or wget to your system or scp the file over from your system.
```bash
cowitness
```

## Customization ⚒️

You can customize CoWitness to fit your specific needs. Here are some possible modifications:

- **Change the default ports**: See the usage help for the cli flags to change the default ports

- **Modify the log file paths**: You can change the paths for the HTTP and DNS log files (http.log and dns.log) by updating the `os.OpenFile` calls in the source code.

<br></br>
### Community & Contributions

We welcome contributions and feedback from the community. If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request on our GitHub repository!

### License

CoWitness is released under the [MIT License](LICENSE).

