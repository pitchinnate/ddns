# Self Hosted DDNS

I built this so that you can run your own dynamic dns service on your home computer. I have a prebuilt windows
binary, but you are welcome to download the source and build your own. Currently only works with Route53 on AWS.

## Description

This is basically a self-hosted Dynamic DNS application. First it hits an api endpoint (`https://api.ipify.org?format=json`) 
once per day, gets your current home IP address. It then updates an A record in Route53 to reflect that ip address. While
some routers offer a built-in DDNS option, you are normally limited to using a DDNS service that could have a cost and also
I personally have found many of them unreliable at best. This can be ran on any computer that sits inside your home/office network.

## Setup

- First you will need a domain name hosted on AWS's Route53.
- You will also need get AWS access keys for IAM account that has access to read/write records to Route53 (I used `AmazonRoute53FullAccess`).
- There is a `.env.example` file with the required environment variables. You can copy this file to `.env` and fill in
  the required environment variables.
  - AWS_ACCESS_KEY_ID = AWS access key
  - AWS_SECRET_ACCESS_KEY = AWS secret access key
  - AWS_REGION = AWS region your Route53 domain is setup
  - HOSTED_ZONE_ID = the route53 Hosted Zone ID
  - RECORD_NAME = the A record you want to store the IP address in, example: `homevpn.test.com`
- You can now simply run the `ddns.exe` (or whatever name you build for) application or you could setup your Operating System
    to run it at startup.
  - for Windows you can use Task Scheduler to have it run the app at login or boot.
  - for Linux there are multiple options to run at startup (systemd, supervisor, etc...)

## Building
For windows you will want use the following build flags so that a terminal window doesn't open when you run 
the application.

**Example Windows build**
```bash
go build -ldflags -H=windowsgui -o ddns.exe
```

**Example Linux build**
```bash
go build -o ddns
```

## Why do you need Dynamic DNS?

Dynamic DNS (DDNS) is commonly used when you want to provide a consistent domain name for a device or service that doesn’t have
a static IP address, such as those on a home or small business network with a dynamic IP. Here are some typical use cases:

- **Remote Access to Home or Office Networks**
    - Personal Servers: Many people set up home servers for personal websites, game servers, media streaming, or file sharing (e.g., Plex or Nextcloud). DDNS allows them to access these services remotely without worrying about IP changes.
    - Remote Desktop or VPN Access: For securely accessing computers or networks remotely, DDNS provides a stable address for VPN connections or remote desktop sessions to a home or small office network.
- **Internet of Things (IoT) Devices**
    - Smart Home Automation: Devices like cameras, thermostats, or smart lighting systems often need remote access. With DDNS, users can easily control and monitor IoT devices without needing to know the current IP address.
    - Security Cameras: Accessing security cameras remotely is a common use case for DDNS, as users can monitor their property through a fixed domain name, even as IP addresses change.
3. **Gaming Servers**
    - Online Gaming: If you’re hosting a private game server (e.g., Minecraft, Counter-Strike), DDNS lets players connect using a consistent domain name instead of requiring you to update everyone when the IP address changes.
    - LAN Parties or Group Gaming: Dynamic DNS is useful for organizing games with friends or setting up game sessions that may need regular connectivity with known addresses.
4. **Hosting Applications on Dynamic IP Networks**
    - Web Development and Testing: Developers often run web applications or APIs on local machines and want external collaborators to access them. With DDNS, they can provide access without setting up static IPs or exposing the local server IP.
    - Small Business Applications: Small businesses often host custom applications like inventory management systems or CRM tools locally. DDNS enables consistent remote access, crucial for users who work from home or on the road.
