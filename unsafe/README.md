# MonkeyMock (unsafe extension) - Putting a square peg through a round hole, by force

Before you ever use this, make sure you've read the [CAVEATS](#CAVEATS)!

This extension provides the same basic functionality as [MonkeyMock](github.com/eshork/monkeymock), but with the added benefit of injecting your mock definitions directly into existing functions and objects (and even future objects), all without requiring you to substitute them within your application logic (typically via dependency injection/drilling).

In short:
- If you have a handle to a function, whether your own or imported from another library, you can force a mock layer upon it.
- If you have a handle to an existing object, whether your own or imported from another library, you can force a mock layer upon it.
- If you have a handle to a struct type, whether your own or imported from another library, you can force a mock layer upon every object instance created from it.

Here's the [CANDY](#Usage)

## But wait!

### Why is this `unsafe`?

Go (golang) is a compiled language; by typical definition, compiled code is immutable at runtime. In reality, all running code is held somewhere within memory, which is volotile and (most importantly) writable. Subverting the original intention of a program at runtime is surprisingly easy once you know how, and that's what this will do. You're effectively exploiting/hacking your code, voluntarily, for the purpose of making better tests.

The Go standard library provides an [`unsafe`](https://golang.org/pkg/unsafe) module, which this submodule uses/abuses to provide most capabilities. To make using the standard `unsafe` module easier, this submodule heavily leverages [bouk/monkey](https://github.com/bouk/monkey), hence the _monkey_ reference in the overall name.

Putting all of these bits together, it made sense to call a module like this `monkeymock/unsafe`, which actually then influenced the original top-level module name of _MonkeyMock_.



# CAVEATS
**Ie: What's the catch?**
- For consistency, you _absolutely must_ disable code inlining compiler optimizations during test builds ([disable inlining during tests](#How_to_disable_inlining_during_tests))
- not safe for production - don't ever do this if you want your code to be reliable (why?)
- some platforms have very strict memory alteration protections, so it wont always work (what?)

## How to disable inlining during tests
reasons and hows

----

# Usage
