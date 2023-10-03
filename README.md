

# Keycloak PGSQL PAM

Fork of https://github.com/kha7iq/kc-ssh-pam but for PostgreSQL

<p align="center">
  <a href="#install">Install</a> •
  <a href="#usage">Usage</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#contributing">Contributing</a> •
</p>

**kc-pgsql-pam** designed to streamline the process of user authentication and enable users to access PGSQL through OIDC. The program integrates with Keycloak to obtain a password grant token based on the user's login credentials, including their username and password. If two-factor authentication is enabled for the user, the program supports OTP code as well.

Once the password grant token is obtained, the program verifies it and passes the necessary parameters so that the user can be authenticated via OIDC and create default user role in PGSQL.

## Install

<details>
    <summary>DEB & RPM</summary>

```bash
# DEB
sudo dpkg -i kc-pgsql-pam_amd64.deb

# RPM
sudo rpm -i kc-pgsql-pam_amd64.rpm

```
</details>


<details>
    <summary>Manual</summary>

```bash
# Chose desired version
export KC_PGSQL_PAM_VERSION="0.1.1"
wget -q https://github.com/aborche/kc-pgsql-pam/releases/download/v${KC_PGSQL_PAM_VERSION}/kc-pgsql-pam_linux_amd64.tar.gz && \
tar -xf kc-pgsql-pam_linux_amd64.tar.gz && \
chmod +x kc-pgsql-pam && \
sudo mkdir -p /opt/kc-pgsql-pam && \
sudo mv kc-pgsql-pam config.toml /opt/kc-pgsql-pam
```
</details>


## Usage
```bash
❯ kc-pgsql-pam --help
Usage: kc-pgsql-pam USERNAME PASSWORD/[OTP]

Generates a password grant token from Keycloak for the given user.

Options:
  -h, --help              Show this help message and exit
  -v, --version           Show version information

Notes:
  For the program to function properly, it needs to locate a configuration file called 'config.toml'.
  The program will search for this file in the current directory, default install '/opt/kc-pgsql-pam', '/etc/kc-pgsql-pam/config.toml',
  and '$HOME/.config/config.toml', in that specific order.

  In addition to defaults, all configuration parameters can also be provided through environment variables.

  $KC_PGSQL_REALM $KC_PGSQL_ENDPOINT $KC_PGSQL_CLIENTID $KC_PGSQL_CLIENTSECRET $KC_PGSQL_CLIENTSCOPE
  
  To use the program, you must create a client in Keycloak and provide the following 
  information in the configuration file: realm, endpoint, client ID, client secret, and 
  client scope is optional.

Arguments:
  USERNAME                The username of the user is taken from $PAM_USER environment variable
  PASSWORD                The password of the user is taken from stdIn
  OTP                     (Optional) The OTP code if two-factor authentication is enabled i.e (password/otp)

  EXAMPLE                 (With otp): echo testpass/717912 | kc-pgsql-pam (Only Password): echo testpass | kc-pgsql-pam
```

## Configuration
  For the program to function properly, it needs to locate a configuration file called `config.toml`.
  
  The program will search for this file in the follwoing order..
  1. Present working directory
  2. Default install location `/opt/kc-pgsql-pam/config.toml`
  3. System `/etc/kc-pgsql-pam/config.toml`,
  4. `$HOME/.config/config.toml`

### Keycloak Client Creation
```bash
Step 1: Log in to the Keycloak Administration Console.

Step 2: Select the realm for which you want to create the client.

Step 3: Click on "Clients" from the left-hand menu, and then click on the "Create" button.

Step 4: In the "Client ID" field, enter "pgsql-pam".

Step 5: Set the "Client Protocol" to "openid-connect".

Step 6: In the "Redirect URIs" field, enter "urn:ietf:wg:oauth:2.0:oob".

Step 7: In the "Access Type" field, select "confidential".

Step 8: In the "Standard Flow Enabled" field, select "ON".

Step 9: In the "Direct Access Grants Enabled" field, select "ON".

Step 10: Click on the "Save" button to create the client.

To get the credentials of the client, follow these steps:

Step 1: Go to the "Clients" page in the Keycloak Administration Console.

Step 2: Select the "pgsql-pam" client from the list.

Step 3: Click on the "Credentials" tab.

Step 4: The client secret will be displayed under the "Client Secret" section.
```

### Config file template

