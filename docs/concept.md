# Concept

## Motivation
Pioneer utilizes a UDP/IP protocol called "Pro DJ Link" to ensure clock and beat synchronization between devices of their DJ/Nexus product line over standard ethernet. This makes it impossible to synchronize common MIDI instruments if you don`t own Toariz-SP16 sampler which can output a MIDI-Masteclock.

Altough there is [CDJ_Clock](https://github.com/g-zi/CDJ_Clock) which allows to sync a DAW to the Pro DJ Link clock, it requires OSX with a running DAW which adds a lot of complexity to the setup.

This project aims to create a small stand-alone Pro DJ Link to MIDI converter device.

## Pro DJ Link

### Physical
The devices are linked with standard CAT-5 ethernet cables.

### Protocol
The devices in link network communicate via udp.
A [protocol analysis](https://github.com/brunchboy/dysentery/raw/master/doc/Analysis.pdf) was done by [James Elliott](https://github.com/brunchboy) for the [dysentery](https://github.com/brunchboy/dysentery) project.

## Hardware concept
Small wallplug powered box with a ethernet-jack on one side and a MIDI-Out on the other. ;)

Raspi3 with pisound https://blokas.io/pisound (???)

## Related Projects
- [dysentery](https://github.com/brunchboy/dysentery) - Exploring ways to participate in a Pioneer Pro DJ Link network
- [CDJ_Clock](https://github.com/g-zi/CDJ_Clock) - With CDJ Clock anything what understands MIDI Beat Clock can be synced to Pioneer CDJs.
- [prolink-go](https://github.com/EvanPurkhiser/prolink-go) - golang library to interface with Pioneers PRO DJ Link network
- [beat-link](https://github.com/brunchboy/beat-link) - Java library for synchronizing with beats from Pioneer DJ Link equipment












