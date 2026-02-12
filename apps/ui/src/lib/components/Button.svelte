<script lang="ts">
  import RightArrow from "$lib/components/icons/RightArrowIcon.svelte";

  export let variant: "primary" | "secondary" = "primary";
  export let size: "sm" | "md" = "md";
  export let type: "button" | "submit" | "reset" = "button";
  export let disabled = false;
  export let href: string | undefined = undefined;
  // Deprecated, kept for compatibility
  export let color: string | undefined = undefined;

  const sizes = new Map<string, string>([
    ["md", "w-36 h-14 text-lg"],
    ["sm", "w-24 h-14 text-base"],
  ]);

  const variants = {
    primary: "bg-sky-400 hover:bg-sky-500 text-white border-transparent",
    secondary: "bg-white hover:bg-gray-50 text-gray-700 border-gray-300 border",
  };

  $: sizeClass = sizes.get(size) || sizes.get("md");
  $: variantClass = variants[variant] || variants.primary;
  $: btnClass = `rounded-full font-semibold font-sans transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center ${sizeClass} ${variantClass} ${$$props.class || ""}`;
</script>

{#if href && !disabled}
  <a {href} class={btnClass} role="button">
    <span class="flex flex-row gap-1 items-center justify-center">
      <span class="text-current">
        <slot />
      </span>
      <slot name="icon" />
    </span>
  </a>
{:else}
  <button {type} {disabled} class={btnClass} on:click>
    <span class="flex flex-row gap-1 items-center justify-center">
      <span class="text-current">
        <slot />
      </span>
      <slot name="icon" />
    </span>
  </button>
{/if}
