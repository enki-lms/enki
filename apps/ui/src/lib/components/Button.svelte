<script lang="ts">
  import RightArrow from "$lib/components/icons/RightArrowIcon.svelte";

  export let color: "blue" = "blue";
  export let size: "sm" | "md" | "lg" = "md";
  export let type: "button" | "submit" | "reset" = "button";
  export let disabled = false;

  const colors = new Map<string, string>([["blue", ""]]);

  const sizes = new Map<string, string>([
    ["lg", "w-48 h-16"],
    ["md", "w-36 h-14"],
    ["sm", "w-24 h-14"],
  ]);
</script>

<button
  {type}
  {disabled}
  class="rounded-full bg-sky-400 hover:bg-sky-500 font-semibold font-sans text-lg {sizes.get(
    size,
  )} text-white disabled:opacity-50 disabled:cursor-not-allowed {$$props.class}"
  on:click
>
  <span class="flex flex-row gap-1 items-center justify-center">
    <span class="text-black">
      <slot />
    </span>
    <slot name="icon" />
  </span>
</button>
