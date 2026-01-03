<script lang="ts">
    export let size: "sm" | "md" | "lg" = "md";
    export let color: string = "#77B6EA";

    const sizes = {
        sm: "w-4 h-4 border-2",
        md: "w-6 h-6 border-3",
        lg: "w-10 h-10 border-4"
    };
</script>

<div class="inline-block {sizes[size]} border-slate-300 border-t-[color:var(--spinner-color)] rounded-full animate-spin {$$props.class}"
     style="--spinner-color: {color}">
</div>
