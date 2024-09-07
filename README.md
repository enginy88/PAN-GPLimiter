# PAN-GPLimiter

A Go program for limiting concurrent remote user logins in a single GP Gateway on a PANOS Firewall.


## What's the motivation?

This one is maybe the most ever wanted feature request of Global Protect for decades! (FR4603-Concurrent Session Limiting) After tons of FR votes, endless requests from customers, lots of reddit messages asks for workarounds, people who are in charge don't have in the same opinion with the technical guys who are on the field as they haven't green lighted for developers to implement this super easy feature for years.

Finally, I ran out of hope and couldn't remain more indifferent to it. So this forces me to create my own home-brewed solution and I give myself the go-ahead.


## A Brief History:

Once I started to implement this program, there was only a PowerShell script [^1] dating from 2018. I haven't tried it by myself but many ones couldn't make it run for some reason. (Or it really doesn't run at all!) Assuming it works, it's also OS (Windows) dependent, inefficient, couldn't handle edge-cases, lacks some features, etc... But besides that, it did its job as it inspired me and led the way to me!

After I created this program, I've found that someone else also created a Python script [^2] in 2020. I was surprised when faced with that since I didn't realize there was such an attempt at all. Honestly if I had known about it, I may never have started at first. You can also check this work since it provides some different features than this one.


## Why Golang?

Because it (cross) compiles into machine code! You can directly run ready-to-go binaries on Windows, Linux, and macOS. No installation, no libraries, no dependencies, no prerequisites... Unlike Bash/PowerShell/Python it is not interpreted at runtime, which drastically reduces runtime overhead compared to scripting languages. The decision to use a compiled language makes it run lightning fast with lower memory usage. Also, due to the statically typed nature of the Go language, it is more error-proof against possible bugs/typos.


## Usage:

Pre-compiled binaries can be downloaded directly from the latest release (Link here: [Latest Release](https://github.com/enginy88/PAN-GPLimiter/releases/latest)). These binaries can readily be used on the systems for which they were compiled. Neither re-compiling any source code nor installing GO is not needed. In case there is no pre-compiled binary presented for your system, you can refer to the [Compilation](#compilation) section.

This program only requires the `appsett.env` file to determine which settings it will run with. This file must be present and accessible.

By default, the program searches for the `appsett.env` file in the working directory. The working directory can be altered by passing `-dir [PATH]` argument. (`PATH` value for `-dir` can be ablsolute or relative.) These options can be explored by passing the `-help` argument to the program:

```shell
# Usage of ./PAN-GPLimiter:
  -dir string
        Path of directory which contains 'appsett.env' file. (Optional)
```

There are multiple settings controlled in environment variable format in the `appsett.env` file. These settings are explained in the [Settings](#settings) section.

How to schedule this program is not within the scope of this program or this documentation. Nevertheless, some hints are shared in the [Hints](#hints) section.


## Settings:

All setting options are provided with the sample `appsett.env` file, along with short descriptions and default values for each. Note that all lines are commented-out in the sample. If any of the options are to be used simply clear the comment token (`#`) and set the preferred setting. First three setting options are mandatory and must be filled prior to use the program.

All options that are provided by the `appsett.env` file can be overriden by environmental variables when needed. In fact, an empty `appsett.env` file may be used if all mandatory and necessary optional options are provided by environmental variables.

```shell
# Basic & Optional Settings:
PANGPLIMITER_FIREWALL_HOST={Enter IP or FQDN of PAN NGFW s MGMT Capable Interface}
PANGPLIMITER_API_KEY={Insert API Key for an Admin Account}
PANGPLIMITER_GP_GATEWAY={Enter Name of PAN Global Protect Gateway}
PANGPLIMITER_VSYS_NO={Enter VSYS Number of Target Device, Use 1 If VSYS Disabled, Default: 1}
PANGPLIMITER_MAX_LOGIN={Enter Max Login Count to Allow, Default: 1}
PANGPLIMITER_EXCLUDED_USERS={Enter List of Primary Usernames to Exclude from Restriction, Space Separated, Default: NONE}

# Advanced & Expert Options:
PANGPLIMITER_LOG_SILENCE={Enter Either TRUE or FALSE to Enable Verbose Log Suppression, Default: FALSE}
PANGPLIMITER_SKIP_VERIFY={Enter Either TRUE or FALSE to Skip TLS Certificate Verification, Default: FALSE}
PANGPLIMITER_KICK_OLDEST={Enter Either TRUE or FALSE to Kick Oldest Sessions Instead of Newests, Default: FALSE}
PANGPLIMITER_LIST_ALL={Enter Either TRUE or FALSE to Show List of All Sessions, Default: FALSE}
PANGPLIMITER_MULTI_THREAD={Enter Either TRUE or FALSE to Enable Parallel Processing, Default: FALSE}
PANGPLIMITER_FAIL_ONERROR={Enter Either TRUE or FALSE to Enable Fail on First Error, Default: FALSE}
PANGPLIMITER_DRY_RUN={Enter Either TRUE or FALSE to Enable Dry Run without Enforcement, Default: FALSE}
PANGPLIMITER_FILE_TEST={Enter Either TRUE or FALSE to Run Offline Test from "test.xml" File, Default: FALSE}
```

<details>

<summary>Long explanation of each settings option: (Expand to view)</summary>

### Explanation of Settings:

**PANGPLIMITER_FIREWALL_HOST** 

TYPE: ```String``` DEFAULT VALUE: ```NONE (Mandatory)``` 

This setting controls which IP Address or FQDN the program will use to connect to the Firewall. This setting option must be provided to start the program.

**PANGPLIMITER_API_KEY**

TYPE: ```String``` DEFAULT VALUE: ```NONE (Mandatory)``` 

This setting controls which API Key the program will use to authenticate with the Firewall. This setting option must be provided to start the program.

**PANGPLIMITER_GP_GATEWAY**

TYPE: ```String``` DEFAULT VALUE: ```NONE (Mandatory)``` 

This setting controls which GP Gateway the program will use to operate on the Firewall. This setting option must be provided to start the program.

**PANGPLIMITER_VSYS_NO**

TYPE: ```Integer``` DEFAULT VALUE: ```1``` 

This setting is to set working VSYS of given GP Gateway. VSYS number and GP Gateway must be matched in order for the program to work. If "Multi Virtual System Capability" is disabled or if the Firewall does not supports this, then VSYS number of 1 should be used which is the default value already.

**PANGPLIMITER_MAX_LOGIN**

TYPE: ```Integer``` DEFAULT VALUE: ```1``` 

When set to more than 1, the behavior of the program changes to allow multiple concurrent sessions for a user limited by the given value. With the default value, only single login will be allowed.

**PANGPLIMITER_EXCLUDED_USERS**

TYPE: ```String``` DEFAULT VALUE: ```NONE``` 

For the given usernames, the program skips processing for specified users, which means they will be exempted and can be use multiple concurrent logins. Usernames must be in the correct format as they seen on the Firewall like DOMAIN\USER. More than one space separated username can be specified.

**PANUSOMXML2EDL_LOG_SILENCE**

TYPE: ```Boolean``` DEFAULT VALUE: ```FALSE``` 

This program has 4 levels of log types: Always, error, warning and log. When this option is set, it will suppress the warning and info level logs. Note that the level of always cannot be silenced. Also, the level of error is shown all the time as it means there is an unrecoverable failure during the operation and the program will be terminated without completing the jobs.

**PANUSOMXML2EDL_SKIP_VERIFY**

TYPE: ```Boolean``` DEFAULT VALUE: ```FALSE``` 

Normally, Firewall should serve the XML API via the HTTPS protocol with a trusted TLS certificate. The default behavior is not to continue when the program encounters an untrusted certificate as it may indicate possible MITM attack, which tries to alter your enforcement. So, it is not advised to change this from the default value. However, it is implemented for possible use with in combination with a Firewall which may be served with an untrusted TLS certificate within your knowledge so that the program can proceed with an insecure connection.

**PANGPLIMITER_KICK_OLDEST**

TYPE: ```Boolean``` DEFAULT VALUE: ```FALSE``` 

When set, this option will change the default behavior of the program to kick the oldest sessions instead of newest sessions. Can be set as desired.

**PANGPLIMITER_LIST_ALL**

TYPE: ```Boolean``` DEFAULT VALUE: ```FALSE``` 

When set, this option will print all active session records to the output. It may be useful to troubleshoot or track active sessions whenever this programs runs.

**PANGPLIMITER_MULTI_THREAD**

TYPE: ```Boolean``` DEFAULT VALUE: ```FALSE``` 

Normally, this program runs in a single thread and operate all tasks sequentially including sending API requests. With this experimental feature, a worker pool will be used to send all kick API requests in parallel. Please note that lower and Firewalls cannot be able to handle requests in parallel. So by choking the management plane of the Firewall, some requests may fail. 

**PANGPLIMITER_FAIL_ONERROR**

TYPE: ```Boolean``` DEFAULT VALUE: ```FALSE``` 

This program sends API requests as many times as the number of people to be kicked. Some requests cannot be fulfilled by the Firewall due to many reasons. This is normal and eventually failed requests will be retried on the next iteration of run. So the program logs the kick errors and continues to work when encountered. For troubleshooting reasons, terminating the program on any error is possible with this.

**PANGPLIMITER_DRY_RUN**

TYPE: ```Boolean``` DEFAULT VALUE: ```FALSE``` 

When set, this program simulates the whole operation without sending kick API requests to the Firewall. It is useful to find out how many or which users are using concurrent sessions. This option can be used to exercise and verify which sessions will be terminated in normal operation mode. 

**PANUSOMXML2EDL_FILE_TEST**

TYPE: ```Boolean``` DEFAULT VALUE: ```FALSE``` 

Normally, this program fetches live active sessions from Firewall via API. When this option set, the program looks for the ```test.xml``` file under working the directory and uses it as input. It may be useful for operate on a preliminary taken information especially the Firewall is air-gapped.

</details>

## Hints:

When using this program under Unix-like OSes like Linux or macOS, Cron can be used to schedule periodic execution of the program. Here is an example of Crontab file:

```shell
# /etc/crontab
# To run the program at every 3 minutes:
*/3 * * * * /path_to_binary/PAN-GPLimiter -dir /path_to_folder_of_appsett.env/ &>> /path_to_any_logging_file.txt
```

Under Windows OSes, the Task Scheduler tool can be used for the same purpose.


## Compilation:

If none of the pre-compiled binaries covers your environment, you can choose to compile from source by yourself. Here are the instructions for that:

```shell
git clone https://github.com/enginy88/PAN-GPLimiter.git
cd PAN-GPLimiter
go get -u ./...
go mod tidy
make local # To compile for your own environment.
make # To compile for pre-selected environments.
```

**NOTE:** To compile from the source code, GO must be installed in the environment. However, it is not necessary for to run compiled binaries! Please consult the GO website for installation instructions: [Installing Go](https://go.dev/doc/install)


## TODO List:

Even I didn't think that those are show blockers, here is my todo list to make this program top-notch:

- [x] Program argument for path to config file: to select config file different then working directory.
- [x] Remove oldest sessions: to provide an option instead of kicking out the newest ones.
- [ ] Separate Mgmt IP support: to provide an alternative way to communicate with HA peers.
- [ ] Daemon mode: to get rid of Cron/Task Scheduler.
- [ ] Syslog and EventLog support: to write logs to OS when daemonized.


## Refs:
[^1]: https://live.paloaltonetworks.com/t5/general-topics/how-to-limit-concurrent-globalprotect-connections-per-user/td-p/202128
[^2]: https://live.paloaltonetworks.com/t5/api-articles/limit-maximum-globalprotect-vpn-sessions/ta-p/332846
