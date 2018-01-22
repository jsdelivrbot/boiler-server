# BoilerGo
[![Build Status](https://travis-ci.com/AzureRelease/boiler-server.svg?token=za632F62BvSeRXAtUssN&branch=master)](https://travis-ci.com/AzureRelease/boiler-server)

## Deploy Manual *(With Ubuntu 16.x)*
### 1. Docker Env Install   
##### *(Install Command-Line CAN NOT Use Shell Script)*
```sbtshell
$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
```

```sbtshell
$ sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
```

```sbtshell
$ sudo apt-get update
```

```sbtshell
$ apt-cache policy docker-ce
```

```sbtshell
$ sudo apt-get install -y docker-ce
```

### 2. Unpacking Docker Image(s)
```sbtshell
$ docker-boiler-run
```

### 3. Startup The Boiler-Server And Boiler-Data
```sbtshell
$ docker-boiler-startup
```
