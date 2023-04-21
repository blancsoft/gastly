import {keymap, highlightSpecialChars, drawSelection, highlightActiveLine, dropCursor,
    rectangularSelection, crosshairCursor,
    lineNumbers, highlightActiveLineGutter} from "https://cdn.jsdelivr.net/npm/@codemirror/view@6.9.5/+esm"
import {EditorState} from "https://cdn.jsdelivr.net/npm/@codemirror/state@6.2.0/+esm"
import {defaultHighlightStyle, syntaxHighlighting, indentOnInput, bracketMatching,
    foldGutter, foldKeymap} from "https://cdn.jsdelivr.net/npm/@codemirror/language@6.6.0/+esm"
import {defaultKeymap, history, historyKeymap} from "https://cdn.jsdelivr.net/npm/@codemirror/commands@6.2.3/+esm"
import {searchKeymap, highlightSelectionMatches} from "https://cdn.jsdelivr.net/npm/@codemirror/search@6.3.0/+esm"
import {autocompletion, completionKeymap, closeBrackets, closeBracketsKeymap} from "https://cdn.jsdelivr.net/npm/@codemirror/autocomplete@6.5.1/+esm"
import {lintKeymap} from "https://cdn.jsdelivr.net/npm/@codemirror/lint@6.2.1/+esm"
import { EditorView } from "https://cdn.jsdelivr.net/npm/@codemirror/view@6.9.5/+esm";

// import { StreamLanguage } from "https://cdn.jsdelivr.net/npm/@codemirror/language@6.6.0/+esm";
// import { go } from "https://cdn.jsdelivr.net/npm/@codemirror/legacy-modes@6.3.2/mode/go.min.js"
// import { dracula } from "https://cdn.jsdelivr.net/npm/thememirror@2.0.1/+esm";

const basicSetup = (() => [
    lineNumbers(),
    highlightActiveLineGutter(),
    highlightSpecialChars(),
    history(),
    foldGutter(),
    drawSelection(),
    dropCursor(),
    EditorState.allowMultipleSelections.of(true),
    indentOnInput(),
    syntaxHighlighting(defaultHighlightStyle, {fallback: true}),
    bracketMatching(),
    closeBrackets(),
    autocompletion(),
    rectangularSelection(),
    crosshairCursor(),
    highlightActiveLine(),
    highlightSelectionMatches(),
    keymap.of([
        ...closeBracketsKeymap,
        ...defaultKeymap,
        ...searchKeymap,
        ...historyKeymap,
        ...foldKeymap,
        ...completionKeymap,
        ...lintKeymap
    ])
])()


window.basicSetup = basicSetup;
window.EditorState = EditorState;
window.EditorView = EditorView;