`config.toml`
```toml
realm = "pgsql-demo"
endpoint = "https://keycloak.example.com"
clientid = "pgsql-pam"
clientsecret = "St0pUs1nGLDAPf0rPostgr3SQL"
clientscope = "openid"
groupsclaim = "groups"
allowedgroups = []
# if you need enable group check, add group to array
# allowedgroups = ["pgsql-user","pgsql-admins"]
alloweddomains = []
# if you need enable domains check, add domain name to array
#alloweddomains = ["example.com", "example1.com"]
```

### Basic local testing

Put your test config to *$HOME/.config/config.toml*

```bash
set +o history
export PAM_USER="user@domain.com"
echo UserPAssw0rd | /opt/kc-pgsql-pam/kc-pgsql-pam
```

### PostgreSQL setting up

* Create or edit `/etc/pam.d/postgresql` and add the following at the top of file
```bash
auth sufficient pam_exec.so expose_authtok      log=/var/log/kc-pgsql-pam.log     /opt/kc-pgsql-pam/kc-pgsql-pam
```
- User is not automatically created during login, so a local user must be present on the system before hand.

To automatically create a user install 
```bash
apt-get install libpam-script
```
Add the following in `/etc/pam.d/postgresql` underneath previous argument
```bash
account optional pam_script.so debug dir=/etc/kc-pgsql-pam/scripts
```

Then, the script itself. In the file `/etc/kc-pgsql-pam/scripts/pam_script_acct`
```bash
#!/bin/sh

CREATEUSER="
DO
\$do\$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_catalog.pg_roles 
      WHERE  rolname = '${PAM_USER}') THEN

      CREATE ROLE \"${PAM_USER}\" WITH LOGIN;
   END IF;
END
\$do\$;"

echo ${CREATEUSER} | psql > /dev/null
```
In PAM modules, username is given in "$PAM_USER" variable.

Make this script executable
```bash
sudo chmod +x /etc/kc-pgsql-pam/scripts/pam_script_auth 
```

Set files permission
```bash
sudo touch /var/log/kc-pgsql-pam.log
sudo chown postgres:postgres /var/log/kc-pgsql-pam.log
sudo chown -R postgres:postgres /etc/kc-pgsql-pam
sudo chown -R root:root /etc/kc-pgsql-pam/scripts /etc/pam.d/postgresql
```

Enable pam check for postgresql. Edit /etc/postgresql/\<version>/\<clustername>/pg_hba.conf and add following line to end of file
```bash
host    all             all             all         pam pamservice=postgresql
```

Reload postgresql service
```bash
sudo systemctl reload postgresql
```

### Check connection and authentication

Open console and run psql with wrong password for user
```bash
$ psql -d postgres -h xxx.xxx.xxx.xxx -U "user@domain.com" -W
Password: 
psql: error: FATAL:  PAM authentication failed for user "user@domain.com"
FATAL:  PAM authentication failed for user "user@domain.com"
```

Run psql with correct password for user
```bash
$ psql -d postgres -h xxx.xxx.xxx.xxx -U "user@domain.com" -W
Password: 
psql (13.11 (Debian 13.11-0+deb11u1), server 12.15 (Ubuntu 12.15-1.pgdg20.04+1))
SSL connection (protocol: TLSv1.3, cipher: TLS_AES_256_GCM_SHA384, bits: 256, compression: off)
Type "help" for help.

postgres=> \conninfo
You are connected to database "postgres" as user "user@domain.com" on host "xxx.xxx.xxx.xxx" at port "5432".
SSL connection (protocol: TLSv1.3, cipher: TLS_AES_256_GCM_SHA384, bits: 256, compression: off)
```

Check postgresql and kc-pgsql-pam logs for errros

```bash
$ tail /var/log/kc-pgsql-pam.log
*** Thu Sep 28 11:44:55 2023
2023/09/28 11:44:55 OIDC Auth: 'user@domain.com' Failed to retrieve token: HTTP request failed with status code 401

*** Thu Sep 28 11:44:59 2023
2023/09/28 11:44:59 OIDC Auth: 'user@domain.com' Token acquired and verified Successfully.
```

:diamond_shape_with_a_dot_inside: Detailed article with screenshots is also [available here](https://lmno.pk/post/kc-sso-pam/)

## Contributing

Contributions, issues and feature requests are welcome!<br/>Feel free to check
[issues page](https://github.com/aborche/kc-pgsql-pam/issues). You can also take a look
at the [contributing guide](https://github.com/aborche/kc-pgsql-pam/blob/master/CONTRIBUTING.md).
