/*
 * Slideshow Configuration
 */
// One of https://github.com/isagalaev/highlight.js/tree/master/src/styles
const HIGHLIGHT_THEME = 'github';

// See options at https://daneden.github.io/animate.css/
const TRANSITIONS = {
  NEXT: 'animate__bounceInUp',
  PREV: 'animate__slideInLeft',
  INCREMENTAL: 'animate__fadeIn'
};

// Get the first title element of a slide and return its text content.
function getTitle(element) {
  const header = element.querySelector('h1, h2, h3, h4, h5');
  let title = '';

  if (header) {
    title = header.textContent;
  }

  return title;
}

// Create the slideshow on the page.
const slideshow = remark.create({
  highlightStyle: HIGHLIGHT_THEME,
  navigation: {
    scroll: false
  }
});

// Set up transitions between slides by monitoring which direction we are
// traveling and whether the slides are incremental (by checking titles).
slideshow.on('beforeShowSlide', (next) => {
  const nextIndex = next.getSlideIndex();
  const prevIndex = slideshow.getCurrentSlideIndex();
  const slides = document.querySelectorAll('.remark-slide-container');
  const nextDiv = slides[nextIndex];
  const prevDiv = slides[prevIndex];
  const nextTitle = getTitle(nextDiv);
  const prevTitle = getTitle(prevDiv);

  let direction = nextIndex > prevIndex ? TRANSITIONS.NEXT : TRANSITIONS.PREV

  if (prevTitle === nextTitle) {
    // Special case, either incremental slides or similar enough.
    direction = TRANSITIONS.INCREMENTAL;
  }

  nextDiv.classList.add('animate__animated');
  nextDiv.classList.add(direction);

  prevDiv.classList.remove('animate__animated');
  prevDiv.classList.remove(TRANSITIONS.NEXT);
  prevDiv.classList.remove(TRANSITIONS.PREV);
  prevDiv.classList.remove(TRANSITIONS.INCREMENTAL);
});

const term = document.querySelector('#terminal');

// Slide the term into or out of the viewport.
function toggleTerm() {
  term.classList.toggle('animate__slideInDown');
  term.classList.toggle('animate__slideOutUp');
}

// Toggle either from keypress or pressing the close button.
window.addEventListener('keyup', (event) => {
  //console.log(event);
  if (event.keyCode === 192 /* Key: ~ */) {
    toggleTerm();

    // If there is a selection, send it to the terminal. To make
    // the code examples a bit nicer, we don't need a `\` at the end
    // of each continued line. Instead, indent subsequent lines and this
    // will replace the newline + join with a single space.
    const selected = document
      .getSelection()
      .toString()
      .replace(/\n  /gm, ' ');

    let msg = '';
    if (selected) {
      msg = `clear\n${selected}\n`;
    }

    document
      .querySelector('iframe')
      .contentWindow
      .postMessage(msg, '*');
  }
});

// Toggle terminal when pressing the `close` button.
term.querySelector('a.btn').addEventListener('click', toggleTerm);

window.addEventListener('load', () => {
  setTimeout(() => {
    // Enable animations
    document.body.classList.remove('preload');
  }, 200);
});

// Setup MathJax
MathJax.Hub.Config({
  asciimath2jax: {
    // Since Markdown makes heavy use of backticks, prefer a syntax that
    // won't conflict with Markdown processing.
    delimiters: [['%%','%%']]
  }
});

MathJax.Hub.Configured();
