#  SSH client

Connect to remote server and run pre-configured commands setup in ./config/config.yml

If config file is present, SSH client will use the config file present and ignore embedded config.

If **no** config file is present the binary will look for embedded config. 


## Usage

```

Configure ./config/config.yml for your remote hosts

./build.sh

```

## Config file format

```yaml
host:
- name: NAME
  connection: IP:PORT # IP address and port SSH server running on remote host
  username: USERNAME # user to login as on remote host
  password: PASSWORD  
  commands:
    - name: Name of command  # User-friendly alias
      string: command string # Command to run on remote host 
      userinput: true # set to true for prompt to append user input to the command string
      whitespace: false # Used to insert white before the end of command string

    - name: List files 
      string: ls -alh 
      userinput: false
      whitespace: false

    - name: Search Logs 
      string: cd /var/log; find . -name * -print0 | xargs -0 grep   
      userinput: true 
      whitespace: true 

```
- **Only password and keyboard-interactive authentication supported at this time. keyboard-interactive will automatically fill the password field provided in /config/config.yml when prompted.**

## How to embed config in binary

Either run ./build.sh or manually embed config. 

Install go-bindata

```
go get -u github.com/go-bindata/go-bindata/...
```

Embed config folder into binary

```
go-bindata -o config.go config
```

