# Setup Dev Environmnet

## Virtual Midi Device 
```bash
$ sudo modprobe snd-virmidi snd_index=1
$ aconnect -io
$ aconnect 20:0 21:0
```
