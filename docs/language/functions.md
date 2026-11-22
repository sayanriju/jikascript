# Functions in JikaScript

Functions in JikaScript are designed to feel natural and modern. You can write traditional function declarations or concise arrow functions.

---

## Function Declarations

Declare functions using standard JavaScript syntax:

```javascript
function calculateAverage(a, b) {
    return (a + b) / 2;
}
```

Notice that you don't even have to write parameter or return types! The compiler automatically figures them out based on how you call the function and what expressions you return.

---

## Call Before Declaration (Hoisting)

In embedded code, it's often nice to organize your high-level tasks at the top of the file and push helper functions to the bottom.

JikaScript pre-collects top-level functions, so you can call a function before its declaration:

```javascript
let result = compute(10); // Perfectly valid!

function compute(val) {
    return val * 2;
}
```

---

## Smart Parameter & Return Inference

The compiler checks all call sites to ensure consistent types:

```javascript
function scale(val) {
    return val * 2;
}

let a = scale(10);  // val inferred as int32
// let b = scale("hello"); // Compile error! Cannot pass string to int parameter
```

Return types are also validated across all branches:
- If you return an integer in one branch, returning a string in another branch will trigger a compile error.
- Functions with mixed bare (`return;`) and value (`return 42;`) returns are safely caught at compile time.

---

## Arrow Functions

Arrow functions can stand alone as task callbacks or be bound to variables for
reuse:

```javascript
const toggleLED = () => {
    led.toggle();
};

task("heartbeat", 500, toggleLED);

const square = x => x * x;
let answer = square(5); // 25
```

Named function declarations and arrow functions are statically resolvable. You
cannot dynamically construct either kind of function. In task callbacks, keep
the callback concise: it runs to completion in the main loop and one long
callback can delay the others.

Clean, readable, and zero boilerplate.

Next up: [Objects & Records](language/objects.md)!
