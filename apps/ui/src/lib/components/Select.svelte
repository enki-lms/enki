<script lang="ts">
  export let id: string | undefined = undefined;
  export let value: string | number = "";
  export let placeholder = "Select...";
  export let disabled = false;
  export let required = false;

  interface Option {
    value: string | number;
    label: string;
  }

  export let options: Option[] = [];
</script>

<div class="relative">
  <select
    {id}
    bind:value
    {disabled}
    {required}
    class="w-full px-4 py-3 pr-10 border-2 border-gray-200 rounded-lg focus:outline-none focus:border-sky-400 transition-colors bg-white appearance-none cursor-pointer disabled:bg-gray-100 disabled:cursor-not-allowed"
  >
    <option value="" disabled>{placeholder}</option>
    {#each options as option}
      <option value={option.value}>{option.label}</option>
    {/each}
  </select>
  <div
    class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none text-gray-500"
  >
    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M19 9l-7 7-7-7"
      />
    </svg>
  </div>
</div>
