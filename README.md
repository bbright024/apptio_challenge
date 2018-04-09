# Apptio Coding Problem A

## Prologue

   This README file got pretty long, so I decided to write a little background
on WHY it's so long.  My girlfriend is a horse trainer, and I like to tag along
and spectate.  There's lots of downtime at horse shows, hence lots of time to
brainstorm and study.  Right now I'm typing this up in a farm house at an
equestrian eventing farm out in Kennewick, WA.

   I want to be as verbose as possible about my design choices and thought
process in this assignment, because for me it's the best way to get quality
feedback.  Even if I don't get this job, I still get 4-5 hours at Apptio HQ to
pick your guys' brains to find out what you all would've done, and I want to make
every second of that opportunity count.  

    Anyway, I decided to do problem A.  Problem B is pretty similar to exercise
4.13 in "The Go Programming Language", and I spent some time a few weeks ago
digesting chapter 4 and working on 4.10 & 4.11.  Those exercises were about
interacting with Go JSON library and the github API, so doing something similar
for this coding assignment would feel a bit like cheating.  Problem A has a
bunch of stuff that I haven't worked on in Go yet, and the question forces me
to think like I was part of a team.  I spent hours brainstorming and planning
solutions to different scenarios, and I'm really excited to see what you guys
think about some of my ideas.

Thanks for your time!

-Brian Bright


## Problem A


	The app devs have created a web app that generates a log file of the format:

	<datetime logtime, string "message">

	Your team has deployed the app in production, where the devs don't have access
	to the log file on the servers local file system due to a security boundary.

	The app devs want access to the contents of the log.

	Write automation/program that will allow the devs to view the logs
	without manual intervention.

## Initial Thoughts

    There are a bunch of ways to tackle this issue.  A few are ruled out due to
the need for automation, and a few others because of manual intervention, a few
more due to their complexity, and a few due to security.

The problem as stated was pretty vague, so I made a few assumptions.

   a) a server admin has root access to every main app host machine
   b) giving all the devs root access to the server isn't a viable solution
   c) Archived log files or the host machine running out of disk space is
       beyond the scope of this assignment
   d) each log entry ends in a new line
   
Please let me know if any of my assumptions are incorrect!  

## Design & Implementation

    My plan of attack is to install an HTTP server on the main app host machine
that would listen on a port that could not be accessed from outside the
firewall.  Pretty simple, just add a rule in iptables/ACL preventing WAN access
on that port, and test it later from the outside with noisy nmap scans.  Thus
the app would only respond to GET requests originating from inside the
organization, and could reply with JSON, plain text, or HTML.

   However, that HTTP design completely depends on the devs accessing the
server from a private IP address.  If the main app were in a container out in
the cloud, this would be a little trickier.  Port forwarding/NAT could be done
for requests coming from a specific range of IP addresses, but that stuff can
be spoofed easily.  Another option would be to write a client side proxy log
server that uses an SSH tunnel to connect with the log server on the main
machine.  A client side proxy would solve many other issues as well, by sort of
trunking requests from devs so that the main machine log server would only need
to track one TCP connection.

    Obviously, there are tons of issues I could work on.  However, the
assignment isn't supposed to be a month long affair, and writing something
overly complex might be frowned upon. My goal is to see what I can write in 3-4
hours, and add stuff on from there.  Overall I expect my time spent on this
assignment to be fairly high, because I'm not an expert in Go and I want to
sharpen my skills.

    Writing a client-side proxy is definitely a goal of mine, if I get that
far. Log entries will never be changed (unless an attacker takes control).  In
the client-side program, caching would reduce network load.  In the server side
program, caching would reduce system calls & disk I/O requests, leaving more
resources available for the main app.

