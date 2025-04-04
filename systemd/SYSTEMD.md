# Systemd Configuration

## Users

```bash
adduser mbooks --shell /usr/sbin/nologin --home /var/www/damned.dev
gpasswd -d mbooks users
chmod 755 /var/www/damned.dev/
cd /var/www/damned.dev/
rm ~mbooks/.bash* ~mbooks/.cloud-locale-test.skip ~mbooks/.profile
```

## Files

```bash
root@damned.dev:/etc/systemd/system# ll /etc/systemd/system/mbooks*.service
-rw-r--r-- 1 root root 421 Sep 30 22:06 mbooks.service

root@damned.dev:/var/www/damned.dev/bin# ll /var/www/damned.dev/bin
-rwxr-xr-x 1 mbooks mbooks 18013459 Oct 13 21:57 musicbook-hub
```

## Install

```bash
root@damned.dev:/etc/systemd/system# systemctl daemon-reload

root@damned.dev:/etc/systemd/system# systemctl status mbooks.service
○ mbooks.service - Musicbooks Web Service
     Loaded: loaded (/etc/systemd/system/mbooks.service; disabled; preset: enabled)
     Active: inactive (dead)

root@damned.dev:/etc/systemd/system# systemctl enable mbooks.service
Created symlink '/etc/systemd/system/multi-user.target.wants/mbooks.service' → '/etc/systemd/system/mbooks.service'.

root@damned.dev:/etc/systemd/system# systemctl status mbooks.service
○ mbooks.service - Musicbooks Web Service
     Loaded: loaded (/etc/systemd/system/mbooks.service; enabled; preset: enabled)
     Active: inactive (dead)

root@damned.dev:/etc/systemd/system# systemctl start mbooks.service

root@damned.dev:/etc/systemd/system# systemctl status mbooks.service
● mbooks.service - Musicbooks Web Service
     Loaded: loaded (/etc/systemd/system/mbooks.service; enabled; preset: enabled)
     Active: active (running) since Thu 2025-03-13 17:39:31 UTC; 3s ago
 Invocation: 07cb01b33aa74daea8ce7b3f95479665
   Main PID: 9889 (mbooks)
      Tasks: 6 (limit: 1109)
     Memory: 1.5M (peak: 1.6M)
        CPU: 18ms
     CGroup: /system.slice/mbooks.service
             └─9889 /var/www/damned.dev/bin/mbooks start server

Mar 13 17:39:31 damned.dev systemd[1]: Started mbooks.service - Musicbooks Web Service.

root@damned.dev:/etc/systemd/system# systemctl stop mbooks.service
root@damned.dev:/etc/systemd/system# systemctl status mbooks.service
○ mbooks.service - Musicbooks Web Service
     Loaded: loaded (/etc/systemd/system/mbooks.service; enabled; preset: enabled)
     Active: inactive (dead) since Thu 2025-03-13 17:41:00 UTC; 3s ago
   Duration: 1min 29.064s
 Invocation: 07cb01b33aa74daea8ce7b3f95479665
    Process: 9889 ExecStart=/var/www/damned.dev/bin/mbooks start server (code=exited, status=0/SUCCESS)
   Main PID: 9889 (code=exited, status=0/SUCCESS)
   Mem peak: 1.7M
        CPU: 20ms

Mar 13 17:41:00 damned.dev systemd[1]: Stopping mbooks.service - Musicbooks Web Service...
Mar 13 17:41:00 damned.dev mbooks[9889]: server.go:85: server: signal terminated: received after 1m29.013996123s
Mar 13 17:41:00 damned.dev mbooks[9889]: server.go:90: server: timeout 5s: creating context (137ns)
Mar 13 17:41:00 damned.dev mbooks[9889]: server.go:95: server: canceling idle connections (43.253µs)
Mar 13 17:41:00 damned.dev mbooks[9889]: server.go:98: server: shutting down the server (52.534µs)
Mar 13 17:41:00 damned.dev mbooks[9889]: server.go:103: server: ¡stopped gracefully! (82.689µs)
Mar 13 17:41:00 damned.dev mbooks[9889]: start.go:72: server: shut down after 1m29.016373338s
Mar 13 17:41:00 damned.dev systemd[1]: mbooks.service: Deactivated successfully.
Mar 13 17:41:00 damned.dev systemd[1]: Stopped mbooks.service - Musicbooks Web Service.
```

## Monitor

```bash
root@damned.dev:/etc/systemd/system# journalctl -f -u mbooks.service
```

## Not used

```bash
# systemctl daemon-reexec
# systemctl enable mbooks.service
# systemctl start mbooks.service
```
