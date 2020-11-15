# Camera Rebooter

Reboot a glitchy UniFi camera daily.

```bash
go run . < config.json
```

See `config.example.json`;

- `address` is the UniFi Protect server (port is 7443) by default
- `username` and `password` are your login credentials for that server
- `camera_name` is what you've called your naughty camera
- `time` is the 24-hour time (HH:MM) for each reboot
- `location` is the time zone in which the reboot time should be scheduled 
