<script lang="ts">
    export let variant: "primary" | "secondary" | "danger" = "primary";
    export let size: "sm" | "md" = "sm";

    const variants = {
        primary: "bg-sky-400 text-white",
        secondary: "bg-slate-300 text-slate-900",
        danger: "bg-red-500 text-white"
    };

    const sizes = {
        sm: "px-2 py-1 text-xs",
        md: "px-3 py-2 text-sm"
    };
</script>

<span class="rounded-full font-semibold {variants[variant]} {sizes[size]} inline-block {$$props.class}">
    <slot />
</span>
