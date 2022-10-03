# PAN-GPLimiter

A Go program for limiting concurrent remote user logins in a single GP Gateway on a PANOS Firewall.


## What's the motivation?

This one is maybe the most ever wanted feature request of Global Protect for decades! (FR4603-Concurrent Session Limiting) After tons of FR votes, endless requests from customers, lots of reddit messages asks for workarounds, PAN PMs don't have in the same opinion with the technical guys who are on the field as they haven't green lighted for developers to implement this super easy feature for years.

Finally, I ran out of hope and couldn't remain more indifferent to it. So this forces me to create my own home-brewed solution and I give myself the go-ahead.


## A Brief History:

Once I started to implement this program, there was only a PowerShell script [^1] dating from 2018.  I haven't tried it by myself but many ones couldn't make it run for some reason. (Or it really doesn't run at all!) Assuming it works, it's also OS (Windows) dependent, inefficient, couldn't handle edge-cases, lacks some features, etc... But besides that, it did its job as it inspired me and led the way to me!

After I created this program, I've found that someone else also created a Python script [^2] in 2020. I was surprised when faced with that since I didn't realize there was such an attempt at all. Honestly if I had known about it, I may never have started at first. You can also check this work since it provides some different features than this one.


## Why Golang?

Because it (cross) compiles into machine code! You can directly run ready-to-go binaries from Windows, Linux and Mac. No installation, no libraries, no dependencies, no prerequisites... Unlike Bash/PowerShell/Python, it's not interpreted on runtime which drastically increases runtime overhead for scripting languages. Decision of using a compiled language makes it run lightning fast with less memory usage. Also due to the statically typed nature of the Go language, it's more error-proof for any possible bugs/typos.


## TODO List:

Even I didn't think that those are show blockers, here is my todo list to make this program top-notch:

- [ ] Program argument for path to config file: to select config file different then working directory.
- [ ] Daemon mode: to get rid of Cron/Task Scheduler.
- [ ] Remove oldest sessions: to provide an option instead of kicking out the newest ones.
- [ ] Separate Mgmt IP support: to provide an alternative way to communicate with HA peers.
- [ ] Syslog and EventLog support: to write logs to OS when daemonized.


## Refs:
[^1]: https://live.paloaltonetworks.com/t5/general-topics/how-to-limit-concurrent-globalprotect-connections-per-user/td-p/202128
[^2]: https://live.paloaltonetworks.com/t5/api-articles/limit-maximum-globalprotect-vpn-sessions/ta-p/332846
