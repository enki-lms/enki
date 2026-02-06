<script lang="ts">
  export let src: string | undefined = undefined;
  export let name: string = "";
  export let width: string = "80";
  export let height: string = "80";

  function getInitials(fullName: string): string {
    const parts = fullName.trim().split(/\s+/);
    if (parts.length === 0 || !parts[0]) return "?";
    if (parts.length === 1) return parts[0].charAt(0).toUpperCase();
    return (
      parts[0].charAt(0) + parts[parts.length - 1].charAt(0)
    ).toUpperCase();
  }

  $: initials = getInitials(name);
  $: fontSize = Math.max(parseInt(width) / 2.5, 12);
</script>

<div
  class="rounded-full overflow-hidden border border-slate-300 flex items-center justify-center bg-gradient-to-br from-blue-500 to-indigo-600 text-white font-semibold"
  style="width: {width}px; height: {height}px; min-width: {width}px; min-height: {height}px; font-size: {fontSize}px;"
>
  {#if src}
    <img
      {src}
      alt="Profile Avatar"
      class="size-full object-cover rounded-full"
    />
  {:else}
    {initials}
  {/if}
</div>
