<script lang="ts">
  import * as monaco from "monaco-editor";
  import { onMount } from "svelte";

  let {
    code = $bindable(""),
    language = "python",
    onSubmit = (_code: string) => {},
  } = $props();

  let editorElement: HTMLDivElement;
  let editor: monaco.editor.IStandaloneCodeEditor;
  let model: monaco.editor.ITextModel;

  function loadCode(code: string, language: string) {
    if (!editor) return;

    model = monaco.editor.createModel(code, language);
    editor.setModel(model);

    model.onDidChangeContent(() => {
      code = model.getValue();
    });
  }

  function handleSubmit() {
    const currentCode = editor.getValue();
    onSubmit(currentCode);
  }

  onMount(() => {
    editor = monaco.editor.create(editorElement, {
      minimap: {
        enabled: false,
      },
      automaticLayout: true,
    });

    loadCode(code, language);

    return () => {
      editor?.dispose();
      model?.dispose();
    };
  });

  $effect(() => {
    if (editor && code !== undefined && code !== editor.getValue()) {
      const position = editor.getPosition();
      editor.setValue(code);
      if (position) editor.setPosition(position);
    }
  });

  $effect(() => {
    if (editor && language) {
      const model = editor.getModel();
      if (model) {
        monaco.editor.setModelLanguage(model, language);
      }
    }
  });
</script>

<div class="flex flex-col w-full h-full">
  <div class="flex-1" bind:this={editorElement}></div>
  <div class="p-4 bg-gray-100 border-t border-gray-200 flex justify-end">
    <button
      class="rounded-full bg-sky-400 hover:bg-sky-500 font-semibold font-sans text-lg px-8 py-3 text-white transition-colors"
      onclick={handleSubmit}
    >
      Submit
    </button>
  </div>
</div>
