<script lang="ts">
    import * as monaco from "monaco-editor";
    import { onMount } from "svelte";

    let { code = "", language = "python" } = $props();

    let editorElement: HTMLDivElement;
    let editor: monaco.editor.IStandaloneCodeEditor;
    let model: monaco.editor.ITextModel;

    function loadCode(code: string, language: string) {
        model = monaco.editor.createModel(code, language);
        editor.setModel(model);
    }

    onMount(() => {
       editor = monaco.editor.create(editorElement);

       loadCode(code, language)
    });
</script>

<div class="flex-grow" bind:this={editorElement}></div>