Further ideas/issues:
     - have the server take a conf file so the admin doesn't need a long string of flags
        for each startup of the program
     - big read requests of the log file could slow down the main app
         Possible solutions:
           - add a limit to how much a user can access at one time
	   - only run the search when the CPU isn't under heavy load
	   - before any requests, cache as much of the log file into memory as possible
	   - cache results (both client side and server side)
	   - depending on the OS, change the 'nice' level of the log server process
     - the server side app will need read access on the log file and any
        archived log files.
     - network traffic of sending log search results could impact performance for main app
     - upgrades to the server side program API shouldn't break clients that assume
         an older API
     - if using a proxy, could batch requests & send as one to the log server
     - filtering the log file kinda reminds me of MapReduce; curious if there's anything I can
        do with that. 

#### Language choice: Go

The server will be written in Go, for a few reasons.
   1. Since Go compiles statically linked binaries, library dependencies
      	will not be an issue, making both programs easier to install.
   2. Go tools make it easy to compile into many different architectures/OS's 
   3. Go has really nice built-in concurrency that the server side program
        can take advantage of if multiple requests come in simultaneously
   4. The proxy server I wrote in C is pretty similar - copy/pasting would be no fun.
   5. Go is a lot more secure than anything I'd write in C

#### Security aside

    The devs put this app in production - so anyone on the web can access it,
right?  And the log file is being stored on the SAME MACHINE as the app, right?
So, if an attacker were able to compromise the app, even if the attacker
couldn't get root the log file could be modified (to varying degrees, depding
on if the app had read/write privileges or just append only). Worst case
scenario, would be if the attacker were to gain root on the machine, nothing in
the log file would be trustworthy.  The situation makes me think of a race
condition: until the log file is secure, it can't be determined if any
information in it is accurate.

    Thus my first step would be to cross my fingers & pray that the app hasn't been
compromised yet.  If I could, I'd write a script to install & configure rsyslog
on every deployed server to get those log files stored remotely ASAP.  If
rsyslog wasn't an option, I'd write a small app that would monitor the log file
for changes and write those changes to a remote destination, making sure that
the receiving app at the destination could only append log files.  There's
still a vulnerability that the receiving log app would be vulnerable to
malformed RPC from a compromised main app, but it's still a lot more secure of
an implementation.

    Remote storage is beneficial in other ways, beyond security. If the main app
crashes or the host is DDOS'ed, you'd want the log files stored remotely so
that the devs could see what led up to the outage.

HOWEVER, this stuff is obviously beyond the scope of the assignment, but it was
really bugging me so I felt like I should write it down.


### Alternatives

 - log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("/usr/share/doc"))))

    found this on the golang.org/pkg/net/http page.  If the log file were used
        in that http.Dir call, the solution go program would be ~10 lines of
        code.
	-problem: pretty brittle.  If the log file gets large or
        archived, more code would need to be added to accomodate changes.
        However, lazyness is a virtue.  It's the simplest & quickest solution
        to the problem we're trying to solve.
	       
- Give every dev root access on the app server
    - obviously stupid.  Gotta limit the circle of trust; there's no telling what could happen
      when all devs have root. 

- Write a caching client proxy server that forwards dev requests.  If all devs go through
   the proxy, that cache is gonna be magnificent - it will reduce the CPU time and network
   bandwidth taken by the log server
   - problem:
	It's a little beyond the scope of the problem we're trying to solve. Without more
	information, it's impossible to know if caching log files external to the server is
	a good thing or a bad thing.
      
- Add rsyslog to the main app host machine & copy the log files to an external database
     -problem: 
     	       Not enough information is known about the situation.  Rsyslog is more
	       resource intensive than my small log server implementation - maybe smaller
	       is better. 

- Add Kerberos to the network & check tickets whenever log files are requested.
   Adding authentication/authorization to the log server would allow for a much
   more dynamic API.  Server admin specific commands could be added,
   e.g. change any set configurations.
      -problems:
        - increases complexity of the system as a whole
        - Even though main app customers wouldn't need to go through kerb, they might still
	  see slowdowns.  A dev might write a script that grabs log files every 10 seconds,
	  and overhead of decrypting/encrypting those messages by the server app on the host machine
	  could totally bog down the system.
	- This is beyond the scope of the main problem - getting the devs access to the log files.
	  Would definitely need a team of devs & testers to implement something like this.
    
