<script lang="ts">
    export let type: "info" | "success" | "warning" | "error" = "info";
    export let title: string = "";
    export let dismissible = false;

    let visible = true;

    const typeStyles = {
        info: {
            bg: "bg-sky-100",
            border: "border-sky-400",
            text: "text-sky-900",
            icon: "text-sky-400"
        },
        success: {
            bg: "bg-green-100",
            border: "border-green-400",
            text: "text-green-900",
            icon: "text-green-400"
        },
        warning: {
            bg: "bg-yellow-100",
            border: "border-yellow-400",
            text: "text-yellow-900",
            icon: "text-yellow-400"
        },
        error: {
            bg: "bg-red-100",
            border: "border-red-400",
            text: "text-red-900",
            icon: "text-red-400"
        }
    };

    const styles = typeStyles[type];
</script>

{#if visible}
    <div class="{styles.bg} border-l-4 {styles.border} {styles.text} p-4 rounded {$$props.class}">
        <div class="flex items-start justify-between">
            <div>
                {#if title}
                    <h3 class="font-semibold mb-1">{title}</h3>
                {/if}
                <p class="text-sm">
                    <slot />
                </p>
            </div>
            {#if dismissible}
                <button
                    on:click={() => (visible = false)}
                    class="ml-4 text-xl font-bold hover:opacity-70"
                >
                    ×
                </button>
            {/if}
        </div>
    </div>
{/if}
