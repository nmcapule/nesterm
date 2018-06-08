# NESTerm

NESTerm is a NES emulator targeted to run on terminals.

## Logs

### Day 1

These references are on the order of when I found them. Hopefully, I can also have a description on each item; what benefits were reaped and how easy was it.

Initially, I did a Google search: `how to write a NES emulator`, which directed me to a small write-up by Michael Fogleman:

- **I made a NES Emulator** (Michael Fogleman, 2015)

    [Link](https://medium.com/@fogleman/i-made-an-nes-emulator-here-s-what-i-learned-about-the-original-nintendo-2e078c9b28fe)

Which is, while a nice introduction, has a nicer reference link:

- **NES Documentation** (Patrick Diskin, 2004)

    [Link](http://nesdev.com/NESDoc.pdf)

This briefly detailed the history of NES and its successor, the SNES. Diving in...

> Log #1: WHoops. Lots of technical details and memory addressing and learned 
> a new word, bank switching.

> Log #2: Finished. Layed out the hardware specs pretty well + cartridge
> hardware specs too! Amazing piece of documentation really.

I did know that Fogleman had already implemented a very nice NES emulator, open-sourced at [github](https://github.com/fogleman/nes). I took a look and became amazed at the quality and readability of the code. I wanna take a stab at this too so I'll probably **NIH**.

Now for 6502 CPU instruction set:

- **6502 Instruction Reference**

    [Link](http://obelisk.me.uk/6502/reference.html#ROR)

Now I  can finally somewhat understand this document! Now onto some coding...

### Day 2

- Read up on 6502 addressing modes.
- Implementing Op Codes lookup tables and Addressing mode lookup table
- Awesome op code reference at [this link](http://www.emulator101.com/reference/6502-reference.html)
- First check-in of Go code
- Lots of manual work and cross-referencing! Darn...

### Day 3

- Read up on 6502 (2A03) memory mapping
- Implementing cpu.Memory struct!

> TODO:
>   - Implement CPU Op Codes
>   - That's a lot of Op Codes. Better get movin.

### Day 4

- Encoding processor status
- My first Op code (ADC)

> NOTE: WTF is the "S" flag?
>   S flag here: http://www.6502.org/tutorials/6502opcodes.html
>   No S flag here: http://nesdev.com/NESDoc.pdf

- **More 6502 Reference**

    [Link](http://nesdev.com/6502.txt)

    Helped finally know that the S flag is the Sign flag, aka., Negative Flag

- Having problems with the sign flag. Why the f does the AND op code has a sign flag change?

> TIL: Sign flag always corresponds to the most-significant bit

- Now having problems differentiating between the carry and the overflow flag.

- My actually first Op code (CLV) - Clear Overflow Flag