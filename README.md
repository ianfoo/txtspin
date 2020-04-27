# txtspinner

A little Go program that implements a spinner for tty-based messages.

Configurable with command line args for animation speed, number of
cycles of the sample messages, the base length of the sample tasks, and
the characters in the animation.

### Tips

The animation loops through the sequence of characters in a single direction,
so you can get spinners using a set of distinct characters, or create a back
and forth by typing out the pattern in reverse except for the last character.

E.g., try `-<C<` and `.oOo` for some different but kinda cool effects.