- Chmod the log files to a group on the host OS, give that group read privileges,
   give devs an account on the server they can SSH into, and add that account to the group. a cron job
   could be added to the host to add log archive files to that group with read privileges.
     -problem:
         - assumes the security policy is ok with devs having ssh access to the machine.
	 - would be hard to track which dev is using the account and at what time

- Add in a sandbox program to let devs pass in their own functions to filter the log
     -problem:
	complexity.  Just thought it'd be a cool idea to do. 
        
### Deployment plan

Deploying such a server would depend heavily on the infrastructure design.  In
the most basic scenario, an admin could write a script that run like so:

a) scp the log server binary that is arch/OS specific to the host machine into
   the same namespace as the log files
    - and don't forget to scp the conf.json file as well
    
b) ssh to the host machine

c) chmod/chown the binary
    - make sure the server binary isn't root
    - create a group that only the binary & the log files are part of
    - give the group read only privileges to the log files

d) write either a cron job or a systemd service that adds any new or archived
   log files to the group

e) have cron or systemd start/restart the server when needed

f) On whatever machine/router connects to the WAN & controls NAT/port
    forwarding, add a rule to the ACL preventing access from WAN IP's to the
    port used by the server.
    
a) through e) would need to be done for each machine the app is run on, so the
script would have to take a list of IP's (either from the CLI or a conf file).

I'm still a beginner with Docker and everything else in the cloud stack, so
coming up with a deployment plan for that kind of infrastructure might take me
a week to learn.  Hopefully I can update this section more before the
interview.

### MONITORING/REPORTING/ALERTS

    Need to log every client IP and time of request/response.
Simplest solution: change the 'log' settings to output log
messages to an append-only file on the local file system.  Interestingly, it
doesn't have to be a local file - output could be sent to a socket, since
log.SetOutput takes an io.Writer interface as the parameter.  Easy remote
logging!


### CONFIGURATION

Parameters set - maximum cache size, maximum concurrent requests, nice level

Until implementing authentication & authorization of the incoming requests to
the log server, every configuration must be set during installation.  If any of
the devs could change log server settings, the settings would be pretty
meaningless.

### Testing

Use noisy nmap scans from outside the organization to make sure the port is
inaccessible from the WAN.
   
Obviously unit testing & coverage checks of the code before deployment.

test what happens when a request to the log server times out

Could use a fuzzer to see if any of the parsing Go does in the libraries
is broken

After passing all the small unit tests, we should test what happens when the
log server and the main app run together on the same machine.  The main app is
live and presumably, we shouldn't take it offline.  Therefore, we need to run a
copy of the main app on an internal server alongside the log server.  We need
to see what happens when both the main app and the log server are operating
under heavy request volume, so tests should be written to simulate both
actions.

### RESOURCES

So far, I've used only code I've seen in "TGPL" and on the golang.org/pkg/*
example sections.  I'm nervous about using libraries or code I don't
understand; I guess I don't want the tech debt involved. For example, I'm not a
fan of how I have to open the requested log file again while handling an HTTP
request.  The file descriptor should still be accessible, so I was confused
about the errors I was getting.  I looked online for a solution, but what I
found was a bunch of clever closure tricks and wrapper stuff.  I don't want to
use that kind of stuff yet until I understand more about how Go works behind
the scenes with the HTTP Handler interface.


My HTTP proxy server solves a similar problem.
   a) take requests from clients over a socket
   b) read in and parse the requests
   c) access the requested resources
   d) write the results back to the client socket
   e) cache results for quicker access next time

"The Go Programming Language" - entire book really

MIT 6.824 Lecture 7 - Guest lecture on Go by Russ Cox of Google/Go

Countless man pages

golang.org/pkg/*

MIT 6.858 Security - Lecture 4: Privelege Seperation in Web Apps
    https://www.youtube.com/watch?v=XnBJc3-N2BU

Googled 'conf files in golang' trying to find out if there was a conf
parsing library for Go, but most comments I found said to use JSON.
