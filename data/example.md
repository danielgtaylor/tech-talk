# Sample Tech Talk

This is an example. Press <kbd>&#8680;</kbd> to continue.

[View example source](https://github.com/danielgtaylor/tech-talk/blob/master/data/example.md)

---

# Slide Title

Slides support Markdown content. **Strong** and *emphasized* text, [links](https://github.com/danielgtaylor/tech-talk), lists, tables, and code syntax highlighting.

1. List item
2. Another item
3. Just use standard Markdown!

```js
// Code example using Javascript
const pi = 3.14159;
let r = 10;

console.log(2 * pi * r * r);
```

[More formatting info](https://github.com/gnab/remark/wiki/Markdown)

---

# Navigation

Key | Description
--- | -----------
<kbd>&#8680;</kbd> | Move to next slide (spacebar also works)
<kbd>&#8678;</kbd> | Move to previous slide
<kbd>&#126;</kbd> | Show terminal (any selected text is copy/pasted)
<kbd>f</kbd> | Go to fullscreen mode
<kbd>esc</kbd> | Exit fullscreen mode
<kbd>p</kbd> | Go to presenter mode (to show slide notes)
<kbd>c</kbd> | Clone window (to display current slide)

---

# Incremental slide

Partial slide updates are also supported by using the `--` delimiter or by simply using the same slide title.

- A list
- Of items
--

- One more item
--

- And one last one!

But you aren't limited to just list items.

---

# Math & Formulas

Complex formulas are easy to display with [AsciiMath](http://asciimath.org/#syntax) using `%%`, or Tex / LaTeX using `$$` delimiters thanks to MathJax.

.center[
**AsciiMath example**
<!-- Notice that we escape the `*` because it has special meaning in Markdown -->
%%i = sum(1.65 \* 0.000125^(o - 1) \* (1 - 2.718^(-0.04t) / 4.15) \* (7490duz) / (100h))%%
]

.center[
**LaTeX example**
$$ax^2 + bx + c = 0$$
]

---

# Images & Videos

Assets in the same folder as the Markdown slides can be referenced relative to the root of the server.

.center[
![Amy Schumer](/static/www/amy.gif)]

```html
<img src="/my-image.png"/>
<video src="/my-video.mp4"/>
```

---
class: center, middle

Optional classes can control layout

Now go and give great tech talks!
