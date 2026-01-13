<script>
    import CodeEditor from "$lib/components/CodeEditor.svelte";

    let code = `
    import math

    print(math.sqrt(10))
    `
</script>

<CodeEditor code={code}></CodeEditor